package mounts

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/containerd/containerd/v2/core/mount"
	"github.com/moby/buildkit/cache"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/identity"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/secrets"
	"github.com/moby/buildkit/session/sshforward"
	"github.com/moby/buildkit/snapshot"
	"github.com/moby/buildkit/solver/pb"
	"github.com/moby/buildkit/util/bklog"
	"github.com/moby/buildkit/util/grpcerrors"
	"github.com/moby/locker"
	"github.com/moby/sys/user"
	"github.com/moby/sys/userns"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
)

func NewMountManager(name string, cm cache.Manager, sm *session.Manager) *MountManager {
	return &MountManager{
		cm:          cm,
		sm:          sm,
		cacheMounts: map[string]*cacheRefShare{},
		managerName: name,
	}
}

type MountManager struct {
	cm            cache.Manager
	sm            *session.Manager
	cacheMountsMu sync.Mutex
	cacheMounts   map[string]*cacheRefShare
	managerName   string
}

func (mm *MountManager) getRefCacheDir(ctx context.Context, ref cache.ImmutableRef, id string, m *pb.Mount, sharing pb.CacheSharingOpt, s session.Group) (mref cache.MutableRef, err error) {
	name := fmt.Sprintf("cached mount %s from %s", m.Dest, mm.managerName)
	if id != m.Dest {
		name += fmt.Sprintf(" with id %q", id)
	}
	g := &cacheRefGetter{
		locker:          &mm.cacheMountsMu,
		cacheMounts:     mm.cacheMounts,
		cm:              mm.cm,
		globalCacheRefs: sharedCacheRefs,
		name:            name,
		session:         s,
	}
	return g.getRefCacheDir(ctx, ref, id, sharing)
}

type cacheRefGetter struct {
	locker          sync.Locker
	cacheMounts     map[string]*cacheRefShare
	cm              cache.Manager
	globalCacheRefs *cacheRefs
	name            string
	session         session.Group
}

func (g *cacheRefGetter) getRefCacheDir(ctx context.Context, ref cache.ImmutableRef, id string, sharing pb.CacheSharingOpt) (mref cache.MutableRef, err error) {
	key := id
	if ref != nil {
		key += ":" + ref.ID()
	}
	mu := g.locker
	mu.Lock()
	defer mu.Unlock()

	if ref, ok := g.cacheMounts[key]; ok {
		return ref.clone(ctx), nil
	}
	defer func() {
		if err == nil {
			share := &cacheRefShare{MutableRef: mref, refs: map[*cacheRef]struct{}{}}
			g.cacheMounts[key] = share
			mref = share.clone(ctx)
		}
	}()

	switch sharing {
	case pb.CacheSharingOpt_SHARED:
		return g.globalCacheRefs.get(ctx, key, func() (cache.MutableRef, error) {
			return g.getRefCacheDirNoCache(ctx, key, ref, id, false)
		})
	case pb.CacheSharingOpt_PRIVATE:
		return g.getRefCacheDirNoCache(ctx, key, ref, id, false)
	case pb.CacheSharingOpt_LOCKED:
		return g.getRefCacheDirNoCache(ctx, key, ref, id, true)
	default:
		return nil, errors.Errorf("invalid cache sharing option: %s", sharing.String())
	}
}

func (g *cacheRefGetter) getRefCacheDirNoCache(ctx context.Context, key string, ref cache.ImmutableRef, id string, block bool) (cache.MutableRef, error) {
	makeMutable := func(ref cache.ImmutableRef) (cache.MutableRef, error) {
		newRef, err := g.cm.New(ctx, ref, g.session, cache.WithRecordType(client.UsageRecordTypeCacheMount), cache.WithDescription(g.name), cache.CachePolicyRetain)
		if err != nil {
			return nil, err
		}
		bklog.G(ctx).Debugf("created new ref for cache dir %q: %s", id, newRef.ID())
		return newRef, nil
	}

	cacheRefsLocker.Lock(key)
	defer cacheRefsLocker.Unlock(key)
	for {
		sis, err := SearchCacheDir(ctx, g.cm, key, false)
		if err != nil {
			return nil, err
		}
		locked := false
		for _, si := range sis {
			if mRef, err := g.cm.GetMutable(ctx, si.ID()); err == nil {
				bklog.G(ctx).Debugf("reusing ref for cache dir %q: %s", id, mRef.ID())
				return mRef, nil
			} else if errors.Is(err, cache.ErrLocked) {
				locked = true
			} else {
				bklog.G(ctx).WithError(err).Errorf("failed to get reuse ref for cache dir %q: %s", id, si.ID())
			}
		}
		if block && locked {
			cacheRefsLocker.Unlock(key)
			select {
			case <-ctx.Done():
				cacheRefsLocker.Lock(key)
				return nil, context.Cause(ctx)
			case <-time.After(100 * time.Millisecond):
				cacheRefsLocker.Lock(key)
			}
		} else {
			break
		}
	}
	mRef, err := makeMutable(ref)
	if err != nil {
		return nil, err
	}

	md := CacheRefMetadata{mRef}
	if err := md.setCacheDirIndex(key); err != nil {
		mRef.Release(context.TODO())
		return nil, err
	}
	return mRef, nil
}

