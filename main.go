package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

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

func (r *OthelloRecode) InitBoard() {
	// r.Board = matrix.New(8, 8, )
	board := make([]float64, 64)
	count := 0
	for _, c := range r.BO {
		if c == '-' {
			board[count] = BLANK
			count++
		} else if c == 'O' {
			board[count] = WHITE
			count++
		} else if c == '*' {
			board[count] = BLACK
			count++
		}
		if count >= 64 {
			break
		}
	}
	r.Board = matrix.New(8, 8, board)
}

func GetHand(place string) Hand {
	hand := new(Hand)
	if place == "pass" {
		hand.Pass = true
		return *hand
	}
	if len(place) != 2 {
		fmt.Println(place)
		panic("SOMTEHING ERROR HAPPEN AT GETPOINT")
	}
	hand.Y = int(place[0] - 64)
	hand.X = int(place[1] - 48)
	return *hand
}

func (r *OthelloRecode) BlackRecode() []Hand {
	r.BHand = make([]Hand, 0)
	for _, rec := range r.B {
		c := strings.Split(rec, "/")
		if c[0] == "" {
			continue
		}
		r.BHand = append(r.BHand, GetHand(c[0]))
	}
	return r.BHand
}

func (r *OthelloRecode) WhiteRecode() []Hand {
	r.WHand = make([]Hand, 0)
	for _, rec := range r.W {
		c := strings.Split(rec, "/")
		if c[0] == "" {
			continue
		}
		r.WHand = append(r.WHand, GetHand(c[0]))
	}
	return r.WHand
}

func (r *OthelloRecode) ShowBoard() {
	for i, b := range r.Board.RawMatrix() {
		if b == WHITE {
			fmt.Printf("●")
		} else if b == BLACK {
			fmt.Printf("○")
		} else if b == BLANK {
			fmt.Printf("-")
		}
		if (i+1)%8 == 0 {
			fmt.Println()
		}
	}
}

func (r *OthelloRecode) Opponent(turn int) int {
	if turn == WHITE {
		return BLACK
	}
	return WHITE
}

func (r *OthelloRecode) reverse(turn, y, x, addY, addX int) bool {
	val, err := r.Board.At(y+addY, x+addX)
	if err != nil {
		return false
	} else if int(val) == turn {
		return true
	} else if int(val) == BLANK {
		return false
	} else if int(val) == r.Opponent(turn) {
		reverse := r.reverse(turn, y+addY, x+addX, addY, addX)
		if reverse {
			r.Board.Set(y+addY, x+addX, float64(turn))
		}
		return reverse
	}
	return false
}

func (r *OthelloRecode) SetHand(hand Hand, turn int) {
	if hand.Pass {
		return
	}
	r.Board.Set(hand.Y, hand.X, float64(turn))
	r.reverse(turn, hand.Y, hand.X, 1, 0)
	r.reverse(turn, hand.Y, hand.X, -1, 0)
	r.reverse(turn, hand.Y, hand.X, 0, 1)
	r.reverse(turn, hand.Y, hand.X, 0, -1)
	r.reverse(turn, hand.Y, hand.X, 1, 1)
	r.reverse(turn, hand.Y, hand.X, -1, 1)
	r.reverse(turn, hand.Y, hand.X, 1, -1)
	r.reverse(turn, hand.Y, hand.X, -1, -1)
}

func (r *OthelloRecode) Winner() string {
	white := 0
	black := 0
	for _, val := range r.Board.RawMatrix() {
		if val == WHITE {
			white++
		} else if val == BLACK {
			black++
		}
	}
	if white > black {
		return "white"
	} else if white < black {
		return "black"
	}
	return "draw"
}

func (r *OthelloRecode) ShowGame() {
	fmt.Println("Start game")
	for i, val := range r.WHand {
		fmt.Println(r.BHand[i])
		if !r.BHand[i].Pass {
			fmt.Printf("BLACK: %v, %v\n", r.BHand[i].Y, r.BHand[i].X)
			r.SetHand(r.BHand[i], BLACK)
		}
		r.ShowBoard()
		time.Sleep(3 * time.Second)
		if !val.Pass {
			fmt.Printf("WHITE: %v, %v\n", val.Y, val.X)
			r.SetHand(val, WHITE)
		}
		r.ShowBoard()

		time.Sleep(3 * time.Second)
	}
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
