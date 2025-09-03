package core

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

const CRLF = "\r\n"

var RespNil = []byte("$-1\r\n")

func encodeString(s string) []byte {
	return []byte(fmt.Sprintf("$%d%s%s%s", len(s), CRLF, s, CRLF))
}

func encodeStringArray(sa []string) []byte {
	var b []byte
	buf := bytes.NewBuffer(b)
	for _, s := range sa {
		buf.Write(encodeString(s))
	}
	return []byte(fmt.Sprintf("*%d%s%s", len(sa), CRLF, buf.Bytes()))
}

func Encode(value interface{}, isSimpleString bool) []byte {
	switch v := value.(type) {
	case string:
		if isSimpleString {
			return []byte(fmt.Sprintf("+%s%s", v, CRLF))
		}
		return encodeString(v)
	case int64, int32, int16, int8, int:
		return []byte(fmt.Sprintf(":%d%s", v, CRLF))
	case error:
		return []byte(fmt.Sprintf("-%s%s", v, CRLF))
	case []string:
		return encodeStringArray(v)
	case [][]string:
		var b []byte
		buf := bytes.NewBuffer(b)
		for _, sa := range v {
			buf.Write(encodeStringArray(sa))
		}
		return []byte(fmt.Sprintf("*%d%s%s", len(v), CRLF, buf.Bytes()))
	case []interface{}:
		var b []byte
		buf := bytes.NewBuffer(b)
		for _, x := range v {
			buf.Write(Encode(x, false))
		}
		return []byte(fmt.Sprintf("*%d%s%s", len(v), CRLF, buf.Bytes()))
	default:
		return RespNil
	}
}

func readSimpleString(data []byte) (string, int, error) {
	// +OK\r\n => OK, nextPos
	pos := 1
	for data[pos] != '\r' {
		pos++
	}
	return string(data[1:pos]), pos + 2, nil
}

func readInt64(data []byte) (int64, int, error) {
	// :123\r\n => 123, nextPos
	pos := 1
	var sign int64 = 1
	if data[pos] == '-' {
		sign = -1
		pos += 1
	}
	if data[pos] == '+' {
		pos += 1
	}
	var res int64 = 0
	for data[pos] != '\r' {
		res = res*10 + int64(data[pos]-'0')
		pos++
	}
	return sign * res, pos + 2, nil
}

func readError(data []byte) (string, int, error) {
	return readSimpleString(data)
}

func readLen(data []byte) (int, int) {
	// $5\r\nhello\r\n => 5, 4
	res, pos, _ := readInt64(data)
	return int(res), pos
}

func readBulkString(data []byte) (string, int, error) {
	// $5\r\nhello\r\n => hello, nextPos
	len, pos := readLen(data)
	return string(data[pos : pos+len]), pos + len + 2, nil
}

func readArray(data []byte) (interface{}, int, error) {
	// *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n => {"hello", "world"}
	len, pos := readLen(data)
	var res []interface{} = make([]interface{}, len)
	for i := range res {
		elem, delta, err := DecodeOne(data[pos:])
		if err != nil {
			return nil, 0, err
		}

		res[i] = elem
		pos += delta
	}
	return res, pos, nil
}

func DecodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("No data")
	}

	switch data[0] {
	case '+':
		return readSimpleString(data)
	case ':':
		return readInt64(data)
	case '-':
		return readError(data)
	case '$':
		return readBulkString(data)
	case '*':
		return readArray(data)
	}
	return nil, 0, nil
}

func Decode(data []byte) (interface{}, error) {
	res, _, err := DecodeOne(data)
	return res, err
}

func ParseCmd(data []byte) (*Command, error) {
	value, err := Decode(data)
	if err != nil {
		return nil, err
	}
	arr := value.([]interface{})
	tokens := make([]string, len(arr))
	for i := range tokens {
		tokens[i] = arr[i].(string)
	}
	return &Command{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil
}
