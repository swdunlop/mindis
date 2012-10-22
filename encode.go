package mindis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

var crlf = []byte{'\r', '\n'}

func marshalRequest(verb string, args ...interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 4096))
	buf.WriteString(fmt.Sprint("*", len(args)+1, "\r\n$", len(verb), "\r\n", verb))
	for _, arg := range args {
		enc, err := marshalItem(arg)
		if err != nil {
			return nil, err
		}
		buf.WriteString(fmt.Sprint("\r\n$", len(enc), "\r\n"))
		buf.Write(enc)
	}
	buf.Write(crlf)
	return buf.Bytes(), nil
}

func marshalItem(arg interface{}) ([]byte, error) {
	switch v := arg.(type) {
	case Marshaler:
		return v.MarshalRedis()
	case json.Marshaler:
		return v.MarshalJSON()
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	case int:
		return []byte(strconv.FormatInt(int64(v), 10)), nil
	case int32:
		return []byte(strconv.FormatInt(int64(v), 10)), nil
	case int64:
		return []byte(strconv.FormatInt(v, 10)), nil
	case uint:
		return []byte(strconv.FormatUint(uint64(v), 10)), nil
	case uint32:
		return []byte(strconv.FormatUint(uint64(v), 10)), nil
	case uint64:
		return []byte(strconv.FormatUint(v, 10)), nil
	case bool:
		return []byte(strconv.FormatBool(v)), nil
	case float32:
		return []byte(strconv.FormatFloat(float64(v), 'g', -1, 32)), nil
	case float64:
		return []byte(strconv.FormatFloat(v, 'g', -1, 64)), nil
	}

	return nil, TYPE_ERROR
}

/* this interface indicates a value knows how to encode itself as bytes for transmission to redis */
type Marshaler interface {
	MarshalRedis() ([]byte, error)
}
