# chip8go
Golang and SDL implementation of Chip8 VM. Passes all tests of [Timendus' Chip-8 Test Suite](https://github.com/Timendus/chip8-test-suite).

# Build
```shell
$ go build -ldflags='-s -w' -o ./chip8go ./cmd/main.go
```

# Run program
```shell
$ ./chip8go -r=<path-to-rom-file>
```

# Keymap
. | . | . | .
--- | --- | --- | ---
1 | 2 | 3 | 4
Q | W | E | T
A | S | D | F
Z | X | C | V
