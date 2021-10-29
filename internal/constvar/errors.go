package constvar

import "github.com/go-errors/errors"

var (
	ErrNoBranch = errors.New("no branch")
	ErrNoFlow   = errors.New("no flows")
)
