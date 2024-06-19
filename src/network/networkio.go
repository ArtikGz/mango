package network

import (
	"bytes"
	"compress/zlib"
	"io"
	dt "mango/src/network/datatypes"
)

func ReadFrom(r io.Reader, compression int) ([]byte, error) {
	pkLen, _, err := dt.ReadVarInt(r)
	if err != nil {
		return nil, err
	}

	if pkLen == 0 {
		return []byte{}, nil
	}

	if compression == -1 {
		return readRaw(r, int64(pkLen))
	}

	return readCompressed(r, int64(pkLen))
}

func WriteTo(w io.Writer, pk []byte, compression int) error {
	if compression == -1 {
		return writeRaw(w, pk)
	}
	return writeCompressed(w, pk, compression)
}

// raw packet format: [length, pid, data...]
func readRaw(r io.Reader, length int64) ([]byte, error) {
	buf := make([]byte, length)

	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// compressed packet format:
// - [length, 0, pid, data...]  (if raw length < compression threshold)
// - [length, rawLength, zip[pid, data...]]
func readCompressed(r io.Reader, length int64) ([]byte, error) {
	rawLength, n, err := dt.ReadVarInt(r)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, length-n)
	if _, err = io.ReadFull(r, buf); err != nil {
		return nil, err
	}

	// decompress buffer if needed
	if rawLength != 0 {
		rawBuf := make([]byte, rawLength)
		br := bytes.NewReader(buf)
		zr, err := zlib.NewReader(br)
		if err != nil {
			return nil, err
		}

		if _, err = io.ReadFull(zr, rawBuf); err != nil {
			return nil, err
		}

		if err = zr.Close(); err != nil {
			return nil, err
		}

		buf = rawBuf
	}

	br := bytes.NewReader(buf)

	if rawLength != 0 {
		return readRaw(br, int64(rawLength))
	}

	return readRaw(br, length-n)
}

// raw packet format: [length, pid, data...]
func writeRaw(w io.Writer, pk []byte) error {
	pkLen := dt.VarInt(len(pk))
	pkLenLen := pkLen.Length()
	buf := make([]byte, int(pkLen)+pkLenLen)

	copy(buf[:pkLenLen], pkLen.Bytes())
	copy(buf[pkLenLen:pkLenLen+int(pkLen)], pk)

	_, err := w.Write(buf)
	return err
}

// compressed packet format:
// - [length, 0, pid, data...]  (if raw length < compression threshold)
// - [length, rawLength, zlib[pid, data...]]
func writeCompressed(w io.Writer, pk []byte, compression int) error {
	rawLen := dt.VarInt(len(pk))

	// write uncompressed packet
	if int(rawLen) < compression {
		pkLen := rawLen + 1
		pkLenLen := pkLen.Length()
		buf := make([]byte, pkLenLen+int(pkLen))

		copy(buf[:pkLenLen], pkLen.Bytes())
		buf[pkLenLen] = 0
		copy(buf[pkLenLen+1:pkLenLen+1+int(pkLen)], pk)

		_, err := w.Write(buf)
		return err
	}

	// write compressed packet to buffer
	compressedBuf := bytes.NewBuffer(make([]byte, len(pk)))
	zw := zlib.NewWriter(compressedBuf)

	if _, err := zw.Write(pk); err != nil {
		return err
	}

	if err := zw.Close(); err != nil {
		return err
	}

	// write buffer to writer
	rawLenLen := rawLen.Length()
	pkLen := dt.VarInt(rawLenLen + compressedBuf.Len())
	pkLenLen := pkLen.Length()

	buf := make([]byte, pkLenLen+int(pkLen))

	copy(buf[:pkLenLen], pkLen.Bytes())
	copy(buf[pkLenLen:pkLenLen+rawLenLen], rawLen.Bytes())
	copy(buf[pkLenLen+rawLenLen:], compressedBuf.Bytes())

	_, err := w.Write(buf)
	return err
}
