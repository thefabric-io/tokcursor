package tokcursor

import (
	"fmt"
	"github.com/pkg/errors"
)

var ErrTokenFormatIncorrect = errors.New("token format incorrect")

type Cursor interface {
	fmt.Stringer
	Token() string
	Values() map[string]string
	PageSize() int32
}
