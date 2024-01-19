package screen

import (
	"chip8go/machine"
	"log"
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const usPerFrame = 16666

var keyMap = map[sdl.Keycode]byte{
	sdl.K_1: 0x0,
	sdl.K_2: 0x1,
	sdl.K_3: 0x2,
	sdl.K_4: 0x3,

	sdl.K_q: 0x4,
	sdl.K_w: 0x5,
	sdl.K_e: 0x6,
	sdl.K_r: 0x7,

	sdl.K_a: 0x8,
	sdl.K_s: 0x9,
	sdl.K_d: 0xA,
	sdl.K_f: 0xB,

	sdl.K_z: 0xC,
	sdl.K_x: 0xD,
	sdl.K_c: 0xE,
	sdl.K_v: 0xF,
}

type Screen struct {
	VideoMem *[256]byte
	KeyChan  chan machine.KeyEvent

	doneChan chan bool
	window   *sdl.Window
}

func (s *Screen) Init() {

	log.Println("Initing sdl")

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	log.Println("Initiated sdl")

	window, err := sdl.CreateWindow("Chip8 Machine", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 512, 256, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	s.window = window

	log.Println("Created window")

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	err = surface.FillRect(nil, 0)
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(time.Duration(usPerFrame) * time.Microsecond)

	palette, err := sdl.AllocPalette(2)
	if err != nil {
		panic(err)
	}

	err = palette.SetColors([]sdl.Color{{A: 255}, {R: 255, G: 255, B: 255, A: 255}})
	if err != nil {
		panic(err)
	}

	s.doneChan = make(chan bool)

	go func() {
		for {
			select {
			case <-s.doneChan:
				log.Println("draw loop done")
				return
			case <-ticker.C:
				oneBitPPSurface, _ := sdl.CreateRGBSurfaceWithFormatFrom(unsafe.Pointer(s.VideoMem), 64, 32, 1, 8, sdl.PIXELFORMAT_INDEX1MSB)

				err = oneBitPPSurface.SetPalette(palette)
				if err != nil {
					panic(err)
				}

				newSurf, err := sdl.CreateRGBSurfaceWithFormat(0, 64, 32, 1, sdl.PIXELFORMAT_ARGB8888)
				if err != nil {
					panic(err)
				}

				err = oneBitPPSurface.Blit(nil, newSurf, &sdl.Rect{X: 0, Y: 0, W: 64, H: 32})
				if err != nil {
					panic(err)
				}

				err = newSurf.BlitScaled(nil, surface, &sdl.Rect{X: 0, Y: 0, W: 512, H: 256})
				if err != nil {
					panic(err)
				}

				err = window.UpdateSurface()
				if err != nil {
					panic(err)
				}
			}
		}
	}()
}

func (s *Screen) Close() {
	if s.window != nil {
		s.doneChan <- true
		err := s.window.Destroy()
		if err != nil {
			log.Println(err)
		}
	}
	sdl.Quit()
}

func (s *Screen) MainLoop() {
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch sdlEvent := event.(type) {
			case *sdl.QuitEvent:
				log.Println("Quit Event received")
				return
			case *sdl.KeyboardEvent:
				if val, ok := keyMap[sdlEvent.Keysym.Sym]; ok {
					s.KeyChan <- machine.KeyEvent{
						Pressed: sdlEvent.State,
						KeyCode: val,
					}
				}
			}
		}
	}
}
