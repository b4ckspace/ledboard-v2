package screens

import (
	"log"

	"github.com/b4ckspace/ledboard-v2/screens/alarm"
	"github.com/b4ckspace/ledboard-v2/screens/donation"
	"github.com/b4ckspace/ledboard-v2/screens/doorbell"
	"github.com/b4ckspace/ledboard-v2/screens/idle"
	"github.com/b4ckspace/ledboard-v2/screens/laserfinished"
	"github.com/b4ckspace/ledboard-v2/screens/laseroperation"
	"github.com/b4ckspace/ledboard-v2/screens/newmemberregistration"
	"github.com/b4ckspace/ledboard-v2/screens/nowplaying"
	"github.com/b4ckspace/ledboard-v2/screens/pizzatimer"
	"github.com/b4ckspace/ledboard-v2/screens/publicserviceannouncement"
)

// Screens represents the main screens manager
type Screens struct {
	// Add fields for managing different screens
}

// NewScreens creates a new Screens instance
func NewScreens() *Screens {
	return &Screens{}
}

// Init initializes the screens manager
func (s *Screens) Init() {
	log.Println("Screens initialized")
	// TODO: Implement screen initialization logic based on original Screens.js
}

// Alarm generates the command for the Alarm screen.
func (s *Screens) Alarm(message string) string {
	return alarm.GenerateAlarmCommand(message)
}

// Donation generates the command for the Donation screen.
func (s *Screens) Donation() string {
	return donation.GenerateDonationCommand()
}

// DoorBell generates the command for the DoorBell screen.
func (s *Screens) DoorBell() string {
	return doorbell.GenerateDoorBellCommand()
}

// Idle generates the command for the Idle screen.
func (s *Screens) Idle(memberCount int) string {
	return idle.GenerateIdleCommand(memberCount)
}

// LaserFinished generates the command for the LaserFinished screen.
func (s *Screens) LaserFinished(duration int) string {
	return laserfinished.GenerateLaserFinishedCommand(duration)
}

// LaserOperation generates the command for the LaserOperation screen.
func (s *Screens) LaserOperation() string {
	return laseroperation.GenerateLaserOperationCommand()
}

// NewMemberRegistration generates the command for the NewMemberRegistration screen.
func (s *Screens) NewMemberRegistration(nickname string) string {
	return newmemberregistration.GenerateNewMemberRegistrationCommand(nickname)
}

// NowPlaying generates the command for the NowPlaying screen.
func (s *Screens) NowPlaying(message string) string {
	return nowplaying.GenerateNowPlayingCommand(message)
}

// PizzaTimer generates the command for the PizzaTimer screen.
func (s *Screens) PizzaTimer() string {
	return pizzatimer.GeneratePizzaTimerCommand()
}

// PublicServiceAnnouncement generates the command for the PublicServiceAnnouncement screen.
func (s *Screens) PublicServiceAnnouncement(message string) string {
	return publicserviceannouncement.GeneratePublicServiceAnnouncementCommand(message)
}
