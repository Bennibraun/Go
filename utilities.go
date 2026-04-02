package main

import "fmt"


// Just returns the opposite stone color
func oppositeColor(stone Stone) Stone {
	switch stone {
	case Black:
		return White
	case White:
		return Black
	default:
		return Empty
	}
}

// Handy check for if a position exists (to handle edges etc.)
func isOnBoard(board [][]Stone, row int, col int) bool {
	return row >= 0 && row < len(board) && col >= 0 && col < len(board[0])
}

// count the liberties of a group starting at a partiular stone
func countLiberties(board [][]Stone, row int, col int) int {
	stoneColor := board[row][col]
	visited := make(map[string]bool)
	return countLibertiesHelper(board, row, col, stoneColor, visited)
}

// the recursive function that handles liberty counts
func countLibertiesHelper(board [][]Stone, row int, col int, stoneColor Stone, visited map[string]bool) int {
	// If we've already visited this position, return 0 to avoid infinite recursion
	posKey := fmt.Sprintf("%d,%d", row, col)
	if visited[posKey] {
		return 0
	}
	visited[posKey] = true

	// Check the four adjacent positions
	liberties := 0
	for _, dir := range directions {
		newRow := row + dir.Row
		newCol := col + dir.Col
		if isOnBoard(board, newRow, newCol) {
			if board[newRow][newCol] == Empty {
				liberties++
			} else if board[newRow][newCol] == stoneColor {
				liberties += countLibertiesHelper(board, newRow, newCol, stoneColor, visited)
			}
		}
	}
	return liberties
}


// check for captured stones and remove them from the board
func removeCapturedStones(board [][]Stone, row int, col int, stone Stone) int {
	numRemoved := 0

	// Check the four adjacent positions for opponent's stones
	for _, dir := range directions {
		adjRow := row + dir.Row
		adjCol := col + dir.Col
		if isOnBoard(board, adjRow, adjCol) && board[adjRow][adjCol] == oppositeColor(stone) {
			if countLiberties(board, adjRow, adjCol) == 0 {
				numRemoved += removeGroup(board, adjRow, adjCol)
			}
		}
	}

	return numRemoved

}

// remove a group of stones starting at a particular position and return the number of stones removed
// same logic as counting liberties (merge into one function later?)
func removeGroup(board [][]Stone, row int, col int) int {
	stoneColor := board[row][col]
	visited := make(map[string]bool)
	return removeGroupHelper(board, row, col, stoneColor, visited)
}

func removeGroupHelper(board [][]Stone, row int, col int, stoneColor Stone, visited map[string]bool) int {

	// If we've already visited this position, return 0 to avoid infinite recursion
	posKey := fmt.Sprintf("%d,%d", row, col)
	if visited[posKey] {
		return 0
	}

	visited[posKey] = true

	// Remove the stone
	board[row][col] = Empty
	numRemoved := 1

	// Check the four adjacent positions
	for _, dir := range directions {
		newRow := row + dir.Row
		newCol := col + dir.Col
		if isOnBoard(board, newRow, newCol) && board[newRow][newCol] == stoneColor {
			numRemoved += removeGroupHelper(board, newRow, newCol, stoneColor, visited)
		}
	}

	return numRemoved

}