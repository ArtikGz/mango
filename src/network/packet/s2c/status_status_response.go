package s2c

import (
	"bytes"
	dt "mango/src/network/datatypes"
)

type StatusResponse struct {
	JsonPayload dt.String
	StatusData  StatusData
}

type StatusData struct {
	Protocol uint16
}

func (pk StatusResponse) getStatusPayload() string {
	return dt.GetDemoServerStatus(int(pk.StatusData.Protocol))
}

func (pk StatusResponse) Bytes() []byte {
	pk.JsonPayload = dt.String(pk.getStatusPayload())

	return bytes.Join([][]byte{
		dt.VarInt(0x00).Bytes(),
		pk.JsonPayload.Bytes(),
	}, nil)
}

func (pk StatusResponse) Broadcast() bool {
	return false
}
