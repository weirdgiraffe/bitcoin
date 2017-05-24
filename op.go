//
// op.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package bitcoin

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

var InvalidTransaction = errors.New("Transacrion is invalid")

// OpConstants implements all script operations that are constants
// check https://en.bitcoin.it/wiki/Script#Constants
//
// return number of consumed bytes from script
func OpConstants(op byte, script []byte, s *stack) (n int, err error) {
	switch {
	case op == OP_0:
		s.PushSlice([]byte{})
	case op < OP_PUSHDATA1:
		n = int(op)
		s.PushSlice(script[:n])
	case op == OP_PUSHDATA1:
		n = int(script[0]) + 1
		s.PushSlice(script[1:n])
	case op == OP_PUSHDATA2:
		n = int(binary.LittleEndian.Uint16(script)) + 2
		s.PushSlice(script[2:n])
	case op == OP_PUSHDATA4:
		n = int(binary.LittleEndian.Uint32(script)) + 4
		s.PushSlice(script[4:n])
	case op == OP_1NEGATE:
		s.PushByte(0x81)
	case OP_1 <= op && op <= OP_16:
		b := op - OP_1 + 1
		s.PushByte(b)
	default:
		err = fmt.Errorf("0x%02x not a Script Constants op", op)
	}
	return
}

// OpStack implements all script operations that are stack
// check https://en.bitcoin.it/wiki/Script#Stack
func OpStack(op byte, main, alt *stack) error {
	switch op {
	case OP_TOALTSTACK:
		b1 := main.Pop()
		alt.PushSlice(b1)
	case OP_FROMALTSTACK:
		b1 := alt.Pop()
		main.PushSlice(b1)
	case OP_IFDUP:
		b1 := main.Top()
		if slice2bool(b1) {
			main.PushSlice(b1)
		}
	case OP_DEPTH:
		n := main.ItemsCount()
		main.PushSlice(n.Bytes())
	case OP_DROP:
		main.Pop()
	case OP_DUP:
		b1 := main.Top()
		main.PushSlice(b1)
	case OP_NIP:
		b2 := main.Pop()
		main.Pop()
		main.PushSlice(b2)
	case OP_OVER:
		b2 := main.Pop()
		b1 := main.Top()
		main.PushSlice(b2)
		main.PushSlice(b1)
	case OP_PICK:
		n := ScriptIntFromSlice(main.Pop()).Int()
		bn := main.Item(n)
		main.PushSlice(bn)
	case OP_ROLL:
		n := ScriptIntFromSlice(main.Pop()).Int()
		b := make([][]byte, n)
		for i := 0; i < n; i++ {
			b[i] = main.Pop()
		}
		bn := main.Pop()
		for i := n - 1; i >= 0; i-- {
			main.PushSlice(b[i])
		}
		main.PushSlice(bn)
	case OP_ROT:
		b3 := main.Pop()
		b2 := main.Pop()
		b1 := main.Pop()
		main.PushSlice(b2)
		main.PushSlice(b3)
		main.PushSlice(b1)
	case OP_SWAP:
		b2 := main.Pop()
		b1 := main.Pop()
		main.PushSlice(b2)
		main.PushSlice(b1)
	case OP_TUCK:
		b2 := main.Pop()
		b1 := main.Pop()
		main.PushSlice(b2)
		main.PushSlice(b1)
		main.PushSlice(b2)
	case OP_2DROP:
		main.Pop()
		main.Pop()
	case OP_2DUP:
		b2 := main.Pop()
		b1 := main.Pop()
		main.PushSlice(b1)
		main.PushSlice(b2)
		main.PushSlice(b1)
		main.PushSlice(b2)
	case OP_3DUP:
		b3 := main.Pop()
		b2 := main.Pop()
		b1 := main.Pop()
		main.PushSlice(b1)
		main.PushSlice(b2)
		main.PushSlice(b3)
		main.PushSlice(b1)
		main.PushSlice(b2)
		main.PushSlice(b3)
	case OP_2OVER:
		b2 := main.Item(2)
		b1 := main.Item(3)
		main.PushSlice(b1)
		main.PushSlice(b2)
	case OP_2ROT:
		b6 := main.Pop()
		b5 := main.Pop()
		b4 := main.Pop()
		b3 := main.Pop()
		b2 := main.Pop()
		b1 := main.Pop()
		main.PushSlice(b3)
		main.PushSlice(b4)
		main.PushSlice(b5)
		main.PushSlice(b6)
		main.PushSlice(b1)
		main.PushSlice(b2)
	case OP_2SWAP:
		b4 := main.Pop()
		b3 := main.Pop()
		b2 := main.Pop()
		b1 := main.Pop()
		main.PushSlice(b3)
		main.PushSlice(b4)
		main.PushSlice(b1)
		main.PushSlice(b2)
	default:
		return fmt.Errorf("0x%02x not a Script Stack op", op)
	}
	return nil
}

