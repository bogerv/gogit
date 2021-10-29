package cmds

import (
	"gitshell/internal/constvar"
	"gitshell/pkg/colorlog"
)

// CatchPanic panic recover
func CatchPanic() {
	if err := recover(); err != nil {
		colorlog.Error(constvar.ErrNoBranch.Error())
	} else {
		return
	}
}
