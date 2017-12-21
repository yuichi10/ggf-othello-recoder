package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yuichi10/matrix"
)

const (
	BLACK = 1
	WHITE = 2
	BLANK = 3
)

type OthelloRecode struct {
	GM    string
	PC    string
	DT    string
	PB    string
	PW    string
	RB    string
	RW    string
	TI    string
	TY    string
	RE    string
	BO    string
	B     []string
	W     []string
	Board *matrix.Matrix
	BHand []Hand
	WHand []Hand
}

type Hand struct {
	Y    int
	X    int
	Pass bool
}

func ReadRecode(raw io.Reader) (recode *OthelloRecode) {
	reader := bufio.NewReader(raw)
	_, err := reader.ReadString(';')
	if err != nil {
		panic(err)
	}
	recode = new(OthelloRecode)
	recode.B = make([]string, 2)
	recode.W = make([]string, 2)
	for {
		key, err := reader.ReadString('[')
		if err != nil {
			fmt.Println(err)
			break
		}
		body, err := reader.ReadString(']')
		if err != nil {
			fmt.Println(err)
			break
		}
		switch key[:len(key)-1] {
		case "GM":
			recode.GM = body[:len(body)-1]
		case "PC":
			recode.PC = body[:len(body)-1]
		case "DT":
			recode.DT = body[:len(body)-1]
		case "PB":
			recode.PB = body[:len(body)-1]
		case "PW":
			recode.PW = body[:len(body)-1]
		case "RB":
			recode.PB = body[:len(body)-1]
		case "RW":
			recode.RW = body[:len(body)-1]
		case "TI":
			recode.TI = body[:len(body)-1]
		case "TY":
			recode.TY = body[:len(body)-1]
		case "RE":
			recode.RE = body[:len(body)-1]
		case "BO":
			recode.BO = body[:len(body)-1]
		case "B":
			recode.B = append(recode.B, body[:len(body)-1])
		case "W":
			recode.W = append(recode.W, body[:len(body)-1])
		}
	}
	return
}

func (r *OthelloRecode) RecodeGame(recoder *Recoder) {
	fmt.Println("Start Recode")
	for i, val := range r.WHand {
		recoder.Write("black", r.BHand[i].Y, r.BHand[i].X, r.BHand[i].Pass, r.Board)
		if !r.BHand[i].Pass {
			r.SetHand(r.BHand[i], BLACK)
		}
		recoder.Write("white", r.WHand[i].Y, r.WHand[i].X, r.WHand[i].Pass, r.Board)
		if !val.Pass {
			r.SetHand(val, WHITE)
		}
	}
	recoder.WriteToFile(r.Winner())
}

func main() {
	file, err := os.Open("testdata/Othello.01e4.ggf")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	recoder := NewRecoder("Othello.01e4.ggf")
	defer recoder.Close()
	for scanner.Scan() {
		recode := ReadRecode(strings.NewReader(scanner.Text()))
		if recode.BO != "8 -------- -------- -------- ---O*--- ---*O--- -------- -------- -------- *" {
			continue
		}
		recode.InitBoard()
		// recode.ShowBoard()
		recode.BlackRecode()
		recode.WhiteRecode()
		// recode.ShowGame()
		recode.RecodeGame(recoder)
	}
}
