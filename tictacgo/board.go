package tictacgo

const boardSize int32 = 3

type Symbol int
const (
	symbolNone = iota
	symbolX
	symbolO
)

type Board [boardSize][boardSize]Symbol

func (p_board *Board) Clear() {
	for y, row := range p_board {
		for x, _ := range row {
			p_board[y][x] = symbolNone
		}
	}
}

func (p_board *Board) Place(p_x, p_y int32, p_symbol Symbol) bool {
	if p_board[p_y][p_x] != symbolNone {
		return false
	} else {
		p_board[p_y][p_x] = p_symbol

		return true
	}
}

func (p_board *Board) CheckState(p_symbol Symbol) (bool, Symbol) {
	var horiz, vert, diag_left, diag_right int32
	hasEmptyTile := false

	for y, row := range p_board {
		for x, tile := range row {
			if !hasEmptyTile && tile == symbolNone {
				hasEmptyTile = true
			}

			if tile == p_symbol {
				horiz ++
			}

			if p_board[x][y] == p_symbol {
				vert ++
			}
		}

		if horiz >= boardSize || vert >= boardSize {
			return true, p_symbol
		} else {
			horiz = 0
			vert  = 0
		}

		if p_board[y][y] == p_symbol {
			diag_left ++
		}

		if p_board[y][int(boardSize) - y - 1] == p_symbol {
			diag_right ++
		}
	}

	if diag_left >= boardSize || diag_right >= boardSize {
		return true, p_symbol
	} else if !hasEmptyTile {
		return true, symbolNone
	} else {
		return false, p_symbol
	}
}
