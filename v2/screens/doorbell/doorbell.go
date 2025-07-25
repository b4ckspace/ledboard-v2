package doorbell

import (
	"github.com/b4ckspace/ledboard-v2/ledboard"
)

// GenerateDoorBellCommand generates the command string for the doorbell screen.
func GenerateDoorBellCommand() string {
	var cmd string

	cmd += ledboard.ControlPatternIn + ledboard.PatternScrollUp
	cmd += ledboard.ControlPatternOut + ledboard.PatternScrollUp

	cmd += ledboard.FontNormal16x9
	cmd += ledboard.ControlFlash + ledboard.FlashOn
	cmd += ledboard.ControlFontColor + ledboard.FontColorRed
	cmd += "! DOORBELL !"
	cmd += ledboard.ControlFlash + ledboard.FlashOff
	cmd += ledboard.PauseSecond2 + "10"

	return cmd
}
