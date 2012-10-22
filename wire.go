package mindis

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
)

/* reads a reply from the redis server, accepting single line replies, error messages, integers, bulk and multi-bulk replies */

func readMsg(br *bufio.Reader) (byte, [][]byte, error) {
	var b []byte
	var data [][]byte
	flag, err := br.ReadByte()

	switch flag {
	case '+', '-', ':':
		b, err = readLine(br)
		data = [][]byte{b}
	case '$':
		b, err = readBulk(br)
		data = [][]byte{b}
	case '*':
		data, err = readMultiBulk(br)
	default:
		err = PROTOCOL_ERROR
	}

	return flag, data, err
}

/* reads a multibulk reply, after the required '*' was read */

func readMultiBulk(br *bufio.Reader) ([][]byte, error) {
	count, err := readInt(br)
	switch {
	case err != nil:
		return nil, err
	case count < 0:
		return nil, PROTOCOL_ERROR
	}

	data := make([][]byte, 0, count)
	for count > 0 {
		f, err := br.ReadByte()
		if f != '$' {
			return nil, PROTOCOL_ERROR
		}
		b, err := readBulk(br)
		if err != nil {
			return nil, err
		}
		data = append(data, b)
		count--
	}

	return data, nil
}

/* reads a bulk reply from br, after the required '$' was read */

func readBulk(br *bufio.Reader) ([]byte, error) {
	size, err := readInt(br)
	if err != nil {
		return nil, err
	}

	data := make([]byte, size+2)
	_, err = br.Read(data)
	switch {
	case err != nil:
		return nil, err
	case data[size] != '\r':
		return nil, PROTOCOL_ERROR
	case data[size+1] != '\n':
		return nil, PROTOCOL_ERROR
	}

	return data[:size], nil
}

/* reads a integer line from br */

func readInt(br *bufio.Reader) (int, error) {
	line, err := readLine(br)
	if err != nil {
		return -1, err
	}
	return strconv.Atoi(string(line))
}

/* reads a crlf terminated line from br, tolerating embedded lf's which should never happen but exhaustively
   analyzing redis seemed harder than just coding defensively; strips the trailing crlf off  */

func readLine(br *bufio.Reader) ([]byte, error) {
	line := make([]byte, 0, 64)

	for {
		part, err := br.ReadSlice('\n')
		if err != nil {
			return nil, err
		}
		if bytes.HasSuffix(part, crlf) {
			line = append(line, part[:len(part)-2]...)
			break
		}
		line = append(line, part...)
	}

	return line, nil
}

/* 
When dealing with a network stack, you want to write full messages whenever possible; therefore, we take the extra copy costs of composing a complete message
*/

func writeRequest(w io.Writer, verb string, args ...interface{}) error {
	b, err := marshalRequest(verb, args...)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}
