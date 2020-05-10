package tokcursor

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"sort"

	"github.com/pkg/errors"
)

// Cursor is a cursor interface representing a pagination cursor.
type Cursor interface {
	fmt.Stringer
	Token() string
	RawToken() string
	KeyValues() map[string]string
	Key() string
	PageSize() int32
}

// NewB64Cursor returns a base64 key:value cursor implementation.
func NewB64Cursor(token string, pageSize int32) (Cursor, error) {
	c := &B64Cursor{
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

// B64Cursor represents a base64 key:value cursor implementation. It contains:
// - raw token, a base64 encoded representation of the key:value token
// - token, a []byte representation of the key:value values i.e. []byte("key1:value1,key2:value2")
// - keyValues, a map[string]string structured key:value pairs of the cursor
// - pageSize, an int32 representing the size of the page.
type B64Cursor struct {
	rawToken  string
	token     []byte
	keyValues map[string]string
	pageSize  int32
}

// Token returns a string representation of the key:value token.
func (c B64Cursor) Token() string {
	return string(c.token)
}

// RawToken returns a string representation of the base64 raw token.
func (c B64Cursor) RawToken() string {
	return c.rawToken
}

// KeyValues returns the map[string]string structured key:value pairs of the cursor.
func (c B64Cursor) KeyValues() map[string]string {
	return c.keyValues
}

// PageSize returns the page size of the pagination.
func (c B64Cursor) PageSize() int32 {
	return c.pageSize
}

// Key returns a string representing an unique key representation of the cursor structure
// For example, for both cursors "key2:value2, key1:value1" and "key1:value1, key2:value2"
// the returned value will be "key1,key2".
func (c B64Cursor) Key() string {
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

// String returns the readable token including all key:value pairs excluding the redundant keys.
func (c B64Cursor) String() string {
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

// ErrTokenFormatIncorrect is returned when the key:value element in a token decoded
// string is incorrect i.e "key:value:key, key2:value2" is incorrect because of the
// incorrect format of the first element "key:value:key".
var ErrTokenFormatIncorrect = errors.New("token format incorrect")
