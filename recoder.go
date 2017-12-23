package main

import (
	"bytes"
	"fmt"
	"os"

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
	b := board.TextAsOneLine(" ")
	r.bufX.WriteString(fmt.Sprintf("%v\n", b))

	yBoard := matrix.New(1, 65, nil)
	if pass {
		yBoard.Set(1, 65, 1.0)
		fmt.Println("PASS")
	} else {
		yBoard.Set(1, (y-1)*8+x, 1.0)
	}

	r.bufY.WriteString(fmt.Sprintf("%s\n", yBoard.TextAsOneLine(" ")))
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
