package api

import "net"

type Server struct {
	Address string
	Port    string
	Socket  *net.Conn
}

func NewServer(address string, port string) *Server {
	server := Server{
		Address: address,
		Port:    port,
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
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		// declare channel
		ch := make(chan []byte)
		errCh := make(chan error)

		// read request from client
		go readRequest(conn, ch, errCh)

		// wait process and send response
		go sendResponse(ch, errCh)
	}
}

func readRequest(conn net.Conn, ch chan []byte, errCh chan error) {

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
			ch <- processFiles(buff)
			buff = []byte{}
			return
		}
	}
}

func processFiles(buff []byte) []byte {
	return buff
}

func sendResponse(ch chan []byte, errCh chan error) {

	for {

	}
}
