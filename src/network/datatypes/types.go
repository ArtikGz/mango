package datatypes

import (
	"encoding/binary"
	"errors"
	"io"
	"math"

	"mango/src/nbt"
)

func ReadByte(reader io.Reader) (value byte, err error) {
	err = binary.Read(reader, binary.BigEndian, &value)
	return
}

type Byte byte // =====================================================

func (b *Byte) ReadFrom(reader io.Reader) (n int64, err error) {
	val, err := ReadByte(reader)
	*b = Byte(val)
	return
}

func (b *Byte) Bytes() (buffer []byte) {
	buffer = make([]byte, 1)
	buffer[0] = byte(*b)

	return buffer
}

type UByte uint8 // =====================================================

func (b *UByte) ReadFrom(reader io.Reader) (n int64, err error) {
	val, err := ReadByte(reader)
	*b = UByte(val)
	return
}

func (b *UByte) Bytes() (buffer []byte) {
	buffer = make([]byte, 1)
	buffer[0] = byte(*b)

	return buffer
}

type Short int16 // ======================================================

func (s *Short) ReadFrom(reader io.Reader) (n int64, err error) {
	err = binary.Read(reader, binary.BigEndian, &s)
	return
}

func (s *Short) Bytes() (buffer []byte) {
	buffer = make([]byte, 2)
	binary.BigEndian.PutUint16(buffer, uint16(*s))
	return buffer
}

type UShort uint16 // ====================================================

func (s *UShort) ReadFrom(reader io.Reader) (n int64, err error) {
	err = binary.Read(reader, binary.BigEndian, s)
	return
}

func (s *UShort) Bytes() (buffer []byte) {
	buffer = make([]byte, 2)
	binary.BigEndian.PutUint16(buffer, uint16(*s))
	return buffer
}

type Int int32 // ====================================================

func (s *Int) ReadFrom(reader io.Reader) (n int64, err error) {
	err = binary.Read(reader, binary.BigEndian, s)
	return
}

func (s *Int) Bytes() (buffer []byte) {
	buffer = make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, uint32(*s))
	return buffer
}

type UInt uint32 // ====================================================

func (i *UInt) ReadFrom(reader io.Reader) (n int64, err error) {
	err = binary.Read(reader, binary.BigEndian, i)
	return
}

func (i *UInt) Bytes() (buffer []byte) {
	buffer = make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, uint32(*i))
	return buffer
}

type String string // ====================================================

func (s *String) ReadFrom(reader io.Reader) (n int64, err error) {
	var length VarInt
	nn, err := length.ReadFrom(reader)
	if err != nil {
		return nn, err
	}
	n += nn

	stringBytes := make([]byte, length)
	_, err = reader.Read(stringBytes)
	if err != nil {
		return n, err
	}
	n += int64(length)

	*s = String(stringBytes)
	return
}

func (s String) Bytes() (buffer []byte) {

	strBytes := []byte(s)
	length := VarInt(len(strBytes))

	buffer = append(buffer, length.Bytes()...)
	buffer = append(buffer, strBytes...)

	return buffer
}

type VarInt int32 // =====================================================

func ReadVarInt(r io.Reader) (VarInt, int64, error) {
	var vi VarInt
	n, err := vi.ReadFrom(r)
	if err != nil {
		return 0, n, err
	}

	return vi, n, nil
}

func (vi VarInt) Length() int {
	if vi == 0 {
		return 1
	}
	i := 0
	for vi > 0 {
		vi >>= 7
		i++
	}
	return i
}

func (vi *VarInt) ReadFrom(reader io.Reader) (n int64, err error) {
	var value uint32

	for curr := byte(0x80); curr&0x80 != 0; n++ {
		if n > 5 {
			return n, errors.New("VarInt too big")
		}

		curr, err = ReadByte(reader)
		if err != nil {
			return n, err
		}

		value |= uint32(curr&0x7F) << (7 * n)
	}

	*vi = VarInt(value)
	return
}

