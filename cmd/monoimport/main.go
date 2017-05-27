//
// main.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/weirdgiraffe/bitcoin"

	mgo "gopkg.in/mgo.v2"
)

func main() {
	blockDir := os.Args[1]
	dbHost := "localhost"

	if v, ok := os.LookupEnv("MONGO_HOST"); ok {
		dbHost = v
	}
	session, err := mgo.Dial(dbHost)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	tc := session.DB("btc").C("tx")
	bc := session.DB("btc").C("block")

	files, err := ioutil.ReadDir(blockDir)
	if err != nil {
		log.Fatal(err)
	}

	count := 0
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		bf, err := bitcoin.OpenBlockFile(blockDir + "/" + f.Name())
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < bf.BlockCount(); i++ {
			b, err := bf.Block(i)
			if err != nil {
				bf.Close()
				log.Fatal(err)
			}
			err = bc.Insert(b)
			if err != nil {
				bf.Close()
				log.Fatal(err)
			}
			fmt.Printf("%9d BLOCK: %s\r", count, b.Hash)
			os.Stdout.Sync()
			docs := make([]interface{}, b.TxCount())
			for i := 0; i < b.TxCount(); i++ {
				docs[i] = b.Tx(i)
			}
			bulk := tc.Bulk()
			bulk.Insert(docs...)
			_, err = bulk.Run()
			if err != nil {
				bf.Close()
				log.Fatal(err)
			}
			count++
		}
		bf.Close()
	}
	fmt.Println()
}
