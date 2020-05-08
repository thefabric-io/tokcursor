package tokcursor

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"sort"
)

var ErrTokenFormatIncorrect = errors.New("token format incorrect")

type Cursor interface {
	fmt.Stringer
	Token() string
	RawToken() string
	KeyValues() map[string]string
	Key() string
	PageSize() int32
}

type cursor struct {
	rawToken  string
	token     []byte
	keyValues map[string]string
	pageSize  int32
}

func (c cursor) Key() string {
	s := ""
	i := 0
	keys := make([]string, len(c.keyValues))

	for k := range c.keyValues {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	for _, k := range keys {
		s += k + ","
	}

	if len(s) != 0 {
		s = s[:len(s)-1]
	}

	return s
}

func (c *cursor) String() string {
	s := ""
	i := 0

	keys := make([]string, len(c.keyValues))

	for k := range c.keyValues {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	for _, e := range keys {
		s += e + ":" + c.keyValues[e] + ","
	}

	if len(s) != 0 {
		s = s[:len(s)-1]
	}

	return s
}

func (c cursor) Token() string {
	return string(c.token)
}

func (c cursor) RawToken() string {
	return c.rawToken
}

func (c cursor) KeyValues() map[string]string {
	return c.keyValues
}

func (c cursor) PageSize() int32 {
	return c.pageSize
}

func keyValues(token []byte) (map[string]string, error) {
	v := bytes.Replace(token, []byte(" "), []byte(""), -1)
	if len(v) == 0 {
		return nil, nil
	}

	bb := bytes.Split(v, []byte(","))
	keyValues := make(map[string]string, len(bb))

	for i := range bb {
		bbb := bytes.Split(bb[i], []byte(":"))
		if len(bbb) != 2 {
			return nil, errors.Wrapf(ErrTokenFormatIncorrect, "got %s", bb[i])
		}

		keyValues[string(bbb[0])] = string(bbb[1])
	}

	return keyValues, nil
}

func NewCursor(token string, pageSize int32) (Cursor, error) {
	c := &cursor{
		token:    []byte(token),
		pageSize: pageSize,
	}
	c.rawToken = base64.StdEncoding.EncodeToString([]byte(token))

	kv, err := keyValues(c.token)
	if err != nil {
		return nil, err
	}

	c.keyValues = kv

	return c, nil
}
