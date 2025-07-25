package laseroperation

import (
	"github.com/b4ckspace/ledboard-v2/ledboard"
)

// GenerateLaserOperationCommand generates the command string for the laser operation screen.
func GenerateLaserOperationCommand() string {
	var cmd string

	cmd += ledboard.ControlPatternIn + ledboard.PatternRadarScan

	cmd += ledboard.FontNormal15x9
	cmd += ledboard.ControlFontColor + ledboard.FontColorRed

	cmd += ledboard.ControlSpecial + ledboard.SpecialHH + "h "
	cmd += ledboard.ControlSpecial + ledboard.SpecialMIN + "m "
	cmd += ledboard.ControlSpecial + ledboard.SpecialSEC + "s "

	cmd += ledboard.PauseSecond4 + "9999"

	return cmd
}