// OpSplice implements all script operations that are Splice
// check https://en.bitcoin.it/wiki/Script#Splice
func OpSplice(op byte, main, alt *stack) error {
	switch op {
	case OP_SIZE:
		b1 := main.Top()
		n := ScriptInt{int64(len(b1))}
		main.PushSlice(n.Bytes())
	default:
		return fmt.Errorf("0x%02x not a Script Splice op", op)
	}
	return nil
}

// OpBitwise implements all script operations that are Bitwise logic
// check https://en.bitcoin.it/wiki/Script#Bitwise_logic
func OpBitwise(op byte, main, alt *stack) error {
	switch op {
	case OP_EQUAL:
		b1 := main.Pop()
		b2 := main.Pop()
		if bytes.Compare(b1, b2) == 0 {
			main.PushByte(1)
		} else {
			main.PushByte(0)
		}
	case OP_EQUALVERIFY:
		b1 := main.Pop()
		b2 := main.Pop()
		if bytes.Compare(b1, b2) != 0 {
			return InvalidTransaction
		}
	default:
		return fmt.Errorf("0x%02x not a Script Bitwise logic op", op)
	}
	return nil
}

// OpArithmetic implements all script operations that are Arithmetic
// check https://en.bitcoin.it/wiki/Script#Arithmetic
func OpArithmetic(op byte, main, alt *stack) error {
	switch op {
	case OP_1ADD:
		n := ScriptIntFromSlice(main.Pop())
		n.val++
		main.PushSlice(n.Bytes())
	case OP_1SUB:
		n := ScriptIntFromSlice(main.Pop())
		n.val--
		main.PushSlice(n.Bytes())
	case OP_NEGATE:
		b := main.Pop()
		n := len(b)
		if n > 0 {
			if b[n-1]&0x80 > 0 {
				b[n-1] = b[n-1] & 0x7f
			} else {
				b[n-1] = b[n-1] | 0x80
			}
		}
		main.PushSlice(b)
	case OP_ABS:
		b := main.Pop()
		n := len(b)
		if n > 0 {
			b[n-1] = b[n-1] & 0x7f
		}
		main.PushSlice(b)
	case OP_NOT:
		n := ScriptIntFromSlice(main.Pop()).Int()
		if n == 0 {
			main.PushByte(1)
		} else {
			main.PushByte(0)
		}
	case OP_0NOTEQUAL:
		n := ScriptIntFromSlice(main.Pop()).Int()
		if n == 0 {
			main.PushByte(0)
		} else {
			main.PushByte(1)
		}
	case OP_ADD:
		b := ScriptIntFromSlice(main.Pop()).Int64()
		a := ScriptIntFromSlice(main.Pop()).Int64()
		c := ScriptInt{a + b}
		main.PushSlice(c.Bytes())
	case OP_SUB:
		b := ScriptIntFromSlice(main.Pop()).Int64()
		a := ScriptIntFromSlice(main.Pop()).Int64()
		c := ScriptInt{a - b}
		main.PushSlice(c.Bytes())
	case OP_BOOLAND:
		b := slice2bool(main.Pop())
		a := slice2bool(main.Pop())
		if a && b {
			main.PushByte(1)
		} else {
			main.PushByte(0)
		}
	case OP_BOOLOR:
		b := slice2bool(main.Pop())
		a := slice2bool(main.Pop())
		if a || b {
			main.PushByte(1)
		} else {
			main.PushByte(0)
		}
	case OP_NUMEQUAL:
		b := ScriptIntFromSlice(main.Pop()).Int64()
		a := ScriptIntFromSlice(main.Pop()).Int64()
		if a == b {
			main.PushByte(1)
		} else {
			main.PushByte(0)
		}
	case OP_NUMEQUALVERIFY:
		b := ScriptIntFromSlice(main.Pop()).Int64()
		a := ScriptIntFromSlice(main.Pop()).Int64()
		if a != b {
			return InvalidTransaction
		}
	case OP_NUMNOTEQUAL:
		b := ScriptIntFromSlice(main.Pop()).Int64()
		a := ScriptIntFromSlice(main.Pop()).Int64()
		if a != b {
			main.PushByte(1)
		} else {
			main.PushByte(0)
		}
	case OP_LESSTHAN:
		b := ScriptIntFromSlice(main.Pop()).Int64()
		a := ScriptIntFromSlice(main.Pop()).Int64()
		if a < b {
			main.PushByte(1)
		} else {
			main.PushByte(0)
		}
	case OP_GREATERTHAN:
		b := ScriptIntFromSlice(main.Pop()).Int64()
		a := ScriptIntFromSlice(main.Pop()).Int64()
		if a > b {
			main.PushByte(1)
		} else {
			main.PushByte(0)
		}
	case OP_LESSTHANOREQUAL:
		b := ScriptIntFromSlice(main.Pop()).Int64()
		a := ScriptIntFromSlice(main.Pop()).Int64()
		if a <= b {
			main.PushByte(1)
		} else {
			main.PushByte(0)
		}
	case OP_GREATERTHANOREQUAL:
		b := ScriptIntFromSlice(main.Pop()).Int64()
		a := ScriptIntFromSlice(main.Pop()).Int64()
		if a >= b {
			main.PushByte(1)
		} else {
			main.PushByte(0)
		}
	case OP_MIN:
		b := ScriptIntFromSlice(main.Pop())
		a := ScriptIntFromSlice(main.Pop())
		if a.Int64() < b.Int64() {
			main.PushSlice(a.Bytes())
		} else {
			main.PushSlice(b.Bytes())
		}
	case OP_MAX:
		b := ScriptIntFromSlice(main.Pop())
		a := ScriptIntFromSlice(main.Pop())
		if a.Int64() > b.Int64() {
			main.PushSlice(a.Bytes())
		} else {
			main.PushSlice(b.Bytes())
		}
	case OP_WITHIN:
		max := ScriptIntFromSlice(main.Pop()).Int64()
		min := ScriptIntFromSlice(main.Pop()).Int64()
		a := ScriptIntFromSlice(main.Pop()).Int64()
		if min <= a && a < max {
			main.PushByte(1)
		} else {
			main.PushByte(0)
		}
	default:
		return fmt.Errorf("0x%02x not a Script Arithmetic op", op)
	}
	return nil
}

