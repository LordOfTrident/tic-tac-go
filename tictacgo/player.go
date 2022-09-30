package tictacgo

type Player Symbol

func (p_player *Player) Switch() {
	switch *p_player {
	case symbolX: *p_player = symbolO
	case symbolO: *p_player = symbolX

	default: panic("player is not X or O")
	}
}
