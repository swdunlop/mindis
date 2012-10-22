package mindis

/* A simple Error abstraction using constant strings. */

type Error string

func (e Error) Error() string {
	return string(e)
}

const TYPE_ERROR = Error("unsupported data type encountered")
const PROTOCOL_ERROR = Error("an invalid response was read from the server")
