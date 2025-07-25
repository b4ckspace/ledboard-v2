package pizzatimer

import (
	"github.com/b4ckspace/ledboard-v2/ledboard"
)

// GeneratePizzaTimerCommand generates the command string for the pizza timer screen.
func GeneratePizzaTimerCommand() string {
	var cmd string

	cmd += ledboard.FontNormal16x9
	cmd += ledboard.ControlPatternIn + ledboard.PatternScrollUp
	cmd += ledboard.ControlPatternOut + ledboard.PatternScrollUp

	cmd += ledboard.ControlFlash + ledboard.FlashOn
	cmd += ledboard.ControlFontColor + ledboard.FontColorYGRCharacter
	cmd += "PIZZA IS READY!"

	cmd += ledboard.ControlFlash + ledboard.FlashOff

	cmd += ledboard.PauseSecond2 + "10"

	return cmd
}
