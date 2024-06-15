package tcp

import (
	"fmt"
	"mango/src/logger"
	"net"
	"sync"
)

type ( // events
	ClientConnectEvent struct {
		conn *net.TCPConn
	}

	ClientDisconnectEvent struct {
		client *TcpClient
	}

	ClientCrashEvent struct {
		client *TcpClient
		err    error
	}
)

type TcpServer struct {
	clients map[*TcpClient]struct{}
	events  chan any

	listener *net.TCPListener

	quit chan struct{}
	wg   sync.WaitGroup
}

func NewTcpServer(host string, port int) (*TcpServer, error) {
	address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}

	l, err := net.ListenTCP("tcp", address)
	if err != nil {
		return nil, err
	}

	s := &TcpServer{
		clients: make(map[*TcpClient]struct{}),
		events:  make(chan any, 1024),

		listener: l,

		quit: make(chan struct{}),
	}

	return s, nil
}

func (s *TcpServer) Start() {
	s.wg.Add(2)
	go s.eventloop()
	go s.serve()

	s.wg.Wait()
}

// handles all client connection events
func (s *TcpServer) eventloop() {
	logger.Info("eventloop: Starting")
	defer s.wg.Done()
	defer logger.Info("eventloop: Quitting")

	for {
		select {
		case <-s.quit:
			return

		case event := <-s.events:
			switch event.(type) {
			case ClientConnectEvent:
				s.wg.Add(1)
				conn := event.(ClientConnectEvent).conn
				err1 := conn.SetKeepAlive(true)
				err2 := conn.SetNoDelay(true)
				client := NewTcpClient(conn)

				if err1 != nil {
					s.events <- ClientCrashEvent{client, err1}
				} else if err2 != nil {
					s.events <- ClientCrashEvent{client, err2}
				} else {
					s.clients[client] = struct{}{}
					logger.Info("Client %s connected", client.conn.RemoteAddr().String())
				}

			case ClientCrashEvent:
				defer s.wg.Done()
				client := event.(ClientCrashEvent).client
				err := event.(ClientCrashEvent).err
				if err != nil {
					logger.Info("Client %s crashed: %s", client.conn.RemoteAddr().String(), err.Error())
				} else {
					logger.Info("Client %s crashed with an unknown error", client.conn.RemoteAddr().String())
				}
				delete(s.clients, client)

				// handle closing in another thread
				go client.Close()

			case ClientDisconnectEvent:
				defer s.wg.Done()
				client := event.(ClientDisconnectEvent).client
				logger.Info("Client %s disconnected", client.conn.RemoteAddr().String())
				delete(s.clients, client)

				// handle closing in another thread
				go client.Close()
			}
		}
	}
}

func (s *TcpServer) serve() {
	logger.Info("serve: Starting")
	defer s.wg.Done()
	defer logger.Info("serve: Quitting")

	for {
		conn, err := s.listener.AcceptTCP()
		if err != nil {
			select {
			case <-s.quit:
				return
			default:
				logger.Warn("Error while accepting connection: %v", err.Error())
			}
			continue
		}
		s.events <- ClientConnectEvent{conn}
	}
}

func (s *TcpServer) Close() {
	for client := range s.clients {
		select {
		case <-client.quit:
			select {
			case s.events <- ClientCrashEvent{client, client.crash}:
				continue
			default:
				s.wg.Done()
				continue
			}
		default:
			select {
			case s.events <- ClientDisconnectEvent{client}:
				continue
			default:
				s.wg.Done()
				continue
			}
		}
	}

	close(s.quit)
	s.listener.Close()
	s.wg.Wait()
	close(s.events)
}
