//
// stack.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package bitcoin

import "fmt"

type stack struct {
	i int
	// size limit is from https://github.com/bitcoin/bitcoin/blob/master/src/script/interpreter.cpp:1031
	item [1000][]byte
}

func (s *stack) Reset() {
	s.i = -1
}

func (s *stack) PushSlice(b []byte) {
	s.i++
	s.item[s.i] = b
}

func (s *stack) PushByte(b byte) {
	s.PushSlice([]byte{b})
}

func (s *stack) Pop() []byte {
	b := s.item[s.i]
	s.i--
	return b
}

func (s *stack) Top() []byte {
	return s.item[s.i]
}

func (s *stack) Item(n int) []byte {
	return s.item[s.i-n]
}

func (s *stack) ItemsCount() ScriptInt {
	return ScriptInt{int64(s.i + 1)}
}

func (s *stack) String() string {
	str := ""
	for i := 0; i <= s.i; i++ {
		str += fmt.Sprintf("%3d) %s\n", i, byte2hex(s.item[i]))
	}
	return str
}

func byte2hex(b []byte) string {
	s := ""
	for i := range b {
		s += fmt.Sprintf("%02x", b[i])
	}
	return s
}

func slice2bool(b []byte) bool {
	for i := range b {
		if b[i] != 0 {
			// Can be negative zero
			if i == len(b)-1 && b[i] == 0x80 {
				return false
			}
			return true
		}
	}
	return false
}
