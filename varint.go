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

func (v Varint) OutSize() int {
	switch {
	case v < 0xfd:

		return 1
	case v < 0xffff:
		return 3
	case v < 0xffffffff:
		return 5
	default:
		return 9
	}
}

func WriteVarint(w io.Writer, u Varint) (err error) {
	switch {
	case u < 0xfd:
		v := byte(0xff & u)
		_, err = w.Write([]byte{v})
		if err != nil {
			return
		}
	case u < 0xffff:
		v := uint16(u)
		err = binary.Write(w, binary.LittleEndian, &v)
		if err != nil {
			return
		}
	case u < 0xffffffff:
		v := uint32(u)
		err = binary.Write(w, binary.LittleEndian, &v)
		if err != nil {
			return
		}
	default:
		v := uint64(u)
		err = binary.Write(w, binary.LittleEndian, &v)
		if err != nil {
			return
		}
	}
	return nil
}

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
		err = binary.Read(r, binary.LittleEndian, &dw)
		if err != nil {
			return
		}
		*u = Varint(dw)
	default:
		*u = Varint(b[0])
	}
	return nil
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
	if ((int64(0x80)<<offt)&s.val) > 0 && (s.val > 0) {
		s.val = -s.val
	}
	return s
}

func (s ScriptInt) Int64() int64 {
	return s.val
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
			if i == 0 {
				return []byte{0}
			}
			return b[:i]
		}
	}
	return b
}
