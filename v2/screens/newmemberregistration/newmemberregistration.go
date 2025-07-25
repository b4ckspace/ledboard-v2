package newmemberregistration

import (
	"github.com/b4ckspace/ledboard-v2/ledboard"
)

// GenerateNewMemberRegistrationCommand generates the command string for the new member registration screen.
func GenerateNewMemberRegistrationCommand(nickname string) string {
	var cmd string

	cmd += ledboard.ControlPatternIn + ledboard.PatternRadarScan

	colors := []string{ledboard.FontColorGreen, ledboard.FontColorRed, ledboard.FontColorYGRHorizontal}
	for _, color := range colors {
		cmd += ledboard.FontNormal7x6
		cmd += ledboard.ControlFontColor + color
		cmd += "Herzlich Willkommen im backspace!"
		cmd += ledboard.PauseSecond2 + "01"
		cmd += ledboard.ControlFrame
	}

	cmd += ledboard.FontNormal16x9
	cmd += ledboard.ControlFontColor + ledboard.FontColorYellow
	cmd += ledboard.ControlPatternIn + ledboard.PatternMoveUp
	cmd += ledboard.ControlPatternOut + ledboard.PatternMoveLeft
	cmd += nickname
	cmd += ledboard.PauseSecond2 + "30"

	return cmd
}
