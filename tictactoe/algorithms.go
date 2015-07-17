package tictactoe

func winner(board [3][3]string, x, y int, name string) string {
	for i := 0; i < 3; i++ {
		if board[x][i] != name {
			break
		}
		if i == 2 {
			return name
		}
	}

	for i := 0; i < 3; i++ {
		if board[i][y] != name {
			break
		}
		if i == 2 {
			return name
		}
	}

	for i := 0; i < 3; i++ {
		if board[i][i] != name {
			break
		}
		if i == 2 {
			return name
		}
	}

	for i := 0; i < 3; i++ {
		if board[i][2-i] != name {
			break
		}
		if i == 2 {
			return name
		}
	}

	return ""
}
