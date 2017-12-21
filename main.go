package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuichi10/matrix"
)

var (
	writeFile   string
	writeWinner string
	writeTurn   string
)

func init() {
	flag.StringVar(&writeFile, "f", "", "set file path of Othello ggf")
	flag.StringVar(&writeWinner, "w", "black", "which winners recode will you set. if you do not set, it will be same to turn")
	flag.StringVar(&writeTurn, "t", "", "which turn you will recode. if you do not set it will be same to winner")
}

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

func checkArgumentStatus() error {
	if writeWinner != "black" && writeWinner != "white" {
		return fmt.Errorf("the winner just allow to use %q or %q", "black", "white")
	}
	if writeTurn != "black" && writeTurn != "white" {
		return fmt.Errorf("the winner just allow to use %q or %q", "black", "white")
	}
	if _, err := os.Stat(writeFile); err != nil {
		return fmt.Errorf("There is no such file: %v", err)
	}
	return nil
}

func initArgument() error {
	if writeFile == "" {
		return fmt.Errorf("you must set file")
	}
	if writeWinner == "" && writeTurn == "" {
		return fmt.Errorf("you must set winner or turn")
	}
	if writeTurn == "" {
		writeTurn = writeWinner
	} else if writeWinner == "" {
		writeWinner = writeTurn
	}
	if writeFile[0] != '/' {
		var err error
		writeFile, err = filepath.Abs(writeFile)
		if err != nil {
			return fmt.Errorf("failed to abs file")
		}
	}
	return checkArgumentStatus()
}

func recodeFileName() string {
	paths := strings.Split(writeFile, "/")
	return paths[len(paths)-1]
}

func dataFileName(base string) (string, string) {
	return fmt.Sprintf("%s_%s_%s_X", base, writeWinner, writeTurn), fmt.Sprintf("%s_%s_%s_Y", base, writeWinner, writeTurn)
}

func main() {
	flag.Parse()
	err := initArgument()
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Open(writeFile)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	recoder := NewRecoder(dataFileName(recodeFileName()))
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
