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
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

type BlockHeader struct {
	Version    uint32
	PrevBlock  DoubleHash
	MerkleRoot DoubleHash
	UnixTime   UnixTime
	Bits       uint32
	Nonce      uint32
}

const BlockHeaderSize = 80

var BadMagic = errors.New("Bad block magic number")

type Block struct {
	Header BlockHeader
	Hash   DoubleHash
	tx     []*Tx
}

func (b Block) String() string {
	ob, err := json.MarshalIndent(&b, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(ob)
}

func (b *Block) TxCount() int {
	return len(b.tx)
}

func (b *Block) Tx(indx int) *Tx {
	return b.tx[indx]
}

type BlockFile struct {
	block []int64
	f     *os.File
}

func OpenBlockFile(path string) (b *BlockFile, err error) {
	b = new(BlockFile)
	b.f, err = os.Open(path)
	if err != nil {
		return nil, err
	}
	err = b.indexBlocks()
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (b *BlockFile) Close() {
	if b.f != nil {
		err := b.f.Close()
		if err != nil {
			log.Printf("Error closing BlockFile: %v", err)
		}
	}
}

func (b *BlockFile) BlockCount() int {
	return len(b.block)
}

func (b *BlockFile) Block(index int) (*Block, error) {
	if index >= 0 && index < len(b.block) {
		ret, err := b.readBlock(b.block[index])
		if err != nil {
			return nil, err
		}
		return ret, nil
	}
	return nil, fmt.Errorf("Bad block index %d (from %d blocks)", index, len(b.block))
}

func (b *BlockFile) indexBlocks() (err error) {
	var offt int64
	for {
		err = checkMagic(b.f)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return
		}
		b.block = append(b.block, offt)
		var bLen uint32
		err = binary.Read(b.f, binary.LittleEndian, &bLen)
		if err != nil {
			return
		}
		offt, err = b.f.Seek(int64(bLen), os.SEEK_CUR)
		if err != nil {
			return
		}
	}
}

func (b *BlockFile) readBlock(offt int64) (ret *Block, err error) {
	// magic number and length are already checked in indexBlocks()
	_, err = b.f.Seek(offt+4+4, os.SEEK_SET)
	if err != nil {
		return
	}
	ret = &Block{}
	buf := make([]byte, BlockHeaderSize)
	_, err = b.f.Read(buf)
	if err != nil {
		return
	}
	ret.Hash.Update(buf)
	_, err = b.f.Seek(offt+4+4, os.SEEK_SET)
	if err != nil {
		return
	}
	err = binary.Read(b.f, binary.LittleEndian, &ret.Header)
	if err != nil {
		return
	}
	var txCount Varint
	err = ReadVarint(b.f, &txCount)
	if err != nil {
		return
	}
	for i := Varint(txCount); i > 0; i-- {
		var t *Tx
		t, err = ReadTx(b.f)
		if err != nil {
			return
		}
		t.Block = ret.Hash
		ret.tx = append(ret.tx, t)
	}
	return ret, nil
}

func checkMagic(r io.Reader) (err error) {
	var b [4]byte
	_, err = r.Read(b[:])
	if err != nil {
		return
	}
	if bytes.Compare([]byte{0xf9, 0xbe, 0xb4, 0xd9}, b[:]) != 0 {
		return BadMagic
	}
	return nil
}
