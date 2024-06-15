package network

import (
	"net"
)

func ReadFrom(conn *net.TCPConn, compression int) ([]byte, error) {
	b := make([]byte, 1024*4)

	_, err := conn.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func WriteTo(conn *net.TCPConn, msg []byte, compression int) error {
	_, err := conn.Write(msg)
	if err != nil {
		return err
	}

	return nil
}
