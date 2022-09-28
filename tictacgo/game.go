package tictacgo

import "github.com/veandco/go-sdl2/sdl"

const (
	capFPS uint32 = 60

	scale int32 = 13

	scrWidth  int32 = 23
	scrHeight int32 = 30

	gameTitle string = "Tic Tac Go!"
)

const (
	assetBoard = iota
	assetPlayers
	assetWon
	assetTie
	assetsCount
)

type State int
const (
	stateInProgress = iota
	stateXWon
	stateOWon
	stateTie
)

type Game struct {
	window   *sdl.Window
	renderer *sdl.Renderer

	screen    *sdl.Texture
	screenRect sdl.Rect

	assets [assetsCount]Asset

	player Player
	board  Board

	cursor sdl.Point

	state State
	quit  bool
}

func (p_game *Game) LoadAsset(p_path string, p_idx int) {
	p_game.assets[p_idx].Load(p_game.renderer, p_path)
}

func (p_game *Game) Init() {
	err := sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		sdl.Log("%s", sdl.GetError())
		panic(err)
	}

	p_game.window, err = sdl.CreateWindow(gameTitle, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
	                                      scrWidth * scale, scrHeight * scale, sdl.WINDOW_RESIZABLE)
	if err != nil {
		sdl.Log("%s", sdl.GetError())
		panic(err)
	}

	p_game.renderer, err = sdl.CreateRenderer(p_game.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		sdl.Log("%s", sdl.GetError())
		panic(err)
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "nearest")
	err = p_game.renderer.SetLogicalSize(scrWidth, scrHeight)
	if err != nil {
		sdl.Log("%s", sdl.GetError())
		panic(err)
	}

	err = p_game.renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		sdl.Log("%s", sdl.GetError())
		panic(err)
	}

	p_game.screen, err = p_game.renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888,
	                                                   sdl.TEXTUREACCESS_TARGET,
	                                                   scrWidth, scrHeight)
	if err != nil {
		sdl.Log("%s", sdl.GetError())
		panic(err)
	}

	p_game.screenRect.W = scrWidth
	p_game.screenRect.H = scrHeight

	p_game.player = symbolX

	p_game.LoadAsset("./res/board.bmp",   assetBoard);
	p_game.LoadAsset("./res/players.bmp", assetPlayers);
	p_game.LoadAsset("./res/won.bmp",     assetWon);
	p_game.LoadAsset("./res/tie.bmp",     assetTie);

	p_game.assets[assetBoard].dest.Y = 7

	p_game.assets[assetPlayers].dest.W /= 2
	p_game.assets[assetPlayers].src.W  /= 2

	p_game.assets[assetWon].dest.X = p_game.assets[assetPlayers].dest.W
	p_game.assets[assetWon].dest.Y = 1
}

func (p_game *Game) Quit() {
	for _, asset := range p_game.assets {
		asset.Free()
	}

	p_game.window.Destroy()
	p_game.renderer.Destroy()

	sdl.Quit()
}

func (p_game *Game) Run() {
	for !p_game.quit {
		p_game.Render()
		p_game.Input()

		sdl.Delay(1000 / capFPS)
	}
}

func (p_game *Game) Render() {
	p_game.renderer.SetDrawColor(0, 0, 0, sdl.ALPHA_OPAQUE)
	p_game.renderer.Clear()

	p_game.renderer.SetDrawColor(42, 46, 121, sdl.ALPHA_OPAQUE)
	p_game.renderer.SetRenderTarget(p_game.screen)
	p_game.renderer.Clear()

	p_game.RenderBoard()

	if p_game.state == stateInProgress {
		p_game.RenderCursor()
	} else {
		p_game.RenderVictoryText()
	}

	p_game.renderer.SetRenderTarget(nil)
	p_game.renderer.Copy(p_game.screen, nil, &p_game.screenRect)
	p_game.renderer.Present()
}

