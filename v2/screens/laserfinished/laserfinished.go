package laserfinished

import (
	"fmt"
	"github.com/b4ckspace/ledboard-v2/ledboard"
)

// GenerateLaserFinishedCommand generates the command string for the laser finished screen.
func GenerateLaserFinishedCommand(duration int) string {
	var cmd string

	cmd += ledboard.ControlPatternIn + ledboard.PatternScrollUp
	cmd += ledboard.ControlPatternOut + ledboard.PatternScrollUp

	if duration > 10 * 60 {
		cmd += ledboard.FontNormal14x8
		cmd += ledboard.ControlFontColor + ledboard.FontColorGreen

		// Blinkin' for the poor (because FLASH seems buggy)
		for i := 0; i < 3; i++ {
			cmd += "Congratulations!"
			cmd += ledboard.PauseMillisecond4 + "0400"
			cmd += ledboard.ControlFrame
			cmd += " "
			cmd += ledboard.PauseMillisecond4 + "0100"
			cmd += ledboard.ControlFrame
		}
	}

	cmd += ledboard.ControlPatternIn + ledboard.PatternPeelOffL // Original was PEEL_OFF_R, but that's not defined. Assuming PEEL_OFF_L

	cmd += ledboard.FontNormal7x6
	cmd += ledboard.ControlFontColor + ledboard.FontColorGreen
	cmd += ledboard.ControlFlash + ledboard.FlashOff

	cmd += "Laser-Job finished:"

	cmd += ledboard.ControlLineFeed
	cmd += ledboard.ControlFontColor + ledboard.FontColorRed

	hours := duration / 3600
	minutes := (duration % 3600) / 60
	seconds := duration % 60

	if hours > 0 {
		cmd += fmt.Sprintf("%dh ", hours)
	}

	if minutes > 0 {
		cmd += fmt.Sprintf("%dm ", minutes)
	}

	if seconds > 0 {
		cmd += fmt.Sprintf("%ds", seconds)
	}

	cmd += ledboard.PauseSecond4 + "0120"

	return cmd
}
