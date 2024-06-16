package tcp

import (
	"fmt"
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

	state     network.Protocol
	emitEvent func(any)
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
			logger.Warn("Error reading connection")
			c.crash = err
			return
		}

		if len(pkBytes) > 0 {
			select {
			case <-c.quit:
				logger.Debug("Quitting handleIncomming due to c.quit closed.")
				return
			default: // Nothing
			}

			packets := network.HandlePacket(c.state, c.conn, &pkBytes)
			if packets != nil {
				for _, packet := range packets {
					if n, ok := packet.(network.OutgoingPacket); ok {
						c.outgoing <- n.Bytes()
					}

					c.nextState(packet)
					c.processEvent(packet)
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
			logger.Debug("Quittin handleOutgoing due to c.quit closed.")
			return
		case msg, ok := <-c.outgoing:
			if !ok {
				logger.Debug("Quitting handleOutgoing due to !ok")
				return
			}
			if err := network.WriteTo(c.conn, msg, c.compression); err != nil {
				logger.Debug("Quitting handleOutgoing due to WriteTo error")
				c.crash = err
				return
			}

		}
	}
}

func (c *TcpClient) Close() {
	c.once.Do(func() {
		if c.crash != nil {
			c.emitEvent(ClientCrashEvent{c, c.crash})
		} else {
			c.emitEvent(ClientDisconnectEvent{c})
		}

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

func (c *TcpClient) processEvent(packet network.Packet) {
	if n, ok := packet.(s2c.LoginSuccess); ok {
		event := BroadcastPacketEvent{
			packet: s2c.SystemChatMessage{
				Content: fmt.Sprintf("[+] %s joined the server.", n.Username),
				Overlay: false,
			}.Bytes(),
		}

		c.emitEvent(event)
	} else if n, ok := packet.(s2c.BlockUpdate); ok {
		event := BroadcastPacketEvent{
			packet: n.Bytes(),
		}

		c.emitEvent(event)
	}
}