func (mm *MountManager) getSSHMountable(ctx context.Context, m *pb.Mount, g session.Group) (cache.Mountable, error) {
	var caller session.Caller
	err := mm.sm.Any(ctx, g, func(ctx context.Context, _ string, c session.Caller) error {
		if err := sshforward.CheckSSHID(ctx, c, m.SSHOpt.ID); err != nil {
			if m.SSHOpt.Optional {
				return nil
			}
			if grpcerrors.Code(err) == codes.Unimplemented {
				return errors.Errorf("no SSH key %q forwarded from the client", m.SSHOpt.ID)
			}
			return err
		}
		caller = c
		return nil
	})
	if err != nil {
		return nil, err
	}
	if caller == nil {
		return nil, nil
	}
	// because ssh socket remains active, to actually handle session disconnecting ssh error
	// should restart the whole exec with new session
	return &sshMount{mount: m, caller: caller, idmap: mm.cm.IdentityMapping()}, nil
}

type sshMount struct {
	mount  *pb.Mount
	caller session.Caller
	idmap  *user.IdentityMapping
}

func (sm *sshMount) Mount(ctx context.Context, readonly bool, g session.Group) (snapshot.Mountable, error) {
	return &sshMountInstance{sm: sm, idmap: sm.idmap}, nil
}

type sshMountInstance struct {
	sm    *sshMount
	idmap *user.IdentityMapping
}

func (sm *sshMountInstance) Mount() ([]mount.Mount, func() error, error) {
	ctx, cancel := context.WithCancelCause(context.TODO())

	uid := int(sm.sm.mount.SSHOpt.Uid)
	gid := int(sm.sm.mount.SSHOpt.Gid)

	if sm.idmap != nil {
		var err error
		uid, gid, err = sm.idmap.ToHost(uid, gid)
		if err != nil {
			cancel(err)
			return nil, nil, err
		}
	}

	sock, cleanup, err := sshforward.MountSSHSocket(ctx, sm.sm.caller, sshforward.SocketOpt{
		ID:   sm.sm.mount.SSHOpt.ID,
		UID:  uid,
		GID:  gid,
		Mode: int(sm.sm.mount.SSHOpt.Mode & 0777),
	})
	if err != nil {
		cancel(err)
		return nil, nil, err
	}
	release := func() error {
		var err error
		if cleanup != nil {
			err = cleanup()
		}
		cancel(err)
		return err
	}

	return []mount.Mount{{
		Type:    "bind",
		Source:  sock,
		Options: []string{"rbind"},
	}}, release, nil
}

func (sm *sshMountInstance) IdentityMapping() *user.IdentityMapping {
	return sm.idmap
}

func (mm *MountManager) getSecretMountable(ctx context.Context, m *pb.Mount, g session.Group) (cache.Mountable, error) {
	if m.SecretOpt == nil {
		return nil, errors.Errorf("invalid secret mount options")
	}

	id := m.SecretOpt.ID
	if id == "" {
		return nil, errors.Errorf("secret ID missing from mount options")
	}
	var dt []byte
	var err error
	err = mm.sm.Any(ctx, g, func(ctx context.Context, _ string, caller session.Caller) error {
		dt, err = secrets.GetSecret(ctx, caller, id)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, secrets.ErrNotFound) && m.SecretOpt.Optional {
			return nil, nil
		}
		return nil, err
	}
	return &secretMount{mount: m, data: dt, idmap: mm.cm.IdentityMapping()}, nil
}

type secretMount struct {
	mount *pb.Mount
	data  []byte
	idmap *user.IdentityMapping
}

func (sm *secretMount) Mount(ctx context.Context, readonly bool, g session.Group) (snapshot.Mountable, error) {
	return &secretMountInstance{sm: sm, idmap: sm.idmap}, nil
}

