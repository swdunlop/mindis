package main

import (
	".."
	"flag"
	"fmt"
	"net"
	"os"
)

var tcp_addr = flag.String("tcp", "", "specifies a tcp address and port for the server")
var unix_addr = flag.String("unix", "", "specifies a unix domain socket for the server")

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
	}

	r := mustDial()
	defer r.Close()

	args := make([]interface{}, len(flag.Args())-1)
	for i, a := range flag.Args()[1:] {
		args[i] = a
	}

	var ret []string
	if r.Exec(&ret, flag.Arg(0), args...) != nil {
		fmt.Fprintln(os.Stderr, "error:", r.Status())
		os.Exit(5)
	}

	for _, r := range ret {
		fmt.Println(r)
	}
}

func mustDial() *mindis.Conn {
	var nc net.Conn
	var err error

	switch {
	case tcp_addr != nil && *tcp_addr != "":
		nc, err = net.Dial("tcp", *tcp_addr)
	case unix_addr != nil && *unix_addr != "":
		nc, err = net.Dial("unix", *unix_addr)
	default:
		println("!! usage-error: you must specify either -tcp or -unix")
		os.Exit(2)
	}

	if err != nil {
		println("!! dial-error:", err.Error())
		os.Exit(2)
	}

	return mindis.Wrap(nc)
}

func usage() {
	flag.PrintDefaults()
	println("mindis-cli is a command line client for redis that submits one command at a time")
	os.Exit(1)
}
