package netString

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

const maxMsg uint64 = 1024 * 8

type NetStringConn struct {
	Conn     net.Conn
	MaxMsg   uint64
	queueOut chan string
}

type NetStringProcessor interface {
	Connected(*NetStringConn)
	Msg(*NetStringConn, string)
	Disconected(*NetStringConn, error)
}

type NetString struct {
	MaxMsg uint64
	nsp    NetStringProcessor
}

func NewNetString(nsp NetStringProcessor) *NetString {
	ns := NetString{maxMsg, nsp}
	return &ns
}

func (n *NetString) Connect(conn net.Conn) {
	n.process(conn)
}

func (n *NetString) Listen(ln net.Listener) error {
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go n.process(conn)
	}
}

func (n *NetString) process(conn net.Conn) {
	nsc := &NetStringConn{conn, n.MaxMsg, make(chan string, 10)}
	go nsc.process()
	n.nsp.Connected(nsc)
	reader := bufio.NewReader(conn)
	for {
		longs, err := reader.ReadString(':')
		if err != nil {
			n.nsp.Disconected(nsc, err)
			break
		}
		long, err := strconv.ParseUint(longs[0:len(longs)-1], 10, 64)
		if err != nil {
			n.nsp.Disconected(nsc, fmt.Errorf("error in convert len data: %s\n", err.Error()))
			break
		}
		fmt.Printf("max: %v\n", nsc.MaxMsg)
		if long > nsc.MaxMsg {
			n.nsp.Disconected(nsc, fmt.Errorf("Data len (%d) is bigger than max msg len (%d)\n", long, nsc.MaxMsg))
			break
		}
		b := make([]byte, long)
		if nread, err := reader.Read(b); err != nil {
			n.nsp.Disconected(nsc, fmt.Errorf("error in receive data: %s\n", err.Error()))
			break
		} else if nread != len(b) {
			n.nsp.Disconected(nsc, fmt.Errorf("error, receive %d data and expected %d\n", nread, long))
			break
		}
		if by, err := reader.ReadByte(); err != nil {
			n.nsp.Disconected(nsc, fmt.Errorf("error in last byte, must be ',': %s", err.Error()))
			break
		} else if by != ',' {
			n.nsp.Disconected(nsc, fmt.Errorf("error, last byte must be ',' and is '%v'", by))
			break
		}
		n.nsp.Msg(nsc, string(b))
	}
	close(nsc.queueOut)
}

func (n *NetStringConn) process() {
	w := bufio.NewWriter(n.Conn)
	for s := range n.queueOut {
		w.Write([]byte(strconv.Itoa(len(s))))
		w.WriteByte(':')
		w.WriteString(s)
		w.WriteByte(',')
		w.Flush()
	}
}

func (n *NetStringConn) Send(s string) {
	n.queueOut <- s
}
