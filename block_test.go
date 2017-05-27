//
// block_test.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package bitcoin

import (
	"testing"
)

func TestBlockFile(t *testing.T) {
	bf, err := OpenBlockFile("assets/testblock.dat")
	if err != nil {
		t.Fatal(err)
	}
	defer bf.Close()
	_, err = bf.Block(0)
	if err != nil {
		t.Error(err)
	}
}
