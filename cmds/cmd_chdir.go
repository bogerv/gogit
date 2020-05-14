package cmds

import "os"

// Chdir change dir
func (slf *Cmd) Chdir() {
	err := os.Chdir(slf.Path)
	if err != nil {
		panic(err)
	}
}
