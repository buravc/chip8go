package main

import (
	"chip8go/machine"
	"chip8go/screen"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
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

	doneChan := make(chan bool, 1)

	core, err := machine.NewCore(romAsm, doneChan)
	if err != nil {
		panic(err)
	}

	scr := screen.Screen{
		VideoMem: core.GetVRAM(),
		KeyChan:  make(chan screen.KeyEvent),
	}

	go bindInput(core, scr.KeyChan, doneChan)

	go func() {
		for {
			if len(doneChan) > 0 {
				fmt.Println("core cycle loop done")
				return
			}
			time.Sleep(time.Millisecond * 2)
			core.Cycle()
		}
	}()

	scr.Init(doneChan)

	scr.MainLoop()
	log.Println("Main loop ended")
	for i := 0; i < 3; i++ {
		doneChan <- true
	}
	scr.Close()
}

func bindInput(core *machine.Chip8, keyChan chan screen.KeyEvent, doneChan chan bool) {
	for {

		select {
		case <-doneChan:
			log.Println("bindInput done")
			return
		case val := <-keyChan:
			core.SetKey(val.KeyCode, val.Pressed)
		}
	}
}
