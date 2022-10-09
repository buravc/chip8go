package machine

import "testing"

func Test_byteArrToUint16(t *testing.T) {
	type input struct {
		high, low byte
	}
	type expect struct {
		want uint16
	}
	type test struct {
		name string
		input
		expect
	}
	tests := []test{
		{
			"0b1111111100101010",
			input{high: 0b11111111, low: 0b00101010},
			expect{want: 0b1111111100101010},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := byteArrToUint16(tt.high, tt.low)
			if got != tt.want {
				t.Errorf("expect: %08b, got: %08b", tt.expect.want, got)
			}
		})
	}
}

func Test_getNibble(t *testing.T) {
	type input struct {
		number uint16
		index  int
	}
	type expect struct {
		want byte
	}
	type test struct {
		name string
		input
		expect
	}
	tests := []test{
		{
			"0b0101011100011111[3]",
			input{number: 0b0101011100011111, index: 3},
			expect{want: 0b0101},
		},
		{
			"0b0101011100011111[2]",
			input{number: 0b0101011100011111, index: 2},
			expect{want: 0b0111},
		},
		{
			"0b0101011100011111[1]",
			input{number: 0b0101011100011111, index: 1},
			expect{want: 0b0001},
		},
		{
			"0b0101011100011111[0]",
			input{number: 0b0101011100011111, index: 0},
			expect{want: 0b1111},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getNibble(tt.number, tt.index)
			if got != tt.want {
				t.Errorf("expect: %08b, got: %08b", tt.expect.want, got)
			}
		})
	}
}

func Test_getByte(t *testing.T) {
	type input struct {
		number uint16
		index  int
	}
	type expect struct {
		want byte
	}
	type test struct {
		name string
		input
		expect
	}
	tests := []test{
		{
			"0b0101011100011111[1]",
			input{number: 0b0101011100011111, index: 1},
			expect{want: 0b01010111},
		},
		{
			"0b0101011100011111[0]",
			input{number: 0b0101011100011111, index: 0},
			expect{want: 0b00011111},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getByte(tt.number, tt.index)
			if got != tt.want {
				t.Errorf("expect: %08b, got: %08b", tt.expect.want, got)
			}
		})
	}
}
