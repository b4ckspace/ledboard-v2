package idle

import (
	"fmt"
	"github.com/b4ckspace/ledboard-v2/ledboard"
)

// GenerateIdleCommand generates the command string for the idle screen.
func GenerateIdleCommand(memberCount int) string {
	var cmd string

	cmd += ledboard.FontNormal7x6
	cmd += ledboard.ControlPatternIn + ledboard.PatternScrollUp
	cmd += ledboard.ControlPatternOut + ledboard.PatternScrollUp

	cmd += ledboard.ControlFontColor + ledboard.FontColorGreen

	cmd += ledboard.ControlSpecial + ledboard.SpecialYYYY + "-"
	cmd += ledboard.ControlSpecial + ledboard.SpecialMM + "-"
	cmd += ledboard.ControlSpecial + ledboard.SpecialDD + " "

	cmd += ledboard.ControlFontColor + ledboard.FontColorRed

	cmd += ledboard.ControlSpecial + ledboard.SpecialHH + ":"
	cmd += ledboard.ControlSpecial + ledboard.SpecialMIN + ":"
	cmd += ledboard.ControlSpecial + ledboard.SpecialSEC

	cmd += ledboard.ControlLineFeed

	cmd += ledboard.ControlFontColor + ledboard.FontColorYellow
	cmd += fmt.Sprintf("humans present: %d", memberCount)
	cmd += ledboard.PauseSecond4 + "9999"

	return cmd
}
