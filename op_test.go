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
			t.Errorf("case #%d error: %v", i+1, err)
		}
		if n != len(tt[i].script)/2 {
			t.Errorf("case #%d consumed len mismatch %d != %d", i+1, len(tt[i].script)/2, n)
		}
		if compareStack(tt[i].expected, s) != 0 {
			t.Errorf("case #%d main stack mismatch", i+1)
			t.Errorf("expected[\n%s]", tt[i].expected)
			t.Errorf("got[\n%s]", s)
		}
	}
}

func TestOpAltStack(t *testing.T) {
	tt := []struct {
		opcode      byte
		iMain, iAlt *stack
		eMain, eAlt *stack
	}{
		{
			OP_TOALTSTACK,
			StackWithValues("aa"), StackWithValues(),
			StackWithValues(), StackWithValues("aa"),
		},
		{
			OP_FROMALTSTACK,
			StackWithValues(), StackWithValues("bb"),
			StackWithValues("bb"), StackWithValues(),
		},
	}
	for i := range tt {
		err := OpStack(tt[i].opcode, tt[i].iMain, tt[i].iAlt)
		if err != nil {
			t.Errorf("case #%d error: %v", i+1, err)
		}
		if compareStack(tt[i].iMain, tt[i].eMain) != 0 {
			t.Errorf("case #%d main stack mismatch", i+1)
			t.Errorf("expected[\n%s]", tt[i].eMain)
			t.Errorf("got[\n%s]", tt[i].iMain)
		}
		if compareStack(tt[i].iAlt, tt[i].eAlt) != 0 {
			t.Errorf("case #%d alt stack mismatch", i+1)
			t.Errorf("expected[\n%s]", tt[i].eAlt)
			t.Errorf("got[\n%s]", tt[i].iAlt)
		}
	}
}

func TestOpStack(t *testing.T) {
	tt := []struct {
		op           byte
		in, expected *stack
	}{
		{OP_IFDUP, StackWithValues("aa"), StackWithValues("aa", "aa")},
		{OP_IFDUP, StackWithValues("aa", "00"), StackWithValues("aa", "00")},
		{OP_DEPTH, StackWithValues("aa", "00"), StackWithValues("aa", "00", "02")},
		{OP_DROP, StackWithValues("01", "02"), StackWithValues("01")},
		{OP_DROP, StackWithValues("aa"), StackWithValues()},
		{OP_DUP, StackWithValues("ab"), StackWithValues("ab", "ab")},
		{OP_NIP, StackWithValues("01", "02"), StackWithValues("02")},
		{OP_OVER, StackWithValues("01", "02"), StackWithValues("01", "02", "01")},
		{OP_PICK, StackWithValues("01", "02", "03", "00"), StackWithValues("01", "02", "03", "03")},
		{OP_PICK, StackWithValues("01", "02", "03", "01"), StackWithValues("01", "02", "03", "02")},
		{OP_PICK, StackWithValues("01", "02", "03", "02"), StackWithValues("01", "02", "03", "01")},
		{OP_ROLL, StackWithValues("01", "02", "03", "01"), StackWithValues("01", "03", "02")},
		{OP_ROLL, StackWithValues("01", "02", "03", "02"), StackWithValues("02", "03", "01")},
		{OP_ROT, StackWithValues("01", "02", "03"), StackWithValues("02", "03", "01")},
		{OP_SWAP, StackWithValues("01", "02"), StackWithValues("02", "01")},
		{OP_TUCK, StackWithValues("01", "02"), StackWithValues("02", "01", "02")},
		{OP_2DROP, StackWithValues("01", "02"), StackWithValues()},
		{OP_2DUP, StackWithValues("01", "02"), StackWithValues("01", "02", "01", "02")},
		{OP_3DUP, StackWithValues("01", "02", "03"), StackWithValues("01", "02", "03", "01", "02", "03")},
		{OP_2OVER, StackWithValues("01", "02", "03", "04"), StackWithValues("01", "02", "03", "04", "01", "02")},
		{OP_2ROT, StackWithValues("01", "02", "03", "04", "05", "06"), StackWithValues("03", "04", "05", "06", "01", "02")},
		{OP_2SWAP, StackWithValues("01", "02", "03", "04"), StackWithValues("03", "04", "01", "02")},
	}
	alt := &stack{}
	for i := range tt {
		alt.Reset()
		err := OpStack(tt[i].op, tt[i].in, alt)
		if err != nil {
			t.Errorf("case #%d error: %v", i+1, err)
		}
		if compareStack(tt[i].in, tt[i].expected) != 0 {
			t.Errorf("case #%d main stack mismatch", i+1)
			t.Errorf("expected[\n%s]", tt[i].expected)
			t.Errorf("got[\n%s]", tt[i].in)
		}
		if alt.ItemsCount().Int() != 0 {
			t.Errorf("case #%d alt stack has items", i+1)
			t.Errorf("\n%s", alt)
		}
	}
}

func TestOpSplice(t *testing.T) {
	tt := []struct {
		op           byte
		in, expected *stack
	}{
		{OP_SIZE, StackWithValues("0102030405"), StackWithValues("0102030405", "05")},
	}
	alt := &stack{}
	for i := range tt {
		alt.Reset()
		err := OpSplice(tt[i].op, tt[i].in, alt)
		if err != nil {
			t.Fatalf("case #%d error: %v", i+1, err)
		}
		if compareStack(tt[i].in, tt[i].expected) != 0 {
			t.Errorf("case #%d main stack mismatch", i+1)
			t.Errorf("expected[\n%s]", tt[i].expected)
			t.Errorf("got[\n%s]", tt[i].in)
		}
		if alt.ItemsCount().Int() != 0 {
			t.Errorf("case #%d alt stack has items", i+1)
			t.Errorf("\n%s", alt)
		}
	}
}

func TestOpBitwise(t *testing.T) {
	tt := []struct {
		op           byte
		in, expected *stack
		expect_err   bool
	}{
		{OP_EQUAL, StackWithValues("0102", "0102"), StackWithValues("01"), false},
		{OP_EQUAL, StackWithValues("0102", "0103"), StackWithValues("00"), false},
		{OP_EQUALVERIFY, StackWithValues("0102", "0102"), StackWithValues(), false},
		{OP_EQUALVERIFY, StackWithValues("0102", "0103"), StackWithValues(), true},
	}
	alt := &stack{}
	for i := range tt {
		alt.Reset()
		err := OpBitwise(tt[i].op, tt[i].in, alt)
		if err != nil && tt[i].expect_err == false {
			t.Fatalf("case #%d error: %v", i+1, err)
		}
		if compareStack(tt[i].in, tt[i].expected) != 0 {
			t.Errorf("case #%d main stack mismatch", i+1)
			t.Errorf("expected[\n%s]", tt[i].expected)
			t.Errorf("got[\n%s]", tt[i].in)
		}
		if alt.ItemsCount().Int() != 0 {
			t.Errorf("case #%d alt stack has items", i+1)
			t.Errorf("\n%s", alt)
		}
	}
}
