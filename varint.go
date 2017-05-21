//
// varint.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package bitcoin

import (
	"encoding/binary"
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