func (p_game *Game) RenderCursor() {
	x := p_game.cursor.X * (p_game.assets[assetPlayers].dest.W + 1)
	y := p_game.cursor.Y * (p_game.assets[assetPlayers].dest.H + 1) +
	     p_game.assets[assetBoard].dest.Y

	p_game.assets[assetPlayers].texture.SetAlphaMod(128)

	switch p_game.player {
	case symbolX: p_game.RenderX(x, y); break
	case symbolO: p_game.RenderO(x, y); break

	default: panic("player is not X or O")
	}

	p_game.assets[assetPlayers].texture.SetAlphaMod(sdl.ALPHA_OPAQUE)
}

func (p_game *Game) RenderBoard() {
	p_game.assets[assetBoard].Render(p_game.renderer)

	for y, row := range p_game.board {
		for x, tile := range row {
			x := int32(x) * (p_game.assets[assetPlayers].dest.W + 1)
			y := int32(y) * (p_game.assets[assetPlayers].dest.H + 1) +
			     p_game.assets[assetBoard].dest.Y

			switch tile {
			case symbolNone: continue

			case symbolX: p_game.RenderX(x, y); break
			case symbolO: p_game.RenderO(x, y); break

			default: panic("Tile is not X, O or None")
			}
		}
	}
}

func (p_game *Game) RenderVictoryText() {
	switch p_game.state {
	case stateXWon: p_game.RenderX(0, 0); break
	case stateOWon: p_game.RenderO(0, 0); break

	case stateTie: p_game.assets[assetTie].Render(p_game.renderer); break

	default: panic("RenderVictoryText() called when state is not XWon, OWon or Tie")
	}

	p_game.assets[assetWon].Render(p_game.renderer)
}

func (p_game *Game) RenderX(p_x int32, p_y int32) {
	p_game.assets[assetPlayers].dest.X = p_x;
	p_game.assets[assetPlayers].dest.Y = p_y;

	p_game.assets[assetPlayers].src.X = 0;
	p_game.assets[assetPlayers].Render(p_game.renderer)
}

func (p_game *Game) RenderO(p_x int32, p_y int32) {
	p_game.assets[assetPlayers].dest.X = p_x;
	p_game.assets[assetPlayers].dest.Y = p_y;

	p_game.assets[assetPlayers].src.X = p_game.assets[assetPlayers].src.W;
	p_game.assets[assetPlayers].Render(p_game.renderer)
}

func (p_game *Game) Input() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			p_game.quit = true

			break

		case *sdl.MouseMotionEvent: p_game.InputMouseMotion(t.X, t.Y); break
		case *sdl.MouseButtonEvent:
			if t.State == sdl.PRESSED {
				p_game.InputMouseClick();
			}

			break

		case *sdl.KeyboardEvent:
			switch t.Keysym.Sym {
			case sdl.K_ESCAPE: p_game.quit = true; break
			}

			break
		}
	}
}

func (p_game *Game) InputMouseMotion(p_x int32, p_y int32) {
	p_game.cursor.X = int32(float32(p_x) /
	                  (float32(p_game.assets[assetBoard].dest.W) / 3.0))
	p_game.cursor.Y = int32(float32(p_y - p_game.assets[assetBoard].dest.Y) /
	                  (float32(p_game.assets[assetBoard].dest.H) / 3.0))

	if p_game.cursor.X < 0 {
		p_game.cursor.X = 0
	} else if p_game.cursor.X >= boardSize - 1 {
		p_game.cursor.X = boardSize - 1
	}

	if p_game.cursor.Y < 0 {
		p_game.cursor.Y = 0
	} else if p_game.cursor.Y >= boardSize - 1 {
		p_game.cursor.Y = boardSize - 1
	}
}

func (p_game *Game) InputMouseClick() {
	if p_game.state != stateInProgress {
		p_game.board.Clear()
		p_game.state = stateInProgress
	} else if p_game.board.Place(p_game.cursor.X, p_game.cursor.Y, Symbol(p_game.player)) {
		won, symbol := p_game.board.CheckState(Symbol(p_game.player))
		if won {
			switch symbol {
			case symbolX:    p_game.state = stateXWon; break
			case symbolO:    p_game.state = stateOWon; break
			case symbolNone: p_game.state = stateTie;  break

			default: panic("symbol is not X, O or None")
			}

			p_game.player = symbolX
		} else {
			p_game.player.Switch()
		}
	}
}
