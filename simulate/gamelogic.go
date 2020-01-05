package simulate

import (
	"fmt"

	"github.com/chainreaction/game"
	"github.com/chainreaction/utils"
)

func updateBoard(board *[][]game.Pixel, x int, y int) {
	return
}

// ChainReaction is called after each move and spreads the orbs on the board
func ChainReaction(gameInstance *game.Instance, x int, y int) error {
	board := gameInstance.Board

	if x < 0 && y < 0 && x > 31 && y > 31 {
		return fmt.Errorf("Given positions x %v and y %v are out of range", x, y)
	}

	board[x][y].DotCount++

	cnt := board[x][y].DotCount

	switch cnt {
	case 2:
		if utils.IsCorner(len(board), len(board[0]), x, y) {
			updateBoard(&board, x, y)
		}
	case 3:
		if utils.IsOnEdge(len(board), len(board[0]), x, y) {
			updateBoard(&board, x, y)
		}
	case 4:
		updateBoard(&board, x, y)
	}

	return nil
}
