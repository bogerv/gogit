package cmds

import (
	"gitshell/colorlog"
	"gitshell/shellconst"
)

// CatchPanic panic recover
func CatchPanic() {
	if err := recover(); err != nil {
		colorlog.Error(shellconst.ErrNoBranch.Error())
	} else {
		return
	}
}
