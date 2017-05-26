//
// tx.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package bitcoin

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
)

type TxIn struct {
	PrevTx        DoubleHash
	PrevTxOutIndx uint32
	Script        []byte
	SequenceNum   uint32
}

type TxOut struct {
	Value  uint64
	Script []byte
}

type Tx struct {
	Verssion uint32
	In       []TxIn
	Out      []TxOut
	LockTime uint32
	Hash     DoubleHash
}

// IsCoinbase return true if this transaction is a generation transaction
// i.e. input for this transaction is a new generated block
func (t Tx) IsCoinbase() bool {
	null := DoubleHash{}
	return len(t.In) == 1 && bytes.Compare(null[:], t.In[0].PrevTx[:]) == 0
}

func ReadTxIn(r io.Reader) (t *TxIn, err error) {
	t = new(TxIn)
	_, err = r.Read(t.PrevTx[:])
	if err != nil {
		return
	}
	err = binary.Read(r, binary.LittleEndian, &t.PrevTxOutIndx)
	if err != nil {
		return
	}
	var n Varint
	err = ReadVarint(r, &n)
	if err != nil {
		return
	}
	t.Script = make([]byte, int(n))
	_, err = r.Read(t.Script)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.LittleEndian, &t.SequenceNum)
	if err != nil {
		return
	}
	return t, nil
}

func ReadTxOut(r io.Reader) (t *TxOut, err error) {
	t = new(TxOut)
	err = binary.Read(r, binary.LittleEndian, &t.Value)
	if err != nil {
		return
	}
	var n Varint
	err = ReadVarint(r, &n)
	if err != nil {
		return
	}
	t.Script = make([]byte, int(n))
	_, err = r.Read(t.Script)
	if err != nil {
		return
	}
	return t, nil
}

func ReadTx(r io.Reader) (t *Tx, err error) {
	t = new(Tx)
	err = binary.Read(r, binary.LittleEndian, &t.Verssion)
	if err != nil {
		return
	}
	var count Varint
	err = ReadVarint(r, &count)
	if err != nil {
		return
	}
	for i := Varint(0); i < count; i++ {
		var in *TxIn
		in, err = ReadTxIn(r)
		if err != nil {
			return
		}
		t.In = append(t.In, *in)
	}
	err = ReadVarint(r, &count)
	if err != nil {
		return
	}
	for i := Varint(0); i < count; i++ {
		var out *TxOut
		out, err = ReadTxOut(r)
		if err != nil {
			return
		}
		t.Out = append(t.Out, *out)
	}
	err = binary.Read(r, binary.LittleEndian, &t.LockTime)
	if err != nil {
		return
	}
	t.Hash.Update(t.Raw())
	return t, nil
}

func (tx *Tx) Raw() []byte {
	w := new(bytes.Buffer)
	err := binary.Write(w, binary.LittleEndian, tx.Verssion)
	if err != nil {
		panic(err)
	}
	inCount := Varint(len(tx.In))
	err = WriteVarint(w, inCount)
	if err != nil {
		panic(err)
	}
	for i := range tx.In {
		_, err = w.Write(tx.In[i].Raw())
		if err != nil {
			panic(err)
		}
	}
	outCount := Varint(len(tx.Out))
	err = WriteVarint(w, outCount)
	if err != nil {
		panic(err)
	}
	for i := range tx.Out {
		_, err = w.Write(tx.Out[i].Raw())
		if err != nil {
			panic(err)
		}
	}
	err = binary.Write(w, binary.LittleEndian, tx.LockTime)
	if err != nil {
		panic(err)
	}
	return w.Bytes()
}

func (tx *TxIn) Raw() []byte {
	w := new(bytes.Buffer)
	_, err := w.Write(tx.PrevTx[:])
	if err != nil {
		panic(err)
	}
	err = binary.Write(w, binary.LittleEndian, tx.PrevTxOutIndx)
	if err != nil {
		panic(err)
	}
	scriptLen := Varint(len(tx.Script))
	err = WriteVarint(w, scriptLen)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(tx.Script)
	if err != nil {
		panic(err)
	}
	err = binary.Write(w, binary.LittleEndian, tx.PrevTxOutIndx)
	if err != nil {
		panic(err)
	}
	return w.Bytes()
}

func (tx *TxOut) Raw() []byte {
	w := new(bytes.Buffer)
	err := binary.Write(w, binary.LittleEndian, tx.Value)
	if err != nil {
		panic(err)
	}
	scriptLen := Varint(len(tx.Script))
	err = WriteVarint(w, scriptLen)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(tx.Script)
	if err != nil {
		panic(err)
	}
	return w.Bytes()
}

func (tx TxIn) String() string {
	ob, err := json.MarshalIndent(&tx, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(ob)
}

func (tx TxOut) String() string {
	ob, err := json.MarshalIndent(&tx, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(ob)
}

func (tx Tx) String() string {
	ob, err := json.MarshalIndent(&tx, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(ob)
}
