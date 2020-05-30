package simulate

import (
	"fmt"

	"github.com/chainreaction/helpers"
	"github.com/chainreaction/models"
	"github.com/chainreaction/utils"
)

func updateBoard(gameInstance *models.Instance, x int, y int, color string) []int {
	var states []int
	dim := gameInstance.Dimension
	board := &gameInstance.Board
	q := utils.NewQueue()
	q.Enqueue(utils.Pair{x, y})

	for !q.IsEmpty() {
		x, y = q.Dequeue()

		if c := (*board)[dim*x+y].Color; c != "" && c != color {
			gameInstance.DecCellCountOfPlayer(c)
		}

		(*board)[dim*x+y].DotCount++
		(*board)[dim*x+y].Color = color

		gameInstance.IncCellCountOfPlayer(color)

		cnt := (*board)[dim*x+y].DotCount

		states = append(states, []int{x, y, cnt}...)

		if isChaining(cnt, dim, x, y) {
			states = append(states, []int{x, y, 0}...)
			chain(board, x, y, color, q, dim)
			states = append(states, -1)
		}

	}

	return states
}

func isChaining(cnt int, dim int, x int, y int) bool {
	if utils.IsCorner(dim, dim, x, y) && cnt == 2 {
		return true
	} else if utils.IsOnEdge(dim, dim, x, y) && cnt == 3 {
		return true
	} else if cnt == 4 {
		return true
	}
	return false
}

func chain(board *[]models.Pixel, x int, y int, color string, q *utils.Queue, dim int) {
	(*board)[dim*x+y].DotCount = 0
	(*board)[dim*x+y].Color = ""

	if x > 0 {
		q.Enqueue(utils.Pair{x - 1, y})
	}
	if y > 0 {
		q.Enqueue(utils.Pair{x, y - 1})
	}
	if x < dim-1 {
		q.Enqueue(utils.Pair{x + 1, y})
	}
	if y < dim-1 {
		q.Enqueue(utils.Pair{x, y + 1})
	}
}

func checkIfWon(gI *models.Instance, color string) bool {
	won := true
	if !helpers.CheckIfEveryonePlayed(gI) {
		return false
	}
	for i := 0; i < gI.Dimension; i++ {
		for j := 0; j < gI.Dimension; j++ {
			if gI.Board[gI.Dimension*i+j].DotCount != 0 {
				if gI.Board[gI.Dimension*i+j].Color != color {
					won = false
				}
			}
		}
	}
	return won
}

// ChainReaction is called after each move and spreads the orbs on the board
func ChainReaction(gameInstance *models.Instance, move models.MoveMsg) error {
	board := gameInstance.Board
	player, exists := gameInstance.GetPlayerByUsername(move.PlayerUserName)

	if !exists {
		return fmt.Errorf("Username doesn't exists for this game %v", move.PlayerUserName)
	}

	x, y := move.XPos, move.YPos

	if x < 0 && y < 0 && x > gameInstance.Dimension && y > gameInstance.Dimension {
		return fmt.Errorf("Given positions x %v and y %v are out of range", x, y)
	}

	if board[gameInstance.Dimension*x+y].DotCount != 0 &&
		board[gameInstance.Dimension*x+y].Color != player.Color {
		return fmt.Errorf("Invalid move. board[%v][%v] already contains color: %v", x, y, board[gameInstance.Dimension*x+y].Color)
	}

	states := updateBoard(gameInstance, x, y, player.Color)

	won := checkIfWon(gameInstance, player.Color)

	helpers.SetIfAllPlayedOnce(gameInstance, player.UserName)

	gameInstance.UpdatedBoard <- states

	if won {
		gameInstance.SetWinner(player)
	}

	return nil
}
