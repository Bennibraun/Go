package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Making a custom type for stones so we don't have to remember 0,1,2 for colors
// iota is a shortcut that just gives incrementing numbers starting from 0
// so Empty = 0, Black = 1, White = 2
type Stone byte

const (
	Empty Stone = iota
	Black
	White
)

// Define error type for quitting (called "sentinel error" in Go)
var ErrQuit = fmt.Errorf("quit")

// Make the game board at startup
func newBoard(x int, y int) [][]Stone {
	// Arrays have to have a fixed hard-coded size, but slices can be dynamic
	// so we use a 2D slice so the baord can be any size
	board := make([][]Stone, x)
	for i := range board {
		board[i] = make([]Stone, y)
	}
	return board
}

// just print the board
func printBoard(board [][]Stone) {

	go_symbols := map[Stone]string{
		Empty: ".",
		Black: "○",
		White: "●",
	}

	// Needs to be dynamic based on board size
	// kind of tricky but basically just generating the numbers from 0 to board witdh and joins with spaces
	s := make([]string, len(board[0]))
	for i := 0; i < len(board[0]); i++ {
		s[i] = strconv.Itoa(i + 1) // convert int to string and offset by 1 since no 0 column
	}
	fmt.Println("  ", strings.Join(s, " "))

	// Same for the row labels but letters instead of numbers
	rowIndex := make([]string, len(board))
	for i := 0; i < len(board); i++ {
		rowIndex[i] = string('a' + i) // can do arithmetic on ascii
	}

	for i := 0; i < len(board); i++ {
		fmt.Print(rowIndex[i] + " ")
		for j := 0; j < len(board[i]); j++ {
			fmt.Print(go_symbols[board[i][j]])
			fmt.Print(" ")
		}
		fmt.Println("")
	}
}

// Read the move from cli input and convert to row/col indices
func parseMove(move string) (byte, byte, error) {

	// Get rid of whitespace and force to lowercase
	move = strings.TrimSpace(strings.ToLower(move))

	// Check for quit command ("quit", "q", "exit")
	quit_cmds := []string{"quit", "q", "exit"}
	for _, cmd := range quit_cmds {
		if move == cmd {
			return 0, 0, ErrQuit
		}
	}

	// Moves have to be a letter and number
	if len(move) != 2 {
		return 0, 0, fmt.Errorf("invalid move format: %q", move)
	}

	// Grab the row and col
	r := move[0]
	c := move[1]

	// I want to do better than this:
	// coordinate := map[string]byte{
	// 	"a": 0, "1": 0,
	// 	"b": 1, "2": 1,
	// 	"c": 2, "3": 2,
	// 	"d": 3, "4": 3,
	// 	"e": 4, "5": 4,
	// 	"f": 5, "6": 5,
	// 	"g": 6, "7": 6,
	// 	"h": 7, "8": 7,
	// 	"i": 8, "9": 8,
	// }

	// go apparently lets you do arithmetic on letters since they are just bytes, like 'a' + 2 == 'c'
	// So we really don't even need to make a coordinate table

	// Check if the row is in acceptable range
	// We just need to offset the letter by the value of 'a' to convert to 0-8 index
	// TODO: fix hardcoding for 9x9
	if r < 'a' || r > 'i' {
		return 0, 0, fmt.Errorf("invalid row: %q (expected a-i)", r)
	}

	// Check the column
	if c < '1' || c > '9' {
		return 0, 0, fmt.Errorf("invalid column: %q (expected 1-9)", c)
	}

	// Looks good, return indices
	return r - 'a', c - '1', nil
}

// Place a stone on the board, checking for validity
// Passing a pointer to the board with "*" so any changes we make to it affect the original board, not a copy
// Go has cool error handling. we set the function to return an "error" type, and nil means no error
func placeStone(board *[][]Stone, row byte, col byte, color Stone) error {
	// Check if the position is already occupied
	if (*board)[row][col] != Empty {
		return fmt.Errorf("position %c%d is already occupied", 'a'+row, col+1)
	}

	// TODO: add more rule checks if needed

	// Place the stone
	(*board)[row][col] = color
	return nil
}

func main() {

	// var board = [9][9]byte{
	// 	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	// 	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	// 	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	// 	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	// 	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	// 	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	// 	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	// 	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	// 	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	// }

	board := newBoard(9, 9)

	// Keep track of turns. Even = black, odd = white
	var turn byte = 0

	// Initialize reader to get cli input
	reader := bufio.NewReader(os.Stdin)

	// Print welcome message and instructions
	fmt.Println("    simp GO (no captures)")
	fmt.Println("---------------------------")
	fmt.Println("Enter moves in the format: a1, b3, etc.")
	fmt.Println(" \"quit\" to quit anytime")
	fmt.Println("")

	// Setup
	fmt.Println("What size board would you like?")
	fmt.Println("a. 9x9")
	fmt.Println("b. 13x13")
	fmt.Println("c. 19x19")
	fmt.Println("d. Custom")

	board_size_input, _ := reader.ReadString('\n')
	board_size_input = strings.TrimSpace(strings.ToLower(board_size_input))
	switch board_size_input {
	case "a":
		board = newBoard(9, 9)
	case "b":
		board = newBoard(13, 13)
	case "c":
		board = newBoard(19, 19)
	case "d":
		fmt.Println("Enter custom board size in this format: 17x9")
		custom_size_input, _ := reader.ReadString('\n')
		custom_size_input = strings.TrimSpace(strings.ToLower(custom_size_input))
		// Parse the custom size input (kind of error-prone but whatever)
		// split the input on "x" and convert to integers
		size_parts := strings.Split(custom_size_input, "x")
		if len(size_parts) != 2 {
			fmt.Println("Invalid custom size format. Defaulting to 9x9.")
			board = newBoard(9, 9)
		} else {
			rows, err1 := strconv.Atoi(size_parts[0])
			cols, err2 := strconv.Atoi(size_parts[1])
			if err1 != nil || err2 != nil {
				fmt.Println("Invalid custom size format. Defaulting to 9x9.")
				board = newBoard(9, 9)
			} else {
				board = newBoard(rows, cols)
			}
		}
	}

	// Main game loop
	for {

		printBoard(board)

		if turn%2 == 0 {
			fmt.Print("Black (○)")
		} else {
			fmt.Print("White (●)")
		}
		fmt.Println(" to move: ")
		fmt.Println("(ex: a1)")

		// Get move input
		text, _ := reader.ReadString('\n')
		row, col, err := parseMove(text)

		// Check if input was a quit command
		if errors.Is(err, ErrQuit) {
			fmt.Println("Thanks for playing!")
			break
		}

		// Check for other parsing errors
		if err != nil {
			fmt.Println("Error parsing move:", err)
			continue
		}

		// Place the stone on the board
		placeErr := placeStone(&board, row, col, Stone((turn%2)+1))
		// Make sure it worked before incrementing the turn, otherwise if there was an error placing the stone we want the same player to try again
		if placeErr != nil {
			fmt.Println("Error placing stone:", placeErr)
		} else {
			turn++
		}

	}
}
