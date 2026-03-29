
/**
* @class board
* @method move
* @method apply private
* @method applyPiece private
* @method applyCapture private
* @method validate private
* @method validateKo private
* @method validateSuicide private
* @method undo
* @method render
* @method show private
* 
* @method getHashtory private
* @method addHashtory private
*/

import Move

package board

type board interface {
	move(Move) board
	apply(Move) board
	applyPiece(Move) board
	applyCapture(Move) board
	validate(Move) int
	validateKo(Move) int
	validateSuicide(Move) int
	undo() board
	render() board
	show() board
	getHashtory() board
	addHashtory() board
}

func (b board) move(m Move) board {
	if apply(m).validate(m) == 0 {
		b.undo()
	}
	return b
}

func (b board) apply(m) board {
	return applyPiece(m).applyCapture(m)
}

func (b board) validate(m) int {

}

