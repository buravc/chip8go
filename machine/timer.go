package machine

import (
	"fmt"
	"time"
)

const msPerFrame = 16666

func (core *Chip8) initTimer() {
	ticker := time.NewTicker(time.Duration(msPerFrame) * time.Microsecond)

	go func() {
		prevSoundTimer := byte(0)

		for {
			if len(core.DoneChan) > 0 {
				fmt.Println("timer loop done")
			}
			select {
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
