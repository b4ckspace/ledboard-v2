package publicserviceannouncement

import (
	"github.com/b4ckspace/ledboard-v2/ledboard"
	"github.com/b4ckspace/ledboard-v2/utils"
)

// GeneratePublicServiceAnnouncementCommand generates the command string for the public service announcement screen.
func GeneratePublicServiceAnnouncementCommand(message string) string {
	var cmd string

	cmd += ledboard.ControlPatternIn + ledboard.PatternRadarScan
	cmd += ledboard.ControlFlash + ledboard.FlashOn

	cmd += ledboard.FontNormal7x6

	cmd += ledboard.ControlFontColor + ledboard.FontColorYellow
	cmd += "PUBLIC "

	cmd += ledboard.ControlFontColor + ledboard.FontColorRed
	cmd += "SERVICE "

	cmd += ledboard.ControlFontColor + ledboard.FontColorGreen
	cmd += "ANNOUNCEMENT"

	cmd += ledboard.ControlFlash + ledboard.FlashOff

	cmd += ledboard.PauseSecond2 + "05"
	cmd += ledboard.ControlFrame

	cmd += utils.SanitizeUmlauts(message)
	cmd += ledboard.PauseSecond2 + "45"

	return cmd
}
