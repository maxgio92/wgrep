package grep

import (
	"github.com/pkg/errors"
)

var (
	ErrPatternNotSpecified  = errors.New("pattern not specified")
	ErrSeedURLsNotSpecified = errors.New("seed urls not specified")
	ErrSeedURLNotValid      = errors.New("seed url is not valid")
)
