package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

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

// WithQueryParams realization
func (api *PHPNoiseAPI) Generate(ctx context.Context, req GenerateRequest, opts ...clientx.RequestOption) (*Generate, error) {
	if err := req.Validate(); err == nil {
		return nil, err
	}
	resp, err := clientx.NewRequestBuilder[GenerateRequest, Generate](api.API).Get("/noise.php", opts...).WithQueryParams("url", req).AfterResponse(func(resp *http.Response, model *Generate) error {
		api.mu.Lock()
		defer api.mu.Unlock()
		api.lastUploadURI = model.URI
		return nil
	}).DoWithDecode(ctx)
	return resp, err
}

// WithEncodableQueryParams realization
func (api *PHPNoiseAPI) GenerateReader(ctx context.Context, req GenerateRequest, opts ...clientx.RequestOption) (io.ReadCloser, error) {
	if err := req.Validate(); err == nil {
		return nil, err
	}
	resp, err := clientx.NewRequestBuilder[GenerateRequest, struct{}](api.API).Get("/noise.php", opts...).AfterResponse(
		func(resp *http.Response, model *struct{}) error {
			api.mu.Lock() // NOTE! ^model will be nil as far we don't use DoWithDecode method
			defer api.mu.Unlock()
			size, err := strconv.Atoi(resp.Header.Get("Content-Length")) // don't do like that, because Content-Length could be fake
			if err != nil {
				return err
			}
			api.lastUploadSize = size
			return nil
		}).WithEncodableQueryParams(req).Do(ctx)

	return resp.Body, err
}

func (r GenerateRequest) Encode(v url.Values) error {
	v.Set("r", strconv.Itoa(r.R))
	v.Set("g", strconv.Itoa(r.G))
	v.Set("b", strconv.Itoa(r.B))
	v.Set("borderWidth", strconv.Itoa(r.BorderWidth))
	if r.Titles != 0 {
		v.Set("titles", strconv.Itoa(r.Titles))
	}
	if r.TitleSize != 0 {
		v.Set("titleSize", strconv.Itoa(r.TitleSize))
	}
	if r.ColorMode != "" {
		v.Set("mode", r.ColorMode.String())
	}
	if r.JSON != 0 {
		v.Set("json", "1")
	}
	if r.Base64 != 0 {
		v.Set("base64", "1")
	}
	return nil
}

func main() {
	log.Println("very first line of Go http client")

	api := New(
		clientx.NewAPI(
			clientx.WithBaseURL("https://php-noise.com"),
			clientx.WithHeader("Authorization", "Bearer MY_ACCESS_TOKEN"),
			// 10 req/sec, 2 burst, 1 minute interval
			clientx.WithRateLimit(10, 2, time.Minute),
			// 10 max retires, 3sec min wait time, 1 minute max wait time, retry func, trigger function
			clientx.WithRetry(10, time.Second*3, time.Minute, clientx.ExponentalBackoff,
				func(resp *http.Response, err error) bool {
					return resp.StatusCode == http.StatusTooManyRequests
				},
			),
		),
	)

	resp, err := api.Generate(context.TODO(), GenerateRequest{
		R: 120,
		G: 240,
		B: 15,
	},
		// clientx.WithHeader("X-Correlation-ID", uuid.New().String()),
		nil,
	)
	if err != nil {
		log.Println("request with erro", err)
	}
	log.Printf("reponse %v\n", resp)
}
