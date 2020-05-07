package tokcursor

import (
	"bytes"
	"encoding/base64"
	"github.com/pkg/errors"
)

type b64cursor struct {
	token    string
	values   map[string]string
	pageSize int32
}

func (c b64cursor) Token() string {
	return c.token
}

func (c b64cursor) Values() map[string]string {
	return c.values
}

func (c b64cursor) PageSize() int32 {
	return c.pageSize
}

func (c *b64cursor) String() string {
	s := ""
	i := 0

	for k, e := range c.values {
		s += k + ":" + e
		if len(c.values) != i+1 {
			s += ","
		}
		i++
	}

	return s
}

func (c *b64cursor) parseToken() error {
	v, err := base64.StdEncoding.DecodeString(c.token)
	if err != nil {
		return err
	}

	v = bytes.Replace(v, []byte(" "), []byte(""), -1)
	bb := bytes.Split(v, []byte(","))
	c.values = make(map[string]string, len(bb))

	for i := range bb {
		bbb := bytes.Split(bb[i], []byte(":"))
		if len(bbb) != 2 {
			return errors.Wrapf(ErrTokenFormatIncorrect, "got %s", bb[i])
		}

		c.values[string(bbb[0])] = string(bbb[1])
	}

	return nil
}

func NewB64Cursor(token string, pageSize int32) (Cursor, error) {
	c := &b64cursor{
		token:    token,
		pageSize: pageSize,
	}

	if err := c.parseToken(); err != nil {
		return nil, err
	}

	return c, nil
}
