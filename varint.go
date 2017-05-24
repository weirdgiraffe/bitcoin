//
// varint.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package bitcoin

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Varint uint64

func ReadVarint(r io.Reader, u *Varint) (err error) {
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
		*u = Varint(s)
	case 0xfe:
		var w uint32
		err = binary.Read(r, binary.LittleEndian, &w)
		if err != nil {
			return
		}
		*u = Varint(w)
	case 0xff:
		var dw uint64
		err = binary.Read(r, binary.LittleEndian, r)
		if err != nil {
			return
		}
		*u = Varint(dw)
	default:
		*u = Varint(b[0])
	}
	return
}

type ScriptInt struct {
	val int64
}

func ScriptIntFromSlice(b []byte) *ScriptInt {
	if len(b) > 8 {
		panic(fmt.Errorf("ScriptInt is bigger than int64"))
	}
	s := &ScriptInt{}
	var offt uint = 0
	for i := range b {
		s.val |= (int64(b[i]) << offt)
		offt += 8
	}
	return s
}

func (s ScriptInt) Int() int {
	return int(s.val)
}

func (s ScriptInt) Bytes() []byte {
	b := make([]byte, 8)
	var offt uint = 0
	for i := range b {
		b[i] = byte((s.val >> offt) & 0xff)
		offt += 8
		if b[i] == 0 {
			return b[:i]
		}
	}
	return []byte{}
}
