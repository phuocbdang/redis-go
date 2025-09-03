package core_test

import (
	"fmt"
	"redisgo/internal/core"
	"testing"
)

func TestSimpleStringDecode(t *testing.T) {
	cases := map[string]string{
		"+OK\r\n": "OK",
	}

	for k, v := range cases {
		value, _ := core.Decode([]byte(k))
		if value != v {
			t.Fail()
		}
	}
}

func TestError(t *testing.T) {
	cases := map[string]string{
		"-Error message\r\n": "Error message",
	}

	for k, v := range cases {
		value, _ := core.Decode([]byte(k))
		if value != v {
			t.Fail()
		}
	}
}

func TestInt64(t *testing.T) {
	cases := map[string]int64{
		":0\r\n":     0,
		":1000\r\n":  1000,
		":+1000\r\n": 1000,
		":-1000\r\n": -1000,
	}

	for k, v := range cases {
		value, _ := core.Decode([]byte(k))
		if value != v {
			t.Fail()
		}
	}
}

func TestBulkStringDecode(t *testing.T) {
	cases := map[string]string{
		"$5\r\nhello\r\n": "hello",
		"$0\r\n\r\n":      "",
	}

	for k, v := range cases {
		value, _ := core.Decode([]byte(k))
		if value != v {
			t.Fail()
		}
	}
}

func TestArrayDecode(t *testing.T) {
	cases := map[string][]interface{}{
		"*0\r\n":                                                    {},
		"*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n":                      {"hello", "world"},
		"*3\r\n:1\r\n:2\r\n:3\r\n":                                  {int64(1), int64(2), int64(3)},
		"*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$5\r\nhello\r\n":             {int64(1), int64(2), int64(3), int64(4), "hello"},
		"*2\r\n*3\r\n:1\r\n:2\r\n:-3\r\n*2\r\n+Hello\r\n-World\r\n": {[]int64{int64(1), int64(2), int64(-3)}, []interface{}{"Hello", "World"}},
	}

	for k, v := range cases {
		value, _ := core.Decode([]byte(k))
		arr := value.([]interface{})
		if len(arr) != len(v) {
			t.Fail()
		}
		for i := range arr {
			if fmt.Sprintf("%v", v[i]) != fmt.Sprintf("%v", arr[i]) {
				t.Fail()
			}
		}
	}
}

func TestString2DArrayEncode(t *testing.T) {
	decode := [][]string{{"hello", "world"}, {"1", "2", "3"}, {"xyz"}}
	encode := core.Encode(decode, false)
	if string(encode) != "*3\r\n*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n*3\r\n$1\r\n1\r\n$1\r\n2\r\n$1\r\n3\r\n*1\r\n$3\r\nxyz\r\n" {
		t.Fail()
	}
	decodeAgain, _ := core.Decode(encode)
	for i := 0; i < len(decode); i++ {
		for j := 0; j < len(decode[i]); j++ {
			if decodeAgain.([]interface{})[i].([]interface{})[j] != decode[i][j] {
				t.Fail()
			}
		}
	}
}

func TestParseCmd(t *testing.T) {
	cases := map[string]core.Command{
		"*3\r\n$3\r\nput\r\n$5\r\nhello\r\n$5\r\nworld\r\n": core.Command{
			Cmd:  "PUT",
			Args: []string{"hello", "world"},
		},
	}

	for k, v := range cases {
		value, _ := core.ParseCmd([]byte(k))
		if value.Cmd != v.Cmd {
			t.Fail()
		}
		if len(value.Args) != len(v.Args) {
			t.Fail()
		}
		for i := 0; i < len(value.Args); i++ {
			if value.Args[i] != v.Args[i] {
				t.Fail()
			}
		}
	}
}
