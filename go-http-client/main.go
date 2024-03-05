package main

import (
	"errors"
	"log"
	"sync"

	"github.com/0x9ef/clientx"
)

type PHPNoiseAPI struct {
	*clientx.API
	mu             *sync.Mutex
	lastUploadURI  string
	lastUploadSize int
}

func New(api *clientx.API) *PHPNoiseAPI {
	return &PHPNoiseAPI{
		API: api,
		mu:  new(sync.Mutex),
	}
}

type (
	GenerateRequest struct {
		R           int       `url:"r"`
		G           int       `url:"g"`
		B           int       `url:"b"`
		Titles      int       `url:"titles"`
		TitleSize   int       `url:"titleSize"`
		BorderWidth int       `url:"borderWidth"`
		ColorMode   ColorMode `url:"colorMode"`
		JSON        int       `url:"json"`
		Base64      int       `url:"base64"`
	}
	ColorMode string
	Generate  struct {
		URI    string `json:"uri`
		Base64 int    `json:"base64"`
	}
)

const (
	ColorModeBrightness ColorMode = "brightness"
	ColorModeAround     ColorMode = "around"
)

func (mode ColorMode) String() string {
	return string(mode)
}

func (r *GenerateRequest) Validate() error {
	if r.R > 255 {
		return errors.New("R is exceeded >255")
	}
	if r.G > 255 {
		return errors.New("G is exceeded >255")
	}
	if r.B > 255 {
		return errors.New("B is exceeded >255")
	}
	if r.Titles < 1 || r.Titles > 50 {
		return errors.New("invalid Tiles (1-50)")
	}
	if r.TitleSize < 1 || r.TitleSize > 20 {
		return errors.New("invalid TileSize (1-20)")
	}
	if r.BorderWidth > 15 {
		return errors.New("invalid BorderWidth (0-15)")
	}
	if r.ColorMode != ColorModeBrightness && r.ColorMode != ColorModeAround {
		return errors.New("invalid ColorMode, supported: brightness, around")
	}
	return nil
}

func main() {
	log.Println("very first line of Go http client")
}
