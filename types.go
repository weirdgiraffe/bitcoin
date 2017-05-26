//
// types.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package bitcoin

import (
	"crypto/sha256"
	"encoding/hex"
)

type UnixTime uint32
type Hash [32]byte
type DoubleHash [32]byte

func (h Hash) MarshalJSON() ([]byte, error) {
	s := hex.EncodeToString(h[:])
	return []byte("\"" + s + "\""), nil
}

func (h DoubleHash) MarshalJSON() ([]byte, error) {
	s := hex.EncodeToString(h[:])
	return []byte("\"" + s + "\""), nil
}

func (h *DoubleHash) Update(b []byte) {
	h1 := sha256.Sum256(b)
	*h = sha256.Sum256(h1[:])
}
