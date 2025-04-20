package api

import (
	"fmt"
	"net"
	"net/http"
	"videoCompressorServer/utils"
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

		fCh := make(chan []byte)
		errCh := make(chan error)
		resCh := make(chan []byte)

		go processRequest(conn, fCh, errCh)

		go processFiles(fCh, resCh, errCh)

		go sendResponse(conn, resCh, errCh)
	}
}

func processRequest(conn net.Conn, fCh chan []byte, errCh chan error) {

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
			buff = append(buff, tmpBuff[:size]...)
		}

		if len(buff) > 11 {
			if totalSize == 0 {
				totalSize = getTotalSize(buff)
			}
			if totalSize != 0 && totalSize > PACKET_SIZE {
				if len(buff) == totalSize {
					fCh <- buff
					return
				}
			} else if totalSize != 0 && totalSize < PACKET_SIZE {
				fCh <- buff
				return
			}
		}
	}
}

func processFiles(fCh chan []byte, resCh chan []byte, errCh chan error) {
	for {
		file := <-fCh
		savedFile, fileJson, err := SaveFile(file)
		if err != nil {
			errCh <- err
			continue
		}

		formattedFile, err := FormatFile(savedFile, fileJson)
		if err != nil {
			errCh <- err
			continue
		}

		resCh <- utils.NewResponse(uint8(http.StatusOK), formattedFile)
	}
}

func sendResponse(conn net.Conn, resCh chan []byte, errCh chan error) {

	select {
	case res := <-resCh:
		_, err := conn.Write(res)
		if err != nil {
			fmt.Println(err)
		}

	case err := <-errCh:
		conn.Write(utils.NewResponse(uint8(0), nil, err))
		if err != nil {
			fmt.Println(err)
		}
	}
}
