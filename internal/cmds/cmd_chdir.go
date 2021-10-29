package cmds

import (
	"log"
	"os"
)

// Chdir switch to directory where to execute the git command
func (slf *Cmd) Chdir() {
	err := os.Chdir(slf.Path)
	if err != nil {
		log.Printf(err.Error())
		os.Exit(500)
	}
}
