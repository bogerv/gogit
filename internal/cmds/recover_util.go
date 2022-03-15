package cmds

import (
	"gogit/internal/constvar"
	"gogit/pkg/colorlog"
)

// CatchPanic panic recover
func CatchPanic() {
	if err := recover(); err != nil {
		colorlog.Error(constvar.ErrNoBranch.Error())
	} else {
		return
	}
}
