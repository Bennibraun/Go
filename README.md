# Contents

[**1. CLI**](#cli)
	[1.1. Start](#start)
	[1.2. Name](#name)
	[1.3. Rulebook](#rulebook)
	[1.4. Move](#move)
	[1.5. Win](#win)
[**2. Types**](#types)
	[**2.1. Piece**](#piece)
	[**2.2. Move**](#move)
	[**2.3. Size**](#size)
	[**2.4. Side**](#side)
[**3. Classes**](#classes)
	[3.1. Board](#board)
	[3.2. Player](#player)
	[3.3. Game](#game)

---
# CLI

## Start

## Name

## Rulebook
## Move
## Win

---

## Types

### Piece
`0 | 1 | 2`
*(empty, black, white)*

### Move
`nil | [byte, byte]`
*(pass, a1, a2, ..., b1, ...)*

### Size
`0 | 1 | 2`
*(9x9, 13x13, 19x19)*

### Side
`1 | 2`
*(black, white)*

---
# Classes
## Board

```go
interface board {
	Size():    size
	Show():    piece[][]
	History(): move[]
}
```

## Player

```go
interface player {
  Name():     string
  Move(move): player, exception
}
```

## Game

```go
interface game {
	Turn():     side
	Board():    board
	Black():    player
	White():    player
	Rulebook(): rulebook
	Move(move): game, exception
}
```

---
