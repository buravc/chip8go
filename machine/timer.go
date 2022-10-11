package machine

import (
	"fmt"
	"log"
	"time"
)

const msPerFrame = 16666

func (core *Chip8) initTimer() {
	ticker := time.NewTicker(time.Duration(msPerFrame) * time.Microsecond)

	go func() {
		prevSoundTimer := byte(0)

		for {
			select {
			case <-core.DoneChan:
				log.Println("timer loop done")
				return
			case <-ticker.C:
				if core.delayTimer > 0 {
					core.delayTimer--
				}
				if core.soundTimer > 0 {
					if prevSoundTimer == 0 {
						fmt.Println("Beep")
						print("\a")
					}
					core.soundTimer--
				}
			}
			prevSoundTimer = core.soundTimer
		}
	}()
}
