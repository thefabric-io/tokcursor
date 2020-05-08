package tokcursor //nolint:testpackage

import (
	"reflect"
	"testing"
)

func Test_keyValues(t *testing.T) {
	tests := []struct {
		name    string
		token   []byte
		want    map[string]string
		wantErr bool
	}{
		{"empty token", []byte(""), nil, false},
		{"well formatted token", []byte("key1:value1,key2:value2"), map[string]string{"key1": "value1", "key2": "value2"}, false},
		{"well formatted spaced token", []byte("key1: value1, key2: value2"), map[string]string{"key1": "value1", "key2": "value2"}, false},
		{"well formatted redundant token", []byte("key1: value1, key1: value1, key2: value2"), map[string]string{"key1": "value1", "key2": "value2"}, false},
		{"unordered token", []byte("key2: value2, key1: value1"), map[string]string{"key1": "value1", "key2": "value2"}, false},
		{"bad token format ':'", []byte("key2:va:lue2, key1: value1"), nil, true},
		{"bad token format ','", []byte("key2:value2,key1,value1"), nil, true},
	}
	for _, tt := range tests {
		token := tt.token
		wantErr := tt.wantErr
		want := tt.want
		t.Run(tt.name, func(t *testing.T) {
			got, err := keyValues(token)
			if (err != nil) != wantErr {
				t.Errorf("keyValues() error = %v, wantErr %v", err, wantErr)
			}
			if eq := reflect.DeepEqual(want, got); !eq {
				t.Errorf("keyValues() got = %v, want %v", got, want)
			}
		})
	}
}

func Test_cursor_Key(t *testing.T) {
	tests := []struct {
		name string
		kvs  map[string]string
		want string
	}{
		{"empty key values", map[string]string{}, ""},
		{"nil key values", nil, ""},
		{"ordered", map[string]string{"key1": "value1", "key2": "value2"}, "key1,key2"},
		{"unordered", map[string]string{"key2": "value2", "key1": "value1"}, "key1,key2"},
	}
	for _, tt := range tests {
		c := &B64Cursor{
			keyValues: tt.kvs,
		}
		want := tt.want
		t.Run(tt.name, func(t *testing.T) {
			got := c.Key()
			if got != want {
				t.Errorf("Key() got = %v  want %v", got, want)
			}
		})
	}
}

func Test_cursor_String(t *testing.T) {
	tests := []struct {
		name string
		kvs  map[string]string
		want string
	}{
		{"empty key values", map[string]string{}, ""},
		{"nil key values", nil, ""},
		{"ordered", map[string]string{"key1": "value1", "key2": "value2"}, "key1:value1,key2:value2"},
		{"unordered", map[string]string{"key2": "value2", "key1": "value1"}, "key1:value1,key2:value2"},
	}
	for _, tt := range tests {
		c := &B64Cursor{
			keyValues: tt.kvs,
		}
		want := tt.want
		t.Run(tt.name, func(t *testing.T) {
			got := c.String()
			if got != want {
				t.Errorf("String() got = %v  want %v", got, want)
			}
		})
	}
}
