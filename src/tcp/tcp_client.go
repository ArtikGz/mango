package tcp

import (
	"mango/src/logger"
	"mango/src/network"
	"mango/src/network/packet/c2s"
	"mango/src/network/packet/s2c"
	"net"
	"sync"
)

type TcpClient struct {
	conn        *net.TCPConn
	compression int

	incoming chan []byte
	outgoing chan []byte

	crash error

	wg   sync.WaitGroup
	once sync.Once
	quit chan struct{}

	state network.Protocol
}

func NewTcpClient(conn *net.TCPConn) *TcpClient {
	c := &TcpClient{
		conn:        conn,
		compression: -1,
		incoming:    make(chan []byte, 128),
		outgoing:    make(chan []byte, 128),
		quit:        make(chan struct{}),
		state:       network.SHAKE,
	}

	c.wg.Add(2)
	go c.handleIncoming()
	go c.handleOutgoing()

	return c
}

func (c *TcpClient) handleIncoming() {
	logger.Info("handleIncoming: Starting")
	defer func() {
		logger.Info("handleIncoming: Quitting")
		c.wg.Done()
		go c.Close()
	}()

	for {
		pkBytes, err := network.ReadFrom(c.conn, c.compression)
		if err != nil {
			c.crash = err
			return
		}

		packets := network.HandlePacket(c.state, c.conn, &pkBytes)
		if packets != nil {
			for _, packet := range packets {
				c.nextState(packet)

				if n, ok := packet.(network.OutgoingPacket); ok {
					c.outgoing <- n.Bytes()
				}
			}
		}
	}
}

func (c *TcpClient) handleOutgoing() {
	logger.Info("handleOutgoing: Starting")
	defer func() {
		logger.Info("handleOutgoing: Quitting")
		c.wg.Done()
		go c.Close()
	}()

	for {
		select {
		case <-c.quit:
			return
		case msg, ok := <-c.outgoing:
			if !ok {
				return
			}
			if err := network.WriteTo(c.conn, msg, c.compression); err != nil {
				c.crash = err
				return
			}

		}
	}
}

func (c *TcpClient) Close() {
	c.once.Do(func() {
		close(c.quit)
		close(c.incoming)
		close(c.outgoing)
		c.wg.Wait()
	})
}

func (c *TcpClient) nextState(packet network.Packet) {
	if n, ok := packet.(c2s.Handshake); ok {
		c.state = network.Protocol(n.NextState)
	} else if _, ok := packet.(s2c.LoginSuccess); ok {
		c.state = network.PLAY
	}
}
