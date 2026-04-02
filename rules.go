package main

import "fmt"

func checkAllRules(board [][]Stone, moveHashes []uint64, row int, col int, stone Stone) error {
	if checkKoRule(board, moveHashes, row, col, stone) {
		return fmt.Errorf("move at %s%d violates the Ko rule", rowLabel(row), col+1)
	}
	if checkSuicideRule(board, row, col, stone) {
		return fmt.Errorf("move at %s%d is suicidal", rowLabel(row), col+1)
	}

	return nil
}

// true if the move violates ko, false if not
func checkKoRule(board [][]Stone, moveHashes []uint64, row int, col int, stone Stone) bool {
	// Need to have some kind of move history to make sure board state isn't being repeated.
	// Shouldn't store all moves as 2D array, super inefficent
	// So need to come up with an algo that takes move history and current board state

	// Build up the unique board states on the fly and check for repeat?
	// Will try to use hashing

	currentHash := hashMove(board,row,col,stone)

	// Check if current hash is in the list of previous board states
	for _, h := range moveHashes {
		if h == currentHash {
			return true
		}
	}

	return false
}

// true if the move is suicidal, false if it's not
func checkSuicideRule(board [][]Stone, row int, col int, stone Stone) bool {

	// Temporarily place the stone on the board to check for liberties
	board[row][col] = stone

	// A cool Go thing is these anonymous "defer" functions
	// which will run after this function returns, regardless of how it returns
	defer func() {
		// Remove the stone after checking
		board[row][col] = Empty
	}()

	// first check  if the new stone has any liberties. If it does, then it's not suicidal.
	if countLiberties(board, row, col) > 0 {
		return false
	}

	// we need to count the liberties of all other-colored stones surrounding the new stone, and if any of them have 0 liberties after placing the new stone, then the move is not suicidal.
	for _, dir := range directions {
		adjRow := row + dir.Row
		adjCol := col + dir.Col
		if isOnBoard(board, adjRow, adjCol) && board[adjRow][adjCol] == oppositeColor(stone) {
			if countLiberties(board, adjRow, adjCol) == 0 {
				return false
			}
		}
	}

	return true

}