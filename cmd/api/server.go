package api

import (
	"fmt"
	"net"
)

type Server struct {
	Protocol string
	Port     string
	Socket   *net.Conn
}

func NewServer(protocol string, port string) *Server {
	server := Server{
		Protocol: protocol,
		Port:     port,
	}

	return &server
}

func (server *Server) Start() error {

	// server listen
	listener, err := net.Listen("tcp", server.Port)
	if err != nil {
		return nil
	}

	for {
		// accept connection from client
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		// declare channel
		ch := make(chan []byte)
		errCh := make(chan error)

		// read request from client
		go processRequest(conn, ch, errCh)

		// wait process and send response
		go sendResponse(conn, ch, errCh)
	}
}

func processRequest(conn net.Conn, ch chan []byte, errCh chan error) {

	buff := make([]byte, 0)

	for {
		tmpBuff := make([]byte, 0)

		size, err := conn.Read(tmpBuff)
		if err != nil {
			errCh <- err
			return
		}

		if size > 0 {
			buff = append(buff, tmpBuff...)
		}

		if len(buff) == 1440 {
			ch <- processFiles(&buff)
			return
		}
	}
}

func processFiles(buff *[]byte) []byte {
	*buff = make([]byte, 0)
	return []byte("ok")
}

func sendResponse(conn net.Conn, ch chan []byte, errCh chan error) {

	select {
	case res := <-ch:
		fmt.Println(res)
		conn.Write(res)

	case err := <-errCh:
		fmt.Println(err)
		conn.Write([]byte(err.Error()))
	}
}
