package main

import "github.com/LordOfTrident/tic-tac-go/tictacgo"

func main() {
	var game tictacgo.Game

	game.Init()
	defer game.Quit()

	game.Run()
}
