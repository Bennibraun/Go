package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	go_symbols := map[byte]string{
		0: ".",
		1: "○",
		2: "●",
	}

	var board = [9][9]byte{
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	// var board [9][9]byte

	coordinate := map[string]byte{
		"a": 0, "1": 0,
		"b": 1, "2": 1,
		"c": 2, "3": 2,
		"d": 3, "4": 3,
		"e": 4, "5": 4,
		"f": 5, "6": 5,
		"g": 6, "7": 6,
		"h": 7, "8": 7,
		"i": 8, "9": 8,
	}

	var turn byte
move:
	{
		fmt.Println("    simp GO (no captures)")
		fmt.Println("---------------------------")
		fmt.Println(" \"quit\" to quit anytime")
		fmt.Println("")
		fmt.Println("  1 2 3 4 5 6 7 8 9")
		var rowIndex = [9]string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
		for i := 0; i < len(board); i++ {
			fmt.Print(rowIndex[i] + " ")
			for j := 0; j < len(board[i]); j++ {
				fmt.Print(go_symbols[board[i][j]])
				fmt.Print(" ")
			}
			fmt.Println("")
		}

		reader := bufio.NewReader(os.Stdin)
		if turn%2 == 0 {
			fmt.Print("Black (○)")
		} else {
			fmt.Print("White (●)")
		}
		fmt.Println(" to move: ")
		fmt.Println("(ex: a1)")
		text, _ := reader.ReadString('\n')

		if strings.TrimRight(text, "\r\n") != "quit" {
			row, rok := coordinate[string(text[0])]
			if !rok {
				fmt.Printf("row %q is invalid\n", text[0])
			}
			col, cok := coordinate[string(text[1])]
			if !cok {
				fmt.Printf("col %q is invalid\n", text[1])
			}
			if rok && cok {
				if board[row][col] != 0 {
					fmt.Println("there's already a piece there you moron")
				} else {
					board[row][col] = turn%2 + 1
					turn++
				}
			}

			goto move
		}
	}
}
