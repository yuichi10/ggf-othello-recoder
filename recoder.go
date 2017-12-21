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
	bufX  bytes.Buffer
	bufY  bytes.Buffer
	xFile *os.File
	yFile *os.File
}

func NewRecoder(fileNameX, fileNameY string) *Recoder {
	if _, err := os.Stat("output"); err != nil {
		err := os.Mkdir("output", 077)
		if err != nil {
			panic(err)
		}
	}
	xFile, err := os.Create(fmt.Sprintf("output/%s.txt", fileNameX))
	if err != nil {
		panic(err)
	}
	yFile, err := os.Create(fmt.Sprintf("output/%s.txt", fileNameY))
	if err != nil {
		panic(err)
	}
	recode := new(Recoder)
	recode.xFile = xFile
	recode.yFile = yFile
	return recode
}

func (r *Recoder) Write(turn string, y, x int, pass bool, board *matrix.Matrix) {
	if turn != writeTurn {
		return
	}
	b := ""
	for _, val := range board.RawMatrix() {
		b = fmt.Sprintf("%v %v", b, val)
	}
	strings.Trim(b, " ")
	r.bufX.WriteString(fmt.Sprintf("%v\n", b))
	passNum := 0
	if pass {
		passNum = 1
	}
	r.bufY.WriteString(fmt.Sprintf("%v %v %v\n", y, x, passNum))
}

func (r *Recoder) WriteToFile(winner string) {
	if winner != writeWinner {
		return
	}
	r.xFile.Write(r.bufX.Bytes())
	r.yFile.Write(r.bufY.Bytes())
	r.bufX.Reset()
	r.bufY.Reset()
}

func (r *Recoder) Close() {
	r.xFile.Close()
	r.yFile.Close()
}
