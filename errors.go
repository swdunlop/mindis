package mindis

/* A simple Error abstraction using constant strings. */
type Error string

/* An implementation of the "error" interface. */
func (e Error) Error() string {
	return string(e)
}

/* Produced when an invalid type is passed to either Scan, Exec or Send. */
const TYPE_ERROR = Error("unsupported data type encountered")

/* Produced when an invalid response is read from the server. */
const PROTOCOL_ERROR = Error("an invalid response was read from the server")
