package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"log"
)

type P struct {
	X, Y, Z int
	Name    string
}

type Q struct {
	X, Y *int32
	Name string
}

func main() {
	data, err := json.Marshal(&P{
		X:    1,
		Y:    2,
		Z:    3,
		Name: "pp",
	})
	if err != nil {
		log.Panicln(err)
	}
	dec := gob.NewDecoder(bytes.NewReader(data))
	m := P{}
	err = dec.Decode(&m)
	if err != nil {
		log.Println(err)
	}
	log.Println(m)
}
