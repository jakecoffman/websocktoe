package tictactoe

import "testing"

func TestWinner(t *testing.T) {
	board := [3][3]string{}
	board[0] = [3]string{"a", "", ""}
	board[1] = [3]string{"",  "a", ""}
	board[2] = [3]string{"", "", "a"}
	if r := winner(board, 2, 2, "a"); r != "a" {
		t.Error("Unexpected", r)
	}
	board = [3][3]string{}
	board[0] = [3]string{"", "", "a"}
	board[1] = [3]string{"",  "a", ""}
	board[2] = [3]string{"a", "", ""}
	if r := winner(board, 0, 2, "a"); r != "a" {
		t.Error("Unexpected", r)
	}
}
