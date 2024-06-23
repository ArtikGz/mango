package tcp

import (
	"io"
	"mango/src/config"
	"mango/src/logger"
	"mango/src/network"
	dt "mango/src/network/datatypes"
	"mango/src/network/packet/c2s"
	"mango/src/network/packet/s2c"
	"net"
	"sync"
	"time"
)

type TcpClient struct {
	conn        *net.TCPConn
	compression int

	incoming chan []byte
	outgoing chan network.OutgoingPacket

	crash error

	wg   sync.WaitGroup
	once sync.Once
	quit chan struct{}

	state     network.Protocol
	emitEvent func(any)

	username string
}

func (c *TcpClient) Username() string           { return c.username }
func (c *TcpClient) Protocol() network.Protocol { return c.state }

func NewTcpClient(conn *net.TCPConn) *TcpClient {
	c := &TcpClient{
		conn:        conn,
		compression: -1,
		incoming:    make(chan []byte, 128),
		outgoing:    make(chan network.OutgoingPacket, 128),
		quit:        make(chan struct{}),
		state:       network.SHAKE,
	}

	c.wg.Add(2)
	go c.handleIncoming()
	go c.handleOutgoing()

	return c
}

func (c *TcpClient) handleIncoming() {
	logger.Debug("handleIncoming: Starting")
	defer func() {
		logger.Debug("handleIncoming: Quitting")
		c.wg.Done()
		go c.Close()
	}()

	for {
		if err := c.conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
			logger.Error("Error setting the connection deadline: %s", err.Error())
			c.crash = err
			return
		}

		pkBytes, err := network.ReadFrom(c.conn, c.compression)
		if err != nil {
			if err != io.EOF {
				logger.Error("Error reading from client: %s", err.Error())
				c.crash = err
			}
			return
		}

		if len(pkBytes) > 0 {
			select {
			case <-c.quit:
				logger.Debug("Quitting handleIncoming due to c.quit closed.")
				return
			default: // Nothing
			}

			packets, err := network.HandlePacket(c, pkBytes)
			if err != nil {
				logger.Error("Error handling packet from client: %s", err.Error())
				c.crash = err
				return
			}

			for _, packet := range packets {
				if outPk, ok := packet.(network.OutgoingPacket); ok {
					if outPk.Broadcast() {
						c.emitEvent(BroadcastPacketEvent{outPk})
					} else {
						c.outgoing <- outPk
					}
					if lp, ok := packet.(s2c.LoginSuccess); ok {
						c.username = string(lp.Username)
					}
				}
				c.nextState(packet)
			}
		}
	}
}

func (c *TcpClient) handleOutgoing() {
	logger.Debug("handleOutgoing: Starting")
	ticker := time.NewTicker(10 * time.Second)
	var keepAlivePacket s2c.KeepAlive

	defer func() {
		logger.Debug("handleOutgoing: Quitting")
		ticker.Stop()
		c.wg.Done()
		go c.Close()
	}()

	for {
		select {
		case <-c.quit:
			logger.Debug("Quitting handleOutgoing due to c.quit closed.")
			return
		case timestamp := <-ticker.C:
			if c.state == network.PLAY {
				keepAlivePacket.KeepAliveID = dt.Long(timestamp.UTC().UnixNano())
				c.outgoing <- keepAlivePacket
			}
		case msg, ok := <-c.outgoing:
			if !ok {
				logger.Debug("Quitting handleOutgoing due to !ok")
				return
			}
			if err := network.WriteTo(c.conn, msg.Bytes(), c.compression); err != nil {
				logger.Debug("Quitting handleOutgoing due to WriteTo error")
				c.crash = err
				return
			}
			if pk, ok := msg.(s2c.SetCompressionPacket); ok {
				logger.Info("Client %s setting compression to %v", c.conn.RemoteAddr().String(), config.CompressionThreshold())
				c.compression = int(pk.Threshold)
				time.Sleep(50 * time.Millisecond) // sleep 1 tick to ensure the client applies the compression
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
