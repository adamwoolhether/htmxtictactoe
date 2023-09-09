package game

import (
	"errors"
	"fmt"
)

type Player string

const (
	X Player = "X"
	O Player = "0"

	BoardSize = 3
)

type Game interface {
	// Start a game with fresh state.
	Start()
	// Move the player's symbol to the given row/column, returning an error if invalid.
	Move(player Player, row, column int) error
	// GetTurn returns which player whose turn it is.
	GetTurn() Player
	// GetBoard returns the current tic-tac-toe board.
	GetBoard() [BoardSize][BoardSize]*Player
	// GetWinner returns the winner of the tic-tac-toe game.
	GetWinner() *Player
}

type TicTacToe struct {
	turn   Player
	winner *Player
	board  [BoardSize][BoardSize]*Player
}

func New() *TicTacToe {
	ttt := TicTacToe{
		turn:   X,
		winner: nil,
		board:  [BoardSize][BoardSize]*Player{},
	}

	return &ttt
}

func (t *TicTacToe) Start() {
	t.turn = X
	t.winner = nil
	t.board = [BoardSize][BoardSize]*Player{}
}

func (t *TicTacToe) Move(player Player, row, column int) error {
	if t.winner != nil {
		return errors.New("game is already over")
	}

	if player != t.turn {
		return fmt.Errorf("it's not %s's turn", player)
	}

	if !isValidMove(t.board, row, column) {
		return fmt.Errorf("location %d,%d is not empty or out of bound", row, column)
	}

	t.board[row][column] = &player

	t.winner = getWinner(t.board)
	if t.winner == nil {
		t.turn = switchPlayer(t.turn)
	}

	return nil
}

func (t *TicTacToe) GetTurn() Player {
	return t.turn
}

func (t *TicTacToe) GetBoard() [BoardSize][BoardSize]*Player {
	return t.board
}

func (t *TicTacToe) GetWinner() *Player {
	return t.winner
}

func isValidMove(board [BoardSize][BoardSize]*Player, row, column int) bool {
	return row >= 0 && row < BoardSize && column >= 0 && column < BoardSize && board[row][column] == nil
}

// This code sucks
func getWinner(board [BoardSize][BoardSize]*Player) *Player {
	winConditions := [][3][2]int{
		{{0, 0}, {0, 1}, {0, 2}},
		{{1, 0}, {1, 1}, {1, 2}},
		{{2, 0}, {2, 1}, {2, 2}},
		{{0, 0}, {1, 0}, {2, 0}},
		{{0, 1}, {1, 1}, {2, 1}},
		{{0, 2}, {1, 2}, {2, 2}},
		{{0, 0}, {1, 1}, {2, 2}},
		{{0, 2}, {1, 1}, {2, 0}},
	}

	for _, line := range winConditions {
		a, b, c := line[0], line[1], line[2]
		if board[a[0]][a[1]] != nil && board[b[0]][b[1]] != nil && board[c[0]][c[1]] != nil {
			if *board[a[0]][a[1]] == *board[b[0]][b[1]] && *board[a[0]][a[1]] == *board[c[0]][c[1]] {
				return board[a[0]][a[1]]
			}
		}
	}

	return nil
}

func switchPlayer(currentPlayer Player) Player {
	if currentPlayer == X {
		return O
	}

	return X
}
