package mindis

import (
	"encoding/json"
	"strconv"
)

func unmarshalItem(data []byte, arg interface{}) error {
	switch v := arg.(type) {
	case json.Unmarshaler:
		return v.UnmarshalJSON(data)
	case *string:
		*v = string(data)
	case *bool:
		b, err := strconv.ParseBool(string(data))
		*v = b
		return err
	case *int:
		n, err := strconv.ParseInt(string(data), 10, 64)
		*v = int(n)
		return err
	case *int32:
		n, err := strconv.ParseInt(string(data), 10, 32)
		*v = int32(n)
		return err
	case *int64:
		n, err := strconv.ParseInt(string(data), 10, 64)
		*v = n
		return err
	case *uint:
		n, err := strconv.ParseUint(string(data), 10, 64)
		*v = uint(n)
		return err
	case *uint32:
		n, err := strconv.ParseUint(string(data), 10, 32)
		*v = uint32(n)
		return err
	case *uint64:
		n, err := strconv.ParseUint(string(data), 10, 64)
		*v = n
		return err
	default:
		return TYPE_ERROR
	}

	return nil
}

func unmarshalData(reply [][]byte, arg interface{}) error {
	switch v := arg.(type) {
	case Unmarshaler:
		return v.UnmarshalRedis(reply)
	case *[]string:
		r := make([]string, len(reply))
		for n, b := range reply {
			r[n] = string(b)
		}
		*v = r
	case map[string]string:
		if len(reply)%2 != 0 {
			return TYPE_ERROR
		}
		for i := 0; i < len(reply); i += 2 {
			v[string(reply[i])] = string(reply[i+1])
		}
	default:
		if len(reply) == 1 {
			return unmarshalItem(reply[0], arg)
		}
		return TYPE_ERROR
	}

	return nil
}

/* this interface indicates a value knows how to decode itself from bytes received from redis */
type Unmarshaler interface {
	UnmarshalRedis([][]byte) error
}