const (
	OP_0         = 0x00
	OP_FALSE     = OP_0
	OP_PUSHDATA1 = 0x4c
	OP_PUSHDATA2 = 0x4d
	OP_PUSHDATA4 = 0x4e
	OP_1NEGATE   = 0x4f
	OP_RESERVED  = 0x50
	OP_1         = 0x51
	OP_TRUE      = OP_1
	OP_2         = 0x52
	OP_3         = 0x53
	OP_4         = 0x54
	OP_5         = 0x55
	OP_6         = 0x56
	OP_7         = 0x57
	OP_8         = 0x58
	OP_9         = 0x59
	OP_10        = 0x5a
	OP_11        = 0x5b
	OP_12        = 0x5c
	OP_13        = 0x5d
	OP_14        = 0x5e
	OP_15        = 0x5f
	OP_16        = 0x60
	// control
	OP_NOP      = 0x61
	OP_VER      = 0x62
	OP_IF       = 0x63
	OP_NOTIF    = 0x64
	OP_VERIF    = 0x65
	OP_VERNOTIF = 0x66
	OP_ELSE     = 0x67
	OP_ENDIF    = 0x68
	OP_VERIFY   = 0x69
	OP_RETURN   = 0x6a
	// stack ops
	OP_TOALTSTACK   = 0x6b
	OP_FROMALTSTACK = 0x6c
	OP_2DROP        = 0x6d
	OP_2DUP         = 0x6e
	OP_3DUP         = 0x6f
	OP_2OVER        = 0x70
	OP_2ROT         = 0x71
	OP_2SWAP        = 0x72
	OP_IFDUP        = 0x73
	OP_DEPTH        = 0x74
	OP_DROP         = 0x75
	OP_DUP          = 0x76
	OP_NIP          = 0x77
	OP_OVER         = 0x78
	OP_PICK         = 0x79
	OP_ROLL         = 0x7a
	OP_ROT          = 0x7b
	OP_SWAP         = 0x7c
	OP_TUCK         = 0x7d
	// splice ops
	OP_CAT    = 0x7e
	OP_SUBSTR = 0x7f
	OP_LEFT   = 0x80
	OP_RIGHT  = 0x81
	OP_SIZE   = 0x82
	// bit logic
	OP_INVERT      = 0x83
	OP_AND         = 0x84
	OP_OR          = 0x85
	OP_XOR         = 0x86
	OP_EQUAL       = 0x87
	OP_EQUALVERIFY = 0x88
	OP_RESERVED1   = 0x89
	OP_RESERVED2   = 0x8a
	// numeric
	OP_1ADD               = 0x8b
	OP_1SUB               = 0x8c
	OP_2MUL               = 0x8d
	OP_2DIV               = 0x8e
	OP_NEGATE             = 0x8f
	OP_ABS                = 0x90
	OP_NOT                = 0x91
	OP_0NOTEQUAL          = 0x92
	OP_ADD                = 0x93
	OP_SUB                = 0x94
	OP_MUL                = 0x95
	OP_DIV                = 0x96
	OP_MOD                = 0x97
	OP_LSHIFT             = 0x98
	OP_RSHIFT             = 0x99
	OP_BOOLAND            = 0x9a
	OP_BOOLOR             = 0x9b
	OP_NUMEQUAL           = 0x9c
	OP_NUMEQUALVERIFY     = 0x9d
	OP_NUMNOTEQUAL        = 0x9e
	OP_LESSTHAN           = 0x9f
	OP_GREATERTHAN        = 0xa0
	OP_LESSTHANOREQUAL    = 0xa1
	OP_GREATERTHANOREQUAL = 0xa2
	OP_MIN                = 0xa3
	OP_MAX                = 0xa4
	OP_WITHIN             = 0xa5
	// crypto
	OP_RIPEMD160           = 0xa6
	OP_SHA1                = 0xa7
	OP_SHA256              = 0xa8
	OP_HASH160             = 0xa9
	OP_HASH256             = 0xaa
	OP_CODESEPARATOR       = 0xab
	OP_CHECKSIG            = 0xac
	OP_CHECKSIGVERIFY      = 0xad
	OP_CHECKMULTISIG       = 0xae
	OP_CHECKMULTISIGVERIFY = 0xaf
	// expansion
	OP_NOP1                = 0xb0
	OP_CHECKLOCKTIMEVERIFY = 0xb1
	OP_NOP2                = OP_CHECKLOCKTIMEVERIFY
	OP_CHECKSEQUENCEVERIFY = 0xb2
	OP_NOP3                = OP_CHECKSEQUENCEVERIFY
	OP_NOP4                = 0xb3
	OP_NOP5                = 0xb4
	OP_NOP6                = 0xb5
	OP_NOP7                = 0xb6
	OP_NOP8                = 0xb7
	OP_NOP9                = 0xb8
	OP_NOP10               = 0xb9
	// template matching params
	OP_SMALLINTEGER = 0xfa
	OP_PUBKEYS      = 0xfb
	OP_PUBKEYHASH   = 0xfd
	OP_PUBKEY       = 0xfe
	// invalid code
	OP_INVALIDOPCODE = 0xff
)

