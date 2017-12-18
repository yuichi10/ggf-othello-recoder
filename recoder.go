package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/yuichi10/matrix"
)

// Recoder recode othello recode
type Recoder struct {
	buf  bytes.Buffer
	file *os.File
}

func NewRecoder(fileName string) *Recoder {
	if _, err := os.Stat("output"); err != nil {
		err := os.Mkdir("output", 077)
		if err != nil {
			panic(err)
		}
	}
	file, err := os.Create(fmt.Sprintf("output/%s.txt", fileName))
	if err != nil {
		panic(err)
	}
	recode := new(Recoder)
	recode.file = file
	return recode
}

func (r *Recoder) Write(turn string, y, x int, pass bool, board *matrix.Matrix) {
	b := ""
	for _, val := range board.RawMatrix() {
		b = fmt.Sprintf("%v %v", b, val)
	}
	strings.Trim(b, " ")
	r.buf.WriteString(fmt.Sprintf("%v %v %v %v %v\n", turn, y, x, pass, b))
}

func (r *Recoder) WriteToFile(winner string) {
	r.file.Write([]byte(fmt.Sprintf("Winner %s\n", winner)))
	r.file.Write(r.buf.Bytes())
	r.buf.Reset()
}

func (r *Recoder) Close() {
	r.file.Close()
}
