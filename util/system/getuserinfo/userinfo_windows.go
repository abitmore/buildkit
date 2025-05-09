package getuserinfo

import (
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/sys/windows"

	"github.com/moby/sys/reexec"
)

const (
	getUserInfoCmd = "get-user-info"
)

func init() {
	reexec.Register(getUserInfoCmd, userInfoMain)
}

func userInfoMain() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: get-user-info usernameOrGroup")
		os.Exit(1)
	}
	username := os.Args[1]
	sid, _, _, err := windows.LookupSID("", username)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	var ident struct {
		SID string
	}
	ident.SID = sid.String()

	asJSON, err := json.Marshal(ident)
	if err != nil {
		fmt.Println(err)
		os.Exit(5)
	}
	fmt.Fprintf(os.Stdout, "%s", string(asJSON))
}
