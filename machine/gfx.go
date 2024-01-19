package machine

const width = 64
const height = 32

type gfx struct {
	vram *[256]byte
}

func (g *gfx) clear() {
	for i := range g.vram {
		g.vram[i] = 0
	}
}

// 1bpp depth
// each byte in the vram corresponds to 8 pixels
func (g *gfx) setPixel(x, y int, value byte) {
	//fmt.Printf("x:%d y:%d val:%d\n", x, y, value)
	bitOffset := y*width + x
	bytes := bitOffset / 8
	modBit := bitOffset % 8
	g.vram[bytes] = setBit(g.vram[bytes], value, 7-modBit)
}

func (g *gfx) getPixel(x, y int) byte {
	bitOffset := y*width + x
	bytes := bitOffset / 8
	modBit := bitOffset % 8

	return getBit(g.vram[bytes], 7-modBit)
}

func (g *gfx) setRow(x, y byte, value byte) {
	g.vram[y*height+x] = value
}

func (g *gfx) get8Pixel(x, y int) byte {
	bitOffset := y*width + x
	bytes := bitOffset / 8
	modBit := bitOffset % 8

	pBytes := g.vram[bytes]
	if modBit == 0 {
		return pBytes
	}

	sBytes := g.vram[bytes+1]

	pBytes <<= modBit
	sBytes >>= 8 - modBit
	return pBytes | sBytes
}
