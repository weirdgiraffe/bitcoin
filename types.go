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

func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

func (h Hash) MarshalJSON() ([]byte, error) {
	return []byte("\"" + h.String() + "\""), nil
}

func (h DoubleHash) MarshalJSON() ([]byte, error) {
	return []byte("\"" + h.String() + "\""), nil
}

func (h DoubleHash) String() string {
	return hex.EncodeToString(h[:])
}

func (h *DoubleHash) Update(b []byte) {
	h1 := sha256.Sum256(b)
	*h = sha256.Sum256(h1[:])
}
