package machine

import "testing"

func Test_draw(t *testing.T) {
	core := Chip8{gfx: gfx{vram: make([]byte, 64*32)}}

	core.i = 0x28
	core.memory[0] = byte(0xD1)
	core.memory[1] = byte(0x21)
	core.v[1] = 1
	core.v[2] = 2
	core.memory[core.i] = 0b11110001
	core.Cycle()

	if core.v[0xF] != 1 {
		t.Fail()
	}

	pixel8 := core.gfx.get8Pixel(1, 2)
	coreMemI := core.memory[core.i]
	if pixel8 != coreMemI {
		t.Fail()
	}
}
