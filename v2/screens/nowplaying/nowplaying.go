package nowplaying

import (
	"github.com/b4ckspace/ledboard-v2/ledboard"
	"github.com/b4ckspace/ledboard-v2/utils"
)

// GenerateNowPlayingCommand generates the command string for the now playing screen.
func GenerateNowPlayingCommand(message string) string {
	var cmd string

	cmd += ledboard.ControlPatternIn + ledboard.PatternRadarScan

	cmd += ledboard.FontNormal7x6

	cmd += ledboard.ControlFontColor + ledboard.FontColorYellow
	cmd += "NOW PLAYING"

	cmd += ledboard.PauseSecond2 + "05"
	cmd += ledboard.ControlFrame

	cmd += utils.SanitizeUmlauts(message)
	cmd += ledboard.PauseSecond2 + "45"

	return cmd
}
