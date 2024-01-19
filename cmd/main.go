package main

import (
	"chip8go/machine"
	"chip8go/screen"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var romFile string
	flag.StringVar(&romFile, "r", "", "path to rom file")

	flag.Parse()

	romFile, err := filepath.Abs(romFile)
	if err != nil {
		os.Exit(1)
	}

	log.Printf("trying to load %s\n", romFile)
	rom, err := os.Open(romFile)
	if err != nil {
		os.Exit(1)
	}

	romAsm, err := ioutil.ReadAll(rom)

	if err != nil {
		os.Exit(1)
	}

	core, err := machine.NewCore(romAsm)
	if err != nil {
		panic(err)
	}

	scr := screen.Screen{
		VideoMem: core.GetVRAM(),
		KeyChan:  core.InputChan,
	}

	core.Run()

	scr.Init()

	scr.MainLoop()
	log.Println("Main loop ended")
	scr.Close()
	core.Shutdown()
}
