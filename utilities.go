package main

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
		newRow := row + dir[0]
		newCol := col + dir[1]
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

