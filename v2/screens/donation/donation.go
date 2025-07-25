package donation

import (
	"github.com/b4ckspace/ledboard-v2/ledboard"
)

// GenerateDonationCommand generates the command string for the donation screen.
func GenerateDonationCommand() string {
	var cmd string

	cmd += ledboard.FontNormal16x9
	cmd += ledboard.ControlPatternIn + ledboard.PatternScrollUp
	cmd += ledboard.ControlPatternOut + ledboard.PatternScrollUp

	cmd += ledboard.ControlFlash + ledboard.FlashOn

	cmd += ledboard.ControlFontColor + ledboard.FontColorYGRCharacter
	cmd += "\\o/ Spende! \\o/"

	cmd += ledboard.ControlFlash + ledboard.FlashOff

	cmd += ledboard.PauseSecond2 + "04"

	return cmd
}