type secretMountInstance struct {
	sm    *secretMount
	root  string
	idmap *user.IdentityMapping
}

func (sm *secretMountInstance) Mount() ([]mount.Mount, func() error, error) {
	dir, err := os.MkdirTemp("", "buildkit-secrets")
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create temp dir")
	}
	cleanupDir := func() error {
		return os.RemoveAll(dir)
	}

	if err := os.Chmod(dir, 0711); err != nil {
		cleanupDir()
		return nil, nil, err
	}

	var mountOpts []string
	if sm.sm.mount.SecretOpt.Mode&0o111 == 0 {
		mountOpts = append(mountOpts, "noexec")
	}

	tmpMount := mount.Mount{
		Type:    "tmpfs",
		Source:  "tmpfs",
		Options: append([]string{"nodev", "nosuid", fmt.Sprintf("uid=%d,gid=%d", os.Geteuid(), os.Getegid())}, mountOpts...),
	}

	if userns.RunningInUserNS() {
		tmpMount.Options = nil
	}

	if err := mount.All([]mount.Mount{tmpMount}, dir); err != nil {
		cleanupDir()
		return nil, nil, errors.Wrap(err, "unable to setup secret mount")
	}
	sm.root = dir

	cleanup := func() error {
		if err := mount.Unmount(dir, 0); err != nil {
			return err
		}
		return cleanupDir()
	}

	randID := identity.NewID()
	fp := filepath.Join(dir, randID)
	if err := os.WriteFile(fp, sm.sm.data, 0600); err != nil {
		cleanup()
		return nil, nil, err
	}

	uid := int(sm.sm.mount.SecretOpt.Uid)
	gid := int(sm.sm.mount.SecretOpt.Gid)

	if sm.idmap != nil {
		uid, gid, err = sm.idmap.ToHost(uid, gid)
		if err != nil {
			cleanup()
			return nil, nil, err
		}
	}

	if err := os.Chown(fp, uid, gid); err != nil {
		cleanup()
		return nil, nil, err
	}

	if err := os.Chmod(fp, os.FileMode(sm.sm.mount.SecretOpt.Mode&0777)); err != nil {
		cleanup()
		return nil, nil, err
	}

	return []mount.Mount{{
		Type:    "bind",
		Source:  fp,
		Options: append([]string{"ro", "rbind", "nodev", "nosuid"}, mountOpts...),
	}}, cleanup, nil
}

func (sm *secretMountInstance) IdentityMapping() *user.IdentityMapping {
	return sm.idmap
}

func (mm *MountManager) MountableCache(ctx context.Context, m *pb.Mount, ref cache.ImmutableRef, g session.Group) (cache.MutableRef, error) {
	if m.CacheOpt == nil {
		return nil, errors.Errorf("missing cache mount options")
	}
	return mm.getRefCacheDir(ctx, ref, m.CacheOpt.ID, m, m.CacheOpt.Sharing, g)
}

func (mm *MountManager) MountableTmpFS(m *pb.Mount) cache.Mountable {
	return newTmpfs(mm.cm.IdentityMapping(), m.TmpfsOpt)
}

func (mm *MountManager) MountableSecret(ctx context.Context, m *pb.Mount, g session.Group) (cache.Mountable, error) {
	return mm.getSecretMountable(ctx, m, g)
}

func (mm *MountManager) MountableSSH(ctx context.Context, m *pb.Mount, g session.Group) (cache.Mountable, error) {
	return mm.getSSHMountable(ctx, m, g)
}

func newTmpfs(idmap *user.IdentityMapping, opt *pb.TmpfsOpt) cache.Mountable {
	return &tmpfs{idmap: idmap, opt: opt}
}

type tmpfs struct {
	idmap *user.IdentityMapping
	opt   *pb.TmpfsOpt
}

func (f *tmpfs) Mount(ctx context.Context, readonly bool, g session.Group) (snapshot.Mountable, error) {
	return &tmpfsMount{readonly: readonly, idmap: f.idmap, opt: f.opt}, nil
}

type tmpfsMount struct {
	readonly bool
	idmap    *user.IdentityMapping
	opt      *pb.TmpfsOpt
}

