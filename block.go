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
type DoubleHash [32]byte

type Block struct {
	Header  BlockHeader
	TxCount Varint
}

type BlockHeader struct {
	Version    uint32
	PrevBlock  Hash
	MerkleRoot Hash
	UnixTime   UnixTime
	Bits       uint32
	Nonce      uint32
}

type TxIn struct {
	PrevTx        DoubleHash
	PrevTxOutIndx uint32
	ScriptLen     Varint
	Script        []byte
	SequenceNum   uint32
}

type TxOut struct {
	Value     [8]byte
	ScriptLen Varint
	Script    []byte
}

type Tx struct {
	Verssion uint32
	CountIn  Varint
	In       []TxIn
	CountOut Varint
	Out      []TxOut
	LockTime uint32
}

// IsCoinbase return true if this transaction is a generation transaction
// i.e. input for this transaction is a new generated block
func (t Tx) IsCoinbase() bool {
	null := DoubleHash{}
	return len(t.In) == 1 && bytes.Compare(null[:], t.In[0].PrevTx[:]) == 0
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
