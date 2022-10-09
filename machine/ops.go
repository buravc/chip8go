package machine

import (
	"math/rand"
)

func notImplemented() {
	panic("not implemented instruction called")
}

func (core *Chip8) advance() {
	core.pc += 2
}

// 0x00E0
func (core *Chip8) clearDisplay() {
	core.gfx.clear()
	core.advance()
}

// 0x00EE
func (core *Chip8) subReturn() {
	core.pc = core.pop()
}

// 0x0NNN
func (core *Chip8) call() {
	notImplemented()
}

// 0x1NNN
// 0xBNNN
func (core *Chip8) jmp(msn byte) {
	core.pc = core.cycleInfo.nnn
	switch msn {
	case 0x1:
		break
	case 0xB:
		core.pc += uint16(core.v[0])
		break
	default:
		notImplemented()
	}
}

// 0x2NNN
func (core *Chip8) sub() {
	core.push(core.pc + 2)
	core.pc = core.cycleInfo.nnn
}

// 0x3XNN
// 0x4XNN
// 0x5XY0
// 0x9XY0
// 0xEX9E
// 0xEXA1
func (core *Chip8) skip(msn byte) {
	switch true {
	case msn == 0x3:
		if core.getVX() == core.cycleInfo.nn {
			core.advance()
		}
		break
	case msn == 0x4:
		if core.getVX() != core.cycleInfo.nn {
			core.advance()
		}
		break
	case msn == 0x5 && core.cycleInfo.n == 0:
		if core.getVX() == core.getVY() {
			core.advance()
		}
		break
	case msn == 0x9 && core.cycleInfo.n == 0:
		if core.getVX() != core.getVY() {
			core.advance()
		}
		break
	case msn == 0xE && core.cycleInfo.nn == 0x9E:
		if core.getKey(core.getVX()) == 1 {
			core.advance()
		}
		break
	case msn == 0xE && core.cycleInfo.nn == 0xA1:
		if core.getKey(core.getVX()) == 0 {
			core.advance()
		}
		break
	default:
		notImplemented()
	}
	core.advance()
}

// 0x6XNN
// 0x7XNN
func (core *Chip8) setXNN(msn byte) {
	switch msn {
	case 0x6:
		core.setVX(core.cycleInfo.nn)
		break
	case 0x7:
		core.setVX(core.getVX() + core.cycleInfo.nn)
		break
	default:
		notImplemented()
	}
	core.advance()
}

// 0x8XY0...E
func (core *Chip8) setXY() {
	switch core.cycleInfo.n {
	case 0x0:
		core.setVX(core.getVY())
		break
	case 0x1:
		core.setVX(core.getVX() | core.getVY())
		break
	case 0x2:
		core.setVX(core.getVX() & core.getVY())
		break
	case 0x3:
		core.setVX(core.getVX() ^ core.getVY())
		break
	case 0x4:
		setCarry := byte(0)
		if (255 - core.getVX()) < core.getVY() {
			setCarry = 1
		} else {
			setCarry = 0
		}
		core.setVX(core.getVX() + core.getVY())
		core.v[0xF] = setCarry
		break
	case 0x5:
		setCarry := byte(0)
		if core.getVX() < core.getVY() {
			setCarry = 0
		} else {
			setCarry = 1
		}
		core.setVX(core.getVX() - core.getVY())
		core.v[0xF] = setCarry
		break
	case 0x6:
		lsb := getBit(core.getVX(), 0)
		core.setVX(core.getVX() >> 1)
		core.v[0xF] = lsb
		break
	case 0x7:
		setCarry := byte(0)
		if core.getVY() < core.getVX() {
			setCarry = 0
		} else {
			setCarry = 1
		}
		core.setVX(core.getVY() - core.getVX())
		core.v[0xF] = setCarry
		break
	case 0xE:
		msb := getBit(core.getVX(), 7)
		core.setVX(core.getVX() << 1)
		core.v[0xF] = msb
		break
	}
	core.advance()
}

// 0xANNN
func (core *Chip8) setINNN() {
	core.Registers.i = core.cycleInfo.nnn
	core.advance()
}

// 0xCXNN
func (core *Chip8) rand() {
	rnd := rand.Intn(256)
	res := byte(rnd) & core.cycleInfo.nn
	core.setVX(res)
	core.advance()
}

// 0xDXYN
func (core *Chip8) draw() {
	core.v[0xF] = 0
	vx := core.getVX() % 64
	vy := core.getVY() % 32
	for i := 0; i < int(core.cycleInfo.n); i++ {
		spriteRow := core.ReadByte(core.i + uint16(i))
		//fmt.Printf("spriteRow: %08b\n", spriteRow)
		for j := 0; j < 8; j++ {
			x := int(vx) + j
			y := int(vy) + i
			if x < 64 && y < 32 {
				spriteBit := getBit(spriteRow, 7-j)
				screenBit := core.gfx.getPixel(x, y)

				if screenBit == 1 && spriteBit == 1 {
					core.v[0xF] = 1
				}

				core.gfx.setPixel(x, y, screenBit^spriteBit)
			}
		}
		//for k := 0; k < 32; k++ {
		//	fmt.Printf("%08b", core.gfx.vram[k*8:k*8+8])
		//	fmt.Println()
		//}
	}
	core.advance()
}

// 0xFX07
// 0xFX0A
func (core *Chip8) setX() {
	switch core.cycleInfo.nn {
	case 0x07:
		core.setVX(core.delayTimer)
		break
	case 0x0A:
		core.setVX(core.getAnyKey())
		break
	default:
		notImplemented()
	}
	core.advance()
}

// 0xFX15
// 0xFX18
func (core *Chip8) setTimer() {
	switch core.cycleInfo.nn {
	case 0x15:
		core.delayTimer = core.getVX()
		break
	case 0x18:
		core.soundTimer = core.getVX()
		break
	default:
		notImplemented()
	}

	core.advance()
}

// 0xFX1E
// 0xFX29
// 0xFX33
func (core *Chip8) setI() {
	switch core.cycleInfo.nn {
	case 0x1E:
		core.i += uint16(core.getVX())
		break
	case 0x29:
		core.i = uint16(core.getVX()) * 5
		break
	case 0x33:
		x := core.getVX()
		hund := x / 100
		ten := (x % 100) / 10
		one := (x % 100) % 10

		core.memory[core.i] = hund
		core.memory[core.i+1] = ten
		core.memory[core.i+2] = one
		break
	default:
		notImplemented()
	}

	core.advance()
}

// 0xFX55
// 0xFX65
func (core *Chip8) reg() {
	switch core.cycleInfo.nn {
	case 0x55:
		for i := uint16(0); i <= uint16(core.cycleInfo.x); i++ {
			core.memory[core.i+i] = core.v[i]
		}
		break
	case 0x65:
		for i := uint16(0); i <= uint16(core.cycleInfo.x); i++ {
			core.v[i] = core.memory[core.i+i]
		}
		break
	default:
		notImplemented()
	}

	core.advance()
}

////////////////////////////////////////
// to make it easy to get and set vx and vy registers

func (core *Chip8) getVX() byte {
	return core.v[core.cycleInfo.x]
}

func (core *Chip8) setVX(val byte) {
	core.v[core.cycleInfo.x] = val
}

func (core *Chip8) getVY() byte {
	return core.v[core.cycleInfo.y]
}

func (core *Chip8) setVY(val byte) {
	core.v[core.cycleInfo.y] = val
}
