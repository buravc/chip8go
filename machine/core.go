package machine

import (
	"errors"
	"fmt"
	"unsafe"
)

const (
	programStart   = 0x200
	stackStart     = 0xEA0
	gfxStart       = 0xF00
	gfxSize        = 0xFF // 256
	maxProgramSize = stackStart - programStart
)

var fonts = [80]byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

type Chip8 struct {
	memory     [4096]byte
	delayTimer byte
	soundTimer byte
	Registers
	gfx       gfx
	cycleInfo CycleInfo
	Input
	DoneChan chan bool
}

type Registers struct {
	pc uint16
	sp uint16
	i  uint16
	v  [16]byte
}

type CycleInfo struct {
	opcode, nnn uint16
	x, y, n, nn byte
}

func (c CycleInfo) String() string {
	return fmt.Sprintf("{opcode:%x x:%x y:%x nnn:%x nn:%x n:%x}", c.opcode, c.x, c.y, c.nnn, c.nn, c.n)
}

func (r Registers) String() string {
	return fmt.Sprintf("{pc:%x sp:%x i:%x v[0]:%x v[1]:%x v[2]:%x v[3]:%x v[4]:%x v[5]:%x v[6]:%x v[7]:%x v[8]:%x v[9]:%x v[A]:%x v[B]:%x v[C]:%x v[D]:%x v[E]:%x v[F]:%x}", r.pc, r.sp, r.i, r.v[0], r.v[1], r.v[2], r.v[3], r.v[4], r.v[5], r.v[6], r.v[7], r.v[8], r.v[9], r.v[10], r.v[11], r.v[12], r.v[13], r.v[14], r.v[15])
}

func NewCore(program []byte, doneChan chan bool) (*Chip8, error) {
	length := len(program)
	if length > maxProgramSize {
		return nil, errors.New("invalid program size")
	}

	if doneChan == nil {
		return nil, errors.New("done channel is nil")
	}

	core := &Chip8{
		DoneChan: doneChan,
	}

	core.Input = &Keypad{Array: [2]byte{}}

	for i := 0; i < len(fonts); i++ {
		core.memory[i] = fonts[i]
	}

	for i := 0; i < length; i++ {
		core.memory[programStart+i] = program[i]
	}

	core.pc = programStart
	core.sp = stackStart

	core.gfx = gfx{vram: *(*[256]byte)(unsafe.Pointer(&core.memory[gfxStart]))}
	core.initTimer()

	return core, nil
}

func (core *Chip8) Cycle() {
	core.cycleInfo.opcode = core.ReadUint16(core.pc)
	core.cycleInfo.x = getNibble(core.cycleInfo.opcode, 2)
	core.cycleInfo.y = getNibble(core.cycleInfo.opcode, 1)
	core.cycleInfo.n = getNibble(core.cycleInfo.opcode, 0)
	core.cycleInfo.nn = getByte(core.cycleInfo.opcode, 0)
	core.cycleInfo.nnn = core.cycleInfo.opcode & 0x0FFF

	opNibble := getNibble(core.cycleInfo.opcode, 3)

	//fmt.Printf("%+v\n", core.cycleInfo)
	//fmt.Printf("%+v\n", core.Registers)
	//fmt.Printf("%08b\n", core.GetKeypadState())
	//fmt.Println()
	//fmt.Println()

	switch opNibble {
	case 0x0:
		switch core.cycleInfo.nnn {
		case 0xE0:
			core.clearDisplay()
			break
		case 0xEE:
			core.subReturn()
			break
		default:
			core.call()
			break
		}
		break
	case 0x1, 0xB:
		core.jmp(opNibble)
		break
	case 0x2:
		core.sub()
		break
	case 0x3, 0x4, 0x5, 0x9, 0xE:
		core.skip(opNibble)
		break
	case 0x6, 0x7:
		core.setXNN(opNibble)
		break
	case 0x8:
		core.setXY()
		break
	case 0xA:
		core.setINNN()
		break
	case 0xC:
		core.rand()
		break
	case 0xD:
		core.draw()
		break
	case 0xF:
		switch core.cycleInfo.nn {
		case 0x07, 0x0A:
			core.setX()
			break
		case 0x15, 0x18:
			core.setTimer()
			break
		case 0x1E, 0x29, 0x33:
			core.setI()
			break
		case 0x55, 0x65:
			core.reg()
			break
		default:
			notImplemented()
			break
		}
		break
	default:
		notImplemented()
		break
	}
}

func (core *Chip8) GetVRAM() *[256]byte {
	return &core.gfx.vram
}

func (core *Chip8) ReadUint16(address uint16) uint16 {
	return byteArrToUint16(core.memory[address], core.memory[address+1])
}

func (core *Chip8) ReadByte(address uint16) byte {
	return core.memory[address]
}

func (core *Chip8) pop() uint16 {
	core.sp++
	highest := core.memory[core.sp]

	core.sp++
	lowest := core.memory[core.sp]

	return byteArrToUint16(highest, lowest)
}

func (core *Chip8) push(b uint16) {
	lowest := getByte(b, 0)
	core.memory[core.sp] = lowest
	core.sp--

	highest := getByte(b, 1)
	core.memory[core.sp] = highest
	core.sp--
}
