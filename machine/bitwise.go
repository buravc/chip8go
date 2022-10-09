package machine

func byteArrToUint16(high, low byte) uint16 {
	var num uint16

	num = num | uint16(high)<<8
	num = num | uint16(low)

	return num
}

func getByte(val uint16, index int) byte {
	return byte(maskUint16(val, 0x00FF, index*8))
}

func getNibble(val uint16, index int) byte {
	return byte(maskUint16(val, 0b0000000000001111, index*4))
}

func getBit(val byte, index int) byte {
	return maskByte(val, 0b00000001, index)
}

func setBit(val, bit byte, index int) byte {
	if bit == 1 {
		bit <<= index
		return val | bit
	}
	if bit == 0 {
		bit = 0xFF & ^(1 << index)
		return val & bit
	}
	notImplemented()
	return 0
}

func maskByte(val, mask byte, shift int) byte {
	mask = mask << shift
	return (val & mask) >> (shift)
}

func maskUint16(val, mask uint16, shift int) uint16 {
	mask = mask << shift
	return (val & mask) >> (shift)
}