func OpcodeName(opcode byte) string {
	switch opcode {
	case OP_0:
		return "OP_0"
	case OP_PUSHDATA1:
		return "OP_PUSHDATA1"
	case OP_PUSHDATA2:
		return "OP_PUSHDATA2"
	case OP_PUSHDATA4:
		return "OP_PUSHDATA4"
	case OP_1NEGATE:
		return "OP_1NEGATE"
	case OP_RESERVED:
		return "OP_RESERVED"
	case OP_1:
		return "OP_1"
	case OP_2:
		return "OP_2"
	case OP_3:
		return "OP_3"
	case OP_4:
		return "OP_4"
	case OP_5:
		return "OP_5"
	case OP_6:
		return "OP_6"
	case OP_7:
		return "OP_7"
	case OP_8:
		return "OP_8"
	case OP_9:
		return "OP_9"
	case OP_10:
		return "OP_10"
	case OP_11:
		return "OP_11"
	case OP_12:
		return "OP_12"
	case OP_13:
		return "OP_13"
	case OP_14:
		return "OP_14"
	case OP_15:
		return "OP_15"
	case OP_16:
		return "OP_16"
	case OP_NOP:
		return "OP_NOP"
	case OP_VER:
		return "OP_VER"
	case OP_IF:
		return "OP_IF"
	case OP_NOTIF:
		return "OP_NOTIF"
	case OP_VERIF:
		return "OP_VERIF"
	case OP_VERNOTIF:
		return "OP_VERNOTIF"
	case OP_ELSE:
		return "OP_ELSE"
	case OP_ENDIF:
		return "OP_ENDIF"
	case OP_VERIFY:
		return "OP_VERIFY"
	case OP_RETURN:
		return "OP_RETURN"
	case OP_TOALTSTACK:
		return "OP_TOALTSTACK"
	case OP_FROMALTSTACK:
		return "OP_FROMALTSTACK"
	case OP_2DROP:
		return "OP_2DROP"
	case OP_2DUP:
		return "OP_2DUP"
	case OP_3DUP:
		return "OP_3DUP"
	case OP_2OVER:
		return "OP_2OVER"
	case OP_2ROT:
		return "OP_2ROT"
	case OP_2SWAP:
		return "OP_2SWAP"
	case OP_IFDUP:
		return "OP_IFDUP"
	case OP_DEPTH:
		return "OP_DEPTH"
	case OP_DROP:
		return "OP_DROP"
	case OP_DUP:
		return "OP_DUP"
	case OP_NIP:
		return "OP_NIP"
	case OP_OVER:
		return "OP_OVER"
	case OP_PICK:
		return "OP_PICK"
	case OP_ROLL:
		return "OP_ROLL"
	case OP_ROT:
		return "OP_ROT"
	case OP_SWAP:
		return "OP_SWAP"
	case OP_TUCK:
		return "OP_TUCK"
	case OP_CAT:
		return "OP_CAT"
	case OP_SUBSTR:
		return "OP_SUBSTR"
	case OP_LEFT:
		return "OP_LEFT"
	case OP_RIGHT:
		return "OP_RIGHT"
	case OP_SIZE:
		return "OP_SIZE"
	case OP_INVERT:
		return "OP_INVERT"
	case OP_AND:
		return "OP_AND"
	case OP_OR:
		return "OP_OR"
	case OP_XOR:
		return "OP_XOR"
	case OP_EQUAL:
		return "OP_EQUAL"
	case OP_EQUALVERIFY:
		return "OP_EQUALVERIFY"
	case OP_RESERVED1:
		return "OP_RESERVED1"
	case OP_RESERVED2:
		return "OP_RESERVED2"
	case OP_1ADD:
		return "OP_1ADD"
	case OP_1SUB:
		return "OP_1SUB"
	case OP_2MUL:
		return "OP_2MUL"
	case OP_2DIV:
		return "OP_2DIV"
	case OP_NEGATE:
		return "OP_NEGATE"
	case OP_ABS:
		return "OP_ABS"
	case OP_NOT:
		return "OP_NOT"
	case OP_0NOTEQUAL:
		return "OP_0NOTEQUAL"
	case OP_ADD:
		return "OP_ADD"
	case OP_SUB:
		return "OP_SUB"
	case OP_MUL:
		return "OP_MUL"
	case OP_DIV:
		return "OP_DIV"
	case OP_MOD:
		return "OP_MOD"
	case OP_LSHIFT:
		return "OP_LSHIFT"
	case OP_RSHIFT:
		return "OP_RSHIFT"
	case OP_BOOLAND:
		return "OP_BOOLAND"
	case OP_BOOLOR:
		return "OP_BOOLOR"
	case OP_NUMEQUAL:
		return "OP_NUMEQUAL"
	case OP_NUMEQUALVERIFY:
		return "OP_NUMEQUALVERIFY"
	case OP_NUMNOTEQUAL:
		return "OP_NUMNOTEQUAL"
	case OP_LESSTHAN:
		return "OP_LESSTHAN"
	case OP_GREATERTHAN:
		return "OP_GREATERTHAN"
	case OP_LESSTHANOREQUAL:
		return "OP_LESSTHANOREQUAL"
	case OP_GREATERTHANOREQUAL:
		return "OP_GREATERTHANOREQUAL"
	case OP_MIN:
		return "OP_MIN"
	case OP_MAX:
		return "OP_MAX"
	case OP_WITHIN:
		return "OP_WITHIN"
	case OP_RIPEMD160:
		return "OP_RIPEMD160"
	case OP_SHA1:
		return "OP_SHA1"
	case OP_SHA256:
		return "OP_SHA256"
	case OP_HASH160:
		return "OP_HASH160"
	case OP_HASH256:
		return "OP_HASH256"
	case OP_CODESEPARATOR:
		return "OP_CODESEPARATOR"
	case OP_CHECKSIG:
		return "OP_CHECKSIG"
	case OP_CHECKSIGVERIFY:
		return "OP_CHECKSIGVERIFY"
	case OP_CHECKMULTISIG:
		return "OP_CHECKMULTISIG"
	case OP_CHECKMULTISIGVERIFY:
		return "OP_CHECKMULTISIGVERIFY"
	case OP_NOP1:
		return "OP_NOP1"
	case OP_CHECKLOCKTIMEVERIFY:
		return "OP_CHECKLOCKTIMEVERIFY"
	case OP_CHECKSEQUENCEVERIFY:
		return "OP_CHECKSEQUENCEVERIFY"
	case OP_NOP4:
		return "OP_NOP4"
	case OP_NOP5:
		return "OP_NOP5"
	case OP_NOP6:
		return "OP_NOP6"
	case OP_NOP7:
		return "OP_NOP7"
	case OP_NOP8:
		return "OP_NOP8"
	case OP_NOP9:
		return "OP_NOP9"
	case OP_NOP10:
		return "OP_NOP10"
	case OP_SMALLINTEGER:
		return "OP_SMALLINTEGER"
	case OP_PUBKEYS:
		return "OP_PUBKEYS"
	case OP_PUBKEYHASH:
		return "OP_PUBKEYHASH"
	case OP_PUBKEY:
		return "OP_PUBKEY"
	case OP_INVALIDOPCODE:
		return "OP_INVALIDOPCODE"
	}
	return fmt.Sprintf("0x%02x", opcode)
}
