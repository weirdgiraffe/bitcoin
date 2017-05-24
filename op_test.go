//
// op_test.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package bitcoin

import (
	"bytes"
	"fmt"
	"testing"
)

func hex2byte(hex string) []byte {
	b := make([]byte, len(hex)/2)
	for i := range b {
		fmt.Sscanf(hex[i*2:i*2+2], "%02x", &b[i])
	}
	return b
}

func StackWithValues(hex ...string) *stack {
	s := &stack{}
	s.Reset()
	for i := range hex {
		s.PushSlice(hex2byte(hex[i]))
	}
	return s
}

func compareStack(a, b *stack) int {
	na := a.ItemsCount().Int()
	nb := b.ItemsCount().Int()
	if na != nb {
		return na - nb
	}
	for i := 0; i < na; i++ {
		ia := a.item[i]
		ib := b.item[i]
		r := bytes.Compare(ia, ib)
		if r != 0 {
			return r
		}
	}
	return 0
}

func TestOpConstants(t *testing.T) {
	tt := []struct {
		opcode   byte
		script   string
		expected *stack
	}{
		{OP_0, "", StackWithValues("")},
		{0x01, "ab", StackWithValues("ab")},
		{OP_PUSHDATA1, "01ab", StackWithValues("ab")},
		{OP_PUSHDATA2, "0100ab", StackWithValues("ab")},
		{OP_PUSHDATA4, "01000000ab", StackWithValues("ab")},
		{OP_1NEGATE, "", StackWithValues("81")},
		{OP_1, "", StackWithValues("01")},
		{OP_16, "", StackWithValues("10")},
	}
	s := &stack{}
	for i := range tt {
		s.Reset()
		n, err := OpConstants(tt[i].opcode, hex2byte(tt[i].script), s)
		if err != nil {
			t.Errorf("%02x%s error: %v", tt[i].opcode, tt[i].script, err)
		}
		if n != len(tt[i].script)/2 {
			t.Errorf("%02x%s consumed len mismatch %d != %d", tt[i].opcode, tt[i].script, len(tt[i].script)/2, n)
		}
		if compareStack(tt[i].expected, s) != 0 {
			t.Errorf("%02x%s stack mismatch", tt[i].opcode, tt[i].script)
			t.Errorf("expected[\n%s]", tt[i].expected)
			t.Errorf("got[\n%s]", s)
		}
	}
}
