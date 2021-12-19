package network

const (
	PLAY int32 = iota - 1
	HANDSHAKE
	STATUS
	LOGIN
)

func Handshake(conn net.Conn) {
	var handshake c2s.Handshake
	handshake.ReadPacket(conn)
	fmt.Println(handshake)

	var request c2s.Request
	request.ReadPacket(conn)
	fmt.Println(request)

	// S2C_response

	var ping c2s.Ping
	ping.ReadPacket(conn)
	fmt.Println(ping)
}
