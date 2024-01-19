package machine

type KeyEvent struct {
	Pressed byte
	KeyCode byte
}

type Input interface {
	getAnyKey() byte
	getKey(index byte) byte
	SetKey(index, pressed byte)
	GetKeypadState() [2]byte
}

type Keypad struct {
	Array [2]byte
}

func (k *Keypad) GetKeypadState() [2]byte {
	return k.Array
}

func (k *Keypad) getAnyKey() byte {
	for {
		for i := 0; i < 16; i++ {
			bit := getBit(k.Array[i/8], i)
			if bit == 1 {
				return 1
			}
		}
	}
}

func (k *Keypad) getKey(index byte) byte {
	bit := getBit(k.Array[index/8], 7-int(index%8))
	return bit
}

func (k *Keypad) SetKey(index, pressed byte) {
	arrayIndex := index / 8
	k.Array[arrayIndex] = setBit(k.Array[arrayIndex], pressed, 7-int(index%8))
}
