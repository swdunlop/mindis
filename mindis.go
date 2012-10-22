package mindis

import (
	"bufio"
	"io"
)

type Conn struct {
	Rwc    io.ReadWriteCloser // the initial ReadWriteCloser provided to Wrap
	r      *bufio.Reader      // our buffered reader
	status error              // internal storage for Status()
	lag    int                // internal count of how many responses we expect for sync
	result [][]byte           // the most recent result, as a list of bytes
}

/* 
Wraps an io.ReadWriteCloser with the necessary state information to interact with Redis. 
*/
func Wrap(rwc io.ReadWriteCloser) *Conn {
	return &Conn{rwc, bufio.NewReader(rwc), nil, 0, nil}
}

/* 
Exec will Sync() to drop old results, Send() a new command, then Next() to check for a response; if ret is non-nil, Scan() will also be used to provide it with a value. 
*/
func (cn *Conn) Exec(ret interface{}, verb string, args ...interface{}) error {
	cn.Sync()
	err := cn.Send(verb, args...)
	if err != nil {
		return err
	}
	cn.Next()
	return cn.Scan(ret)
}

/* 
Scan decodes the current result to src, where Next() or Exec() was used to advance to the next result.  If Status() was not nil, Scan() will stubbornly repeat that error.
*/
func (cn *Conn) Scan(ret interface{}) error {
	switch {
	case cn.status != nil:
		return cn.status
	case is_nil(ret):
		return nil
	}
	return unmarshalData(cn.result, ret)
}

/* 
Simple internal test to determine if a nil value was passed
*/
func is_nil(v interface{}) bool {
	switch v.(type) {
	case nil:
		return true
	}
	return false
}

/* 
Send sends a command to the redis database, but does not wait for a response; this permits spamming a number of commands without waiting on a round trip.  Send will reset the internal Status, reflecting that redis tolerates sending new commands after a previous one has failed.

Send may produce an error if the underlying io.Writer fails, but is not aware of the redis database responce, since it does not wait for a response. 
*/
func (cn *Conn) Send(verb string, args ...interface{}) error {
	cn.status = nil
	err := writeRequest(cn.Rwc, verb, args...)
	if err != nil {
		return cn.updateStatus(err)
	}
	cn.lag++
	return nil
}

/* 
Next waits for a command to complete and returns an error if either the underlying Reader did, or if Redis indicated that the command did not complete successfully.  Next() will return the Status() if it was not nil, since probably will not be a response in the event of a Send() error. 
*/
func (cn *Conn) Next() error {
	if cn.lag > 0 {
		cn.lag--
	}
	if cn.status != nil {
		return cn.status
	}

	t, result, err := readMsg(cn.r)
	if err != nil {
		return cn.updateStatus(err)
	}

	switch t {
	case '+', '$', '*', ':':
		cn.result = result
	case '-':
		if len(result) == 0 {
			return cn.updateStatus(PROTOCOL_ERROR)
		}
		return cn.updateStatus(Error(string(result[0])))
	default:
		return cn.updateStatus(PROTOCOL_ERROR)
	}

	return nil
}

/* 
Sync will Next() through pending redis responses, then returns the resulting Status().  This relies on an internal counter used incremented by Send() and decremented when a response is read.

This is intended for uses where a number of Send() invocations were used, possibly as part of a MULTI .. EXEC block.
*/
func (cn *Conn) Sync() error {
	for cn.lag > 0 {
		cn.Next()
	}
	return cn.status
}

/* 
Status reproduces the oldest error associated with a connection.  The internal field associated with Status is cleared when Send() is invoked. 
*/
func (cn *Conn) Status() error {
	return cn.status
}

/* 
If err is not nil, but cn.status is, update cn.status.
*/
func (cn *Conn) updateStatus(err error) error {
	if err != nil && cn.status == nil {
		cn.status = err
	}
	return err
}

/* 
Forwards to the wrapped io.ReadWriteCloser.
*/
func (cn *Conn) Close() error {
	return cn.Rwc.Close()
}
