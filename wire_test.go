package mindis

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

var flag byte
var msg [][]byte
var err error

func Test_Read(t *testing.T) {
	readMockMsg(t, "+OK")
	expectMsg(t, '+', "OK")
	readMockMsg(t, "-DIE IN A FIRE")
	expectMsg(t, '-', "DIE IN A FIRE")
	readMockMsg(t, ":1001")
	expectMsg(t, ':', "1001")
	readMockMsg(t, "$4", "helo")
	expectMsg(t, '$', "helo")
	readMockMsg(t, "*2", "$4", "helo", "$5", "world")
	expectMsg(t, '*', "helo", "world")

	readMockMsg(t, "*8", "$7", "RamDisk", "$7", "VNCPort", "$11", "sample-test", "$7", "GDBPort", "$12", "ebola_status", "$6", "Memory", "$4", "Core", "$16", "sandcastle_nodes")
	expectMsg(t, '*', "RamDisk", "VNCPort", "sample-test", "GDBPort", "ebola_status", "Memory", "Core", "sandcastle_nodes")
}

func readMockMsg(t *testing.T, lines ...string) {
	flag, msg, err = readMsg(bufio.NewReader(bytes.NewBufferString(strings.Join(lines, "\r\n") + "\r\n")))
}

func expectMsg(t *testing.T, f byte, values ...string) {
	if err != nil {
		t.Fatal("unexpected error:", err.Error())
	}
	if flag != f {
		t.Fatalf("expected flag %#v, got: %#v", f, flag)
	}
	msgstr := make([]string, len(msg))
	for i, b := range msg {
		msgstr[i] = string(b)
	}
	if strings.Join(msgstr, "\000") != strings.Join(values, "\000") {
		t.Fatalf("expected data %#v, got %#v")
	}
}
