# chip8go
Golang and SDL implementation of Chip8 VM. Passes all tests of [Timendus' Chip-8 Test Suite](https://github.com/Timendus/chip8-test-suite).

# Build
```shell
$ go build -ldflags='-s -w' -o ./chip8go ./cmd/main.go
```

# Run program
```shell
$ GODEBUG=cgocheck=0 chip8go -r=<path-to-rom-file>
```

Go runtime doesn't allow managed memory to be passed to C functions, since `1.16` version. `GODEBUG=cgocheck=0` is required to disable this check.

# Keymap
. | . | . | .
--- | --- | --- | ---
1 | 2 | 3 | 4
Q | W | E | T
A | S | D | F
Z | X | C | V
