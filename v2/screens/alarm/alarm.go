package alarm

import (
	"github.com/b4ckspace/ledboard-v2/ledboard"
	"github.com/b4ckspace/ledboard-v2/utils"
)

// GenerateAlarmCommand generates the command string for the alarm screen.
func GenerateAlarmCommand(message string) string {
	var cmd string

	cmd += ledboard.ControlPatternIn + ledboard.PatternRadarScan

	cmd += ledboard.FontNormal16x9
	cmd += ledboard.ControlFlash + ledboard.FlashOn
	cmd += ledboard.ControlFontColor + ledboard.FontColorRed
	cmd += "!  ALARM  !"
	cmd += ledboard.ControlFlash + ledboard.FlashOff
	cmd += ledboard.PauseSecond2 + "04"

	cmd += ledboard.ControlFrame

	cmd += ledboard.FontNormal7x6
	cmd += ledboard.ControlFontColor + ledboard.FontColorGreen
	cmd += ledboard.ControlPatternIn + ledboard.PatternMoveUp
	cmd += ledboard.ControlPatternOut + ledboard.PatternMoveLeft
	cmd += utils.SanitizeUmlauts(message)
	cmd += ledboard.PauseSecond2 + "30"

	return cmd
}
