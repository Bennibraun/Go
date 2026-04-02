package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"hash/fnv"
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

type Move struct {
	Row, Col int
}

// The directions that groups can be connected in
// useful for counting liberties and checking rules
var directions = []Move{
    {Row: -1, Col: 0},
    {Row: 1, Col: 0},
    {Row: 0, Col: -1},
    {Row: 0, Col: 1},
}


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

func rowLabel(i int) string {
	// just calculates the row label if we have more than 26 rows
	// same logic as in parseMove
	var label string = ""
	for i >= 0 {
		label = string('a'+(i%26)) + label
		i = i/26 - 1
	}
	return label
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
	// also need to handle >26 rows
	rowIndex := make([]string, len(board))
	for i := 0; i < len(board); i++ {
		rowIndex[i] = rowLabel(i)
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
func parseMove(move string) (int, int, error) {

	// Get rid of whitespace and force to lowercase
	move = strings.TrimSpace(strings.ToLower(move))

	// Check for quit command ("quit", "q", "exit")
	quit_cmds := []string{"quit", "q", "exit"}
	for _, cmd := range quit_cmds {
		if move == cmd {
			return 0, 0, ErrQuit
		}
	}

	// // Moves have to be a letter and number
	// if len(move) != 2 {
	// 	return 0, 0, fmt.Errorf("invalid move format: %q", move)
	// }

	// // Grab the row and col
	// r := move[0]
	// c := move[1]

	// Because the board can be any size, we need to be able to handle >26 rows and >9 cols
	// for rows we could do a-z, then aa-zz, etc.
	// for cols we just stick to regular integers allowing for more digits
	// So we need to split the move into the letter part and the number part
	var rPart string
	var cPart string
	for _, char := range move {
		if char >= 'a' && char <= 'z' {
			rPart += string(char)
		} else if char >= '0' && char <= '9' {
			cPart += string(char)
		} else {
			return 0, 0, fmt.Errorf("invalid character in move: %q", char)
		}
	}

	// now we convert the rPart to a row index. I think just adding their values together would work
	var r int = 0
	for _, char := range rPart {
		r = r*26 + int(char-'a'+1)
	}
	r = r - 1 // offset by 1 for 0-indexing

	// Convert cPart to column index
	c, err := strconv.Atoi(cPart)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid column: %q", cPart)
	}
	c = c - 1 // offset by 1 for 0-indexing

	return r, c, nil
}

func hashMove(board [][]Stone, row int, col int, stone Stone) uint64 {
	// Place the stone temporarily to get a hash for the board state
	// This will be used to add the board state to hashtory and to check ko rule beforehand

	// start the hashing
	h := fnv.New64a()
	// go through the board to build up the hash
	for r, currentRow := range board {
		for c, cell := range currentRow {
			if r == row && c == col {
				// Write the current color for this cell, since it's the move being checked
				h.Write([]byte{byte(stone)})
			} else {
				// Write what's already there in the board
				h.Write([]byte{byte(cell)})
			}
		}
	}

	// compute the hash value
	return h.Sum64()
}

// Place a stone on the board, checking for validity
// Passing a pointer to the board with "*" so any changes we make to it affect the original board, not a copy
// Go has cool error handling. we set the function to return an "error" type, and nil means no error
func placeStone(board *[][]Stone, moveHashes *[]uint64, row int, col int, color Stone) (int, error) {
	// Check if the position exists
	if !isOnBoard(*board, row, col) {
		return 0, fmt.Errorf("position %s%d is out of bounds", rowLabel(row), col+1)
	}

	// Check if the position is already occupied
	if (*board)[row][col] != Empty {
		return 0, fmt.Errorf("position %s%d is already occupied", rowLabel(row), col+1)
	}

	// Check all the rules for this move. If any rule is violated, return the error message
	if err := checkAllRules(*board, *moveHashes, row, col, color); err != nil {
		return 0, err
	}

	// Place the stone
	(*board)[row][col] = color
	// Add the new board state to the move history
	*moveHashes = append(*moveHashes, hashMove(*board, row, col, color))

	// check for captures and remove any captured stones
	numRemoved := removeCapturedStones(*board, row, col, color)
	if numRemoved > 0 {
		fmt.Printf("Removed %d captured stones.\n", numRemoved)
	}


	return numRemoved, nil
}

func main() {

	// Keep track of turns. Even = black, odd = white
	var turn byte = 0

	// Keep track of moves. Made a struct with two ints, row & col
	moves := []Move{}
	// Also track board states by hashing them
	moveHashes := []uint64{}

	// track scores (number of captured stones) for each color
	var blackScore int = 0
	var whiteScore int = 0

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

	var board [][]Stone

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
		numRemoved, placeErr := placeStone(&board, &moveHashes, row, col, Stone((turn%2)+1))
		// Make sure it worked before incrementing the turn, otherwise if there was an error placing the stone we want the same player to try again
		if placeErr != nil {
			fmt.Println("Error placing stone:", placeErr)
		} else {
			// Save the move for Ko rule
			moves = append(moves, Move{Row: row, Col: col})
			// increment the score for the current color for any captured stones
			if numRemoved > 0 {
				if turn%2 == 0 {
					blackScore += numRemoved
				} else {
					whiteScore += numRemoved
				}
				fmt.Printf("Black score: %d, White score: %d\n", blackScore, whiteScore)
			}

			turn++
		}

	}
}

