package machine

import "testing"

func Test_getPixel(t *testing.T) {
	gfx := gfx{vram: [256]byte{}}

	for i := 0; i < 32; i++ {
		for j := 0; j < 64; j++ {
			if gfx.getPixel(j, i) != 0 {
				t.Fail()
			}

			gfx.setPixel(j, i, 1)

			if gfx.getPixel(j, i) != 1 {
				t.Fail()
			}
		}
	}
}
