//
// types.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package bitcoin

import (
	"encoding/hex"
)

type UnixTime uint32
type Hash [32]byte
type DoubleHash [32]byte

const DoubleHashSize = 32

func (h Hash) MarshalJSON() ([]byte, error) {
	s := hex.EncodeToString(h[:])
	return []byte("\"" + s + "\""), nil
}

func (h *DoubleHash) MarshalJSON() ([]byte, error) {
	s := hex.EncodeToString(h[:])
	return []byte("\"" + s + "\""), nil
}
