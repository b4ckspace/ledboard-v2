package screens

import (
	"fmt"

	"github.com/b4ckspace/ledboard-v2/ledboard"
	"github.com/b4ckspace/ledboard-v2/utils"
)

// Screens represents the main screens manager
type Screens struct {
	// Add fields for managing different screens
}

// NewScreens creates a new Screens instance
func NewScreens() *Screens {
	return &Screens{}
}

// Alarm generates the command for the Alarm screen.
func (s *Screens) Alarm(message string) string {
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

// Donation generates the command string for the donation screen.
func (s *Screens) Donation() string {
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

// DoorBell generates the command string for the doorbell screen.
func (s *Screens) DoorBell() string {
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

// Idle generates the command string for the idle screen.
func (s *Screens) Idle(memberCount int) string {
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

// LaserFinished generates the command string for the laser finished screen.
func (s *Screens) LaserFinished(duration int) string {
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

	cmd += ledboard.ControlPatternIn + ledboard.PatternPeelOffR

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

// LaserOperation generates the command string for the laser operation screen.
func (s *Screens) LaserOperation() string {
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

// NewMemberRegistration generates the command string for the new member registration screen.
func (s *Screens) NewMemberRegistration(nickname string) string {
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

// NowPlaying generates the command string for the now playing screen.
func (s *Screens) NowPlaying(message string) string {
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

// PizzaTimer generates the command string for the pizza timer screen.
func (s *Screens) PizzaTimer() string {
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

// PublicServiceAnnouncement generates the command string for the public service announcement screen.
func (s *Screens) PublicServiceAnnouncement(message string) string {
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
