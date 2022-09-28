package tictacgo

import "github.com/veandco/go-sdl2/sdl"

type Asset struct {
	texture  *sdl.Texture
	src, dest sdl.Rect
}

func (p_asset *Asset) Load(p_renderer *sdl.Renderer, p_path string) {
	surface, err := sdl.LoadBMP(p_path)
	if err != nil {
		sdl.Log("%s", sdl.GetError())
		panic(err)
	} else {
		sdl.Log("Loaded asset '%s'", p_path)
	}

	colorKey := sdl.MapRGB(surface.Format, 255, 0, 255)
	surface.SetColorKey(true, colorKey)

	texture, err := p_renderer.CreateTextureFromSurface(surface)
	if err != nil {
		sdl.Log("%s", sdl.GetError())
		panic(err)
	}

	p_asset.texture = texture
	p_asset.dest    = sdl.Rect{W: surface.W, H: surface.H}
	p_asset.src     = sdl.Rect{W: surface.W, H: surface.H}

	surface.Free()
}

func (p_asset *Asset) Render(p_renderer *sdl.Renderer) {
	p_renderer.Copy(p_asset.texture, &p_asset.src, &p_asset.dest)
}

func (p_asset *Asset) Free() {
	p_asset.texture.Destroy()
}
