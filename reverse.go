package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/yuichi10/matrix"
)

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