func (m *tmpfsMount) Mount() ([]mount.Mount, func() error, error) {
	opt := []string{"nosuid"}
	if m.readonly {
		opt = append(opt, "ro")
	}
	if m.opt != nil {
		if m.opt.Size > 0 {
			opt = append(opt, fmt.Sprintf("size=%d", m.opt.Size))
		}
	}
	return []mount.Mount{{
		Type:    "tmpfs",
		Source:  "tmpfs",
		Options: opt,
	}}, func() error { return nil }, nil
}

func (m *tmpfsMount) IdentityMapping() *user.IdentityMapping {
	return m.idmap
}

var (
	cacheRefsLocker = locker.New()
	sharedCacheRefs = &cacheRefs{}
)

type cacheRefs struct {
	mu     sync.Mutex
	shares map[string]*cacheRefShare
}

// ClearActiveCacheMounts clears shared cache mounts currently in use.
// Caller needs to hold CacheMountsLocker before calling
func ClearActiveCacheMounts() {
	sharedCacheRefs.shares = nil
}

func CacheMountsLocker() sync.Locker {
	return &sharedCacheRefs.mu
}

func (r *cacheRefs) get(ctx context.Context, key string, fn func() (cache.MutableRef, error)) (cache.MutableRef, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.shares == nil {
		r.shares = map[string]*cacheRefShare{}
	}

	share, ok := r.shares[key]
	if ok {
		return share.clone(ctx), nil
	}

	mref, err := fn()
	if err != nil {
		return nil, err
	}

	share = &cacheRefShare{MutableRef: mref, main: r, key: key, refs: map[*cacheRef]struct{}{}}
	r.shares[key] = share
	return share.clone(ctx), nil
}

type cacheRefShare struct {
	cache.MutableRef
	mu   sync.Mutex
	refs map[*cacheRef]struct{}
	main *cacheRefs
	key  string
}

func (r *cacheRefShare) clone(ctx context.Context) cache.MutableRef {
	bklog.G(ctx).WithFields(map[string]any{
		"key":   r.key,
		"stack": bklog.TraceLevelOnlyStack(),
	}).Trace("cloning cache mount ref share")
	cacheRef := &cacheRef{cacheRefShare: r}
	if cacheRefCloneHijack != nil {
		cacheRefCloneHijack()
	}
	r.mu.Lock()
	r.refs[cacheRef] = struct{}{}
	r.mu.Unlock()
	return cacheRef
}

func (r *cacheRefShare) release(ctx context.Context) error {
	bklog.G(ctx).WithFields(map[string]any{
		"key":   r.key,
		"stack": bklog.TraceLevelOnlyStack(),
	}).Trace("releasing cache mount ref share main")
	if r.main != nil {
		delete(r.main.shares, r.key)
	}
	return r.Release(ctx)
}

var (
	cacheRefReleaseHijack func()
	cacheRefCloneHijack   func()
)

type cacheRef struct {
	*cacheRefShare
}

func (r *cacheRef) Release(ctx context.Context) error {
	bklog.G(ctx).WithFields(map[string]any{
		"key":   r.key,
		"stack": bklog.TraceLevelOnlyStack(),
	}).Trace("releasing cache mount ref share")
	if r.main != nil {
		r.main.mu.Lock()
		defer r.main.mu.Unlock()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.refs, r)
	if len(r.refs) == 0 {
		if cacheRefReleaseHijack != nil {
			cacheRefReleaseHijack()
		}
		return r.release(ctx)
	}
	return nil
}

const (
	keyCacheDir   = "cache-dir"
	cacheDirIndex = keyCacheDir + ":"
)

func SearchCacheDir(ctx context.Context, store cache.MetadataStore, id string, withNested bool) ([]CacheRefMetadata, error) {
	var results []CacheRefMetadata
	key := cacheDirIndex + id
	if withNested {
		key += ":"
	}
	mds, err := store.Search(ctx, key, withNested)
	if err != nil {
		return nil, err
	}
	for _, md := range mds {
		if withNested {
			v := md.Get(keyCacheDir)
			// skip partial ids but allow id without ref ID
			if v == nil || v.Index != key && !strings.HasPrefix(v.Index, key) {
				continue
			}
		}
		results = append(results, CacheRefMetadata{md})
	}
	return results, nil
}

type CacheRefMetadata struct {
	cache.RefMetadata
}

func (md CacheRefMetadata) setCacheDirIndex(id string) error {
	return md.SetString(keyCacheDir, id, cacheDirIndex+id)
}

func (md CacheRefMetadata) ClearCacheDirIndex() error {
	return md.ClearValueAndIndex(keyCacheDir, cacheDirIndex)
}
