package netString

import (
	"fmt"
	"net"
	"testing"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Server struct {
	ln net.Listener
}
type Client struct{}

func (s *Server) Connected(nsc *NetStringConn) {
	fmt.Printf("Sconnected: %v\n", nsc.Conn.RemoteAddr())
}
func (s *Server) Msg(nsc *NetStringConn, str string) {
	fmt.Printf("Sreceive %s\n", str)
	nsc.Send(str)
}
func (s *Server) Disconected(nsc *NetStringConn, err error) {
	fmt.Printf("SDisconected: %v %s\n", nsc.Conn.RemoteAddr(), err.Error())
	s.ln.Close()
}

func (c *Client) Connected(nsc *NetStringConn) {
	fmt.Printf("connected: %v\n", nsc.Conn.RemoteAddr())
	nsc.Send("test")
}
func (c *Client) Msg(nsc *NetStringConn, s string) {
	fmt.Printf("receive %s\n", s)
	nsc.Conn.Close()
}
func (c *Client) Disconected(nsc *NetStringConn, err error) {
	fmt.Printf("Disconected: %v %s\n", nsc.Conn.RemoteAddr(), err.Error())
}

func TestMain(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:18000")
	check(err)
	ns := newNetString(&Server{ln})
	go func() {
		conn, err := net.Dial("tcp", "127.0.0.1:18000")
		check(err)
		nsc := newNetString(&Client{})
		go nsc.Connect(conn)
		fmt.Println("client exit")
	}()
	ns.Listen(ln)
}
