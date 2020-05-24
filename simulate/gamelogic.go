package simulate

import (
	"fmt"

	"github.com/chainreaction/helpers"
	"github.com/chainreaction/models"
	"github.com/chainreaction/utils"
)

func updateBoard(board *[]models.Pixel, x int, y int, color string, dim int) []int {
	var states []int
	q := utils.NewQueue()
	q.Enqueue(utils.Pair{x, y})

	(*board)[dim*x+y].DotCount++
	(*board)[dim*x+y].Color = color

	states = append(states, []int{x, y, (*board)[dim*x+y].DotCount}...)

	for !q.IsEmpty() {
		x, y = q.Dequeue()

		cnt := (*board)[dim*x+y].DotCount

		switch cnt {
		case 2:
			if utils.IsCorner(dim, dim, x, y) {
				updatePixelState(board, x, y, color, q, &states, dim)
			}
		case 3:
			if utils.IsOnEdge(dim, dim, x, y) {
				updatePixelState(board, x, y, color, q, &states, dim)
			}
		case 4:
			updatePixelState(board, x, y, color, q, &states, dim)
		}

	}

	return states
}

func updatePixelState(board *[]models.Pixel, x int, y int, color string, q *utils.Queue, t *[]int, dim int) {
	(*board)[dim*x+y].DotCount = 0
	(*board)[dim*x+y].Color = ""

	*t = append(*t, []int{x, y, 0}...)

	var p utils.Pair

	if x > 0 {
		(*board)[(dim*(x-1))+y].DotCount++
		(*board)[(dim*(x-1))+y].Color = color
		p = utils.Pair{x - 1, y}
		q.Enqueue(p)
		tmp := []int{x - 1, y, (*board)[(dim*(x-1))+y].DotCount}
		*t = append(*t, tmp...)
	}
	if y > 0 {
		(*board)[(dim*x)+(y-1)].DotCount++
		(*board)[(dim*x)+(y-1)].Color = color
		p = utils.Pair{x, y - 1}
		q.Enqueue(p)
		tmp := []int{x, y - 1, (*board)[(dim*x)+(y-1)].DotCount}
		*t = append(*t, tmp...)
	}
	if x < dim-1 {
		(*board)[(dim*(x+1))+y].DotCount++
		(*board)[(dim*(x+1))+y].Color = color
		p = utils.Pair{x + 1, y}
		q.Enqueue(p)
		tmp := []int{x + 1, y, (*board)[(dim*(x+1))+y].DotCount}
		*t = append(*t, tmp...)
	}
	if y < dim-1 {
		(*board)[(dim*x)+(y+1)].DotCount++
		(*board)[(dim*x)+(y+1)].Color = color
		p = utils.Pair{x, y + 1}
		q.Enqueue(p)
		tmp := []int{x, y + 1, (*board)[(dim*x)+(y+1)].DotCount}
		*t = append(*t, tmp...)
	}

	// seperator to know which states updated levelwise
	*t = append(*t, -1)
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

	states := updateBoard(&board, x, y, player.Color, gameInstance.Dimension)

	won := checkIfWon(gameInstance, player.Color)

	helpers.SetIfAllPlayedOnce(gameInstance, player.UserName)

	gameInstance.UpdatedBoard <- states

	if won {
		gameInstance.SetWinner(player)
	}

	return nil
}
