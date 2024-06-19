package c2s

import (
	"io"
	dt "mango/src/network/datatypes"
)

type ClientInformation struct {
	Locale              dt.String
	ViewDistance        dt.Byte
	ChatMode            dt.VarInt
	ChatColors          dt.Boolean
	DisplayedSkinsPart  dt.UByte
	MainHand            dt.VarInt
	EnableTextFiltering dt.Boolean
	AllowServerListing  dt.Boolean
}

func ReadClientInformationPacket(r io.Reader) (*ClientInformation, error) {
	var pk ClientInformation
	if _, err := pk.Locale.ReadFrom(r); err != nil {
		return nil, err
	}
	if _, err := pk.ViewDistance.ReadFrom(r); err != nil {
		return nil, err
	}
	if _, err := pk.ChatMode.ReadFrom(r); err != nil {
		return nil, err
	}
	if _, err := pk.ChatColors.ReadFrom(r); err != nil {
		return nil, err
	}
	if _, err := pk.DisplayedSkinsPart.ReadFrom(r); err != nil {
		return nil, err
	}
	if _, err := pk.MainHand.ReadFrom(r); err != nil {
		return nil, err
	}
	if _, err := pk.EnableTextFiltering.ReadFrom(r); err != nil {
		return nil, err
	}
	if _, err := pk.AllowServerListing.ReadFrom(r); err != nil {
		return nil, err
	}
	return &pk, nil
}
