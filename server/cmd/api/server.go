package api

import (
	"fmt"
	"net"
)

const PACKET_SIZE int = 1440

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

	listener, err := net.Listen(server.Protocol, server.Port)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		ch := make(chan []byte)
		errCh := make(chan error)

		go processRequest(conn, ch, errCh)

		go sendResponse(conn, ch, errCh)
	}
}

func processRequest(conn net.Conn, ch chan []byte, errCh chan error) {

	buff := make([]byte, 0)
	totalSize := 0

	for {
		tmpBuff := make([]byte, PACKET_SIZE)

		size, err := conn.Read(tmpBuff)
		if err != nil {
			errCh <- err
			return
		}

		if size > 0 {
			buff = append(buff, tmpBuff...)
		}

		if len(buff) > 11 {
			if totalSize == 0 {
				totalSize = getTotalSize(buff)
			}

			if totalSize != 0 && totalSize > PACKET_SIZE {
				fmt.Printf("buff length %d, totalSize %d\n", len(buff), totalSize)
				if len(buff) > totalSize {
					res, err := processFiles(&buff)
					if err != nil {
						fmt.Println("processError: ", err)
						return
					}

					ch <- res
					return
				}
			} else if totalSize != 0 && totalSize < PACKET_SIZE {
				res, err := processFiles(&buff)
				if err != nil {
					fmt.Println("processError: ", err)
					return
				}
				ch <- res
				return
			}
		}
	}
}

func processFiles(buff *[]byte) ([]byte, error) {
	err := SaveFile(*buff)
	if err != nil {
		return nil, err
	}

	*buff = make([]byte, 0)
	return []byte("ok"), nil
}

func sendResponse(conn net.Conn, ch chan []byte, errCh chan error) {

	select {
	case res := <-ch:
		conn.Write(res)

	case err := <-errCh:
		fmt.Println(err)
		conn.Write([]byte(err.Error()))
	}
}
