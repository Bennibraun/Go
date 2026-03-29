package main

func checkAllRules(board [][]Stone, moveHashes []uint64, row int, col int, stone Stone) error {
	if !checkKoRule(board, moveHashes, row, col, stone) {
		// return error: Ko violation
	}
	// return no error
}

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
			return false
		}
	}

	return true
}