func (vi VarInt) Bytes() (buffer []byte) {
	value := vi

	for i := 0; i < 5; i++ {
		var current byte = byte(value & 0x7F)
		value >>= 7

		if value > 0 {
			current |= 0x80
		}

		buffer = append(buffer, current)

		if value == 0 {
			return buffer
		}

	}
	return
}

type Float float32 // =====================================================

func (f *Float) ReadFrom(reader io.Reader) (n int64, err error) {
	var bits uint32
	err = binary.Read(reader, binary.BigEndian, &bits)
	*f = Float(math.Float32frombits(bits))
	return
}

func (f *Float) Bytes() (buffer []byte) {
	buffer = make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, math.Float32bits(float32(*f)))
	return buffer
}

type Double float64 // =====================================================

func (d *Double) ReadFrom(reader io.Reader) (n int64, err error) {
	var bits uint64
	err = binary.Read(reader, binary.BigEndian, &bits)
	*d = Double(math.Float64frombits(bits))
	return
}

func (d *Double) Bytes() (buffer []byte) {
	buffer = make([]byte, 8)
	binary.BigEndian.PutUint64(buffer, math.Float64bits(float64(*d)))
	return buffer
}

type Long uint64 // =====================================================

func (l *Long) ReadFrom(reader io.Reader) (n int64, err error) {
	err = binary.Read(reader, binary.BigEndian, l)
	return
}

func (l Long) Bytes() (buffer []byte) {
	buffer = make([]byte, 8)
	binary.BigEndian.PutUint64(buffer, uint64(l))
	return buffer
}

type Boolean bool // =====================================================

func (b *Boolean) ReadFrom(reader io.Reader) (n int64, err error) {
	val, err := ReadByte(reader)
	*b = val != 0
	return
}

func (b Boolean) Bytes() (buffer []byte) {
	buffer = make([]byte, 1)
	if b {
		buffer[0] = 1
	} else {
		buffer[0] = 0
	}

	return buffer
}

type NbtCompound nbt.NBTTag // Should have NBTType as a compound ========

func (nc *NbtCompound) ReadFrom(reader io.Reader) (n int64, err error) {
	// TODO
	return
}

func (nc *NbtCompound) Bytes() (buffer []byte) {
	buffer = nbt.Marshal(nbt.NBTTag(*nc))
	return
}

type Position struct{ X, Y, Z int } // =====================================================

func (p *Position) ReadFrom(reader io.Reader) (n int64, err error) {
	var value int64
	err = binary.Read(reader, binary.BigEndian, &value)

	p.X = int(value >> 38)
	p.Y = int(value << 52 >> 52)
	p.Z = int(value << 26 >> 38)

	return
}

func (p *Position) Bytes() (buffer []byte) {
	var value int64
	value |= int64((p.X & 0x3FFFFFF) << 38) // x = 26 MSBs
	value |= int64((p.Z & 0x3FFFFFF) << 12) // z = 26 middle bits
	value |= int64(p.Y & 0xFFF)             // y = 12 LSBs

	/*
		buffer = make([]byte, 8)
		binary.BigEndian.PutUint64(buffer, uint64(value))
	*/
	buffer = make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer, uint64(value))
	return
}

// =====================================================
type BitSet struct {
	Data []Long
}

func (bs *BitSet) Bytes() (buffer []byte) {
	length := VarInt(len(bs.Data))
	buffer = append(buffer, length.Bytes()...)

	data := make([]byte, 0, len(bs.Data)*8)
	for _, l := range bs.Data {
		data = append(data, l.Bytes()...)
	}

	buffer = append(buffer, data...)
	return
}

type Property struct {
	Name      String
	Value     String
	IsSigned  Boolean
	Signature String
}

func (p *Property) Bytes() (buffer []byte) {
	buffer = append(buffer, p.Name.Bytes()...)
	buffer = append(buffer, p.Value.Bytes()...)

	buffer = append(buffer, p.IsSigned.Bytes()...)
	if p.IsSigned {
		buffer = append(buffer, p.Signature.Bytes()...)
	}

	return
}
