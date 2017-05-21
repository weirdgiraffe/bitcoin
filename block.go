//
// block.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package bitcoin

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

type UnixTime uint32
type Hash [32]byte

type BlockHeader struct {
	Version    uint32
	PrevBlock  Hash
	MerkleRoot Hash
	UnixTime   UnixTime
	Bits       uint32
	Nonce      uint32
}

type Block struct {
	Header  BlockHeader
	TxCount uint64
}

func LoadBlock(r io.Reader) (b *Block, err error) {
	err = verifyMagicNumber(r)
	if err != nil {
		return
	}

	var bLen uint32
	err = binary.Read(r, binary.LittleEndian, &bLen)
	if err != nil {
		return
	}

	b = &Block{}
	err = binary.Read(r, binary.LittleEndian, &b.Header)
	if err != nil {
		return nil, err
	}

	err = ReadVarint(r, &b.TxCount)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func verifyMagicNumber(r io.Reader) error {
	var magic = [4]byte{0xf9, 0xbe, 0xb4, 0xd9}
	var m [4]byte
	n, err := r.Read(m[:])
	if err != nil {
		return err
	}
	if n != 4 {
		return fmt.Errorf("failed to read magic")
	}
	if bytes.Compare(m[:], magic[:]) != 0 {
		return fmt.Errorf("magic not match %v != %v", magic, m)
	}
	return nil
}

func ReadVarint(r io.Reader, u *uint64) (err error) {
	b := []byte{0}
	_, err = r.Read(b)
	if err != nil {
		return
	}
	switch b[0] {
	case 0xfd:
		var s uint16
		err = binary.Read(r, binary.LittleEndian, &s)
		if err != nil {
			return
		}
		*u = uint64(s)
	case 0xfe:
		var w uint32
		err = binary.Read(r, binary.LittleEndian, &w)
		if err != nil {
			return
		}
		*u = uint64(w)
	case 0xff:
		err = binary.Read(r, binary.LittleEndian, r)
		if err != nil {
			return
		}
	default:
		*u = uint64(b[0])
	}
	return
}

func (b Block) String() string {
	ob, err := json.MarshalIndent(&b, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(ob)
}

func (h *Hash) MarshalJSON() ([]byte, error) {
	fmt.Println("here")
	s := "\""
	for _, b := range h {
		s += fmt.Sprintf("%02x", b)
	}
	s += "\""
	return []byte(s), nil

}
