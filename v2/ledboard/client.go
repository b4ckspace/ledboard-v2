package ledboard

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/b4ckspace/ledboard-v2/utils"
)

// Client represents the LED board client
type Client struct {
	host string
	port int
	conn *net.UDPConn
}

// NewClient creates a new LED board client instance
func NewClient(host string, port int) *Client {
	return &Client{host: host, port: port}
}

// Init initializes the LED board client
func (c *Client) Init() error {
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", c.host, c.port))
	if err != nil {
		return fmt.Errorf("failed to resolve UDP address: %w", err)
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return fmt.Errorf("failed to dial UDP: %w", err)
	}
	c.conn = conn
	log.Println("LED Board Client initialized")
	return nil
}

// Send sends a datagram to the LED board.
func (c *Client) Send(datagram string) {
	if c.conn == nil {
		log.Println("UDP connection not initialized, cannot send.")
		return
	}

	message := []byte(datagram)

	// Debugging: Print the raw datagram being sent
	log.Printf("Sending to LED board (raw): %q\n", datagram)
	log.Printf("Sending to LED board (hex): %x\n", message)

	_, err := c.conn.Write(message)
	if err != nil {
		log.Printf("Failed sending UDP message: %v", err)
	}
}

// SetDate sets the date on the LED board.
func (c *Client) SetDate(date time.Time) {
	var cmd string
	cmd += "\x01Z00\x02E"
	cmd += "B"

	cmd += utils.Byte2Hex(byte(date.Year()%100)) + utils.Byte2Hex(byte(date.Year()/100))
	cmd += utils.Byte2Hex(byte(date.Month()))
	cmd += utils.Byte2Hex(byte(date.Day()))
	cmd += utils.Byte2Hex(byte(date.Hour()))
	cmd += utils.Byte2Hex(byte(date.Minute()))
	cmd += utils.Byte2Hex(0)
	cmd += utils.Byte2Hex(0)

	cmd += ControlEnd

	c.Send(cmd)
}

// SendScreen sends a single screen command to the LED board.
func (c *Client) SendScreen(screen string) {
	datagram := c.buildDatagram(screen)
	c.Send(datagram)
}

// SendScreens sends multiple screen commands to the LED board, joined by a frame.
func (c *Client) SendScreens(screens []string) {
	combinedScreen := ""
	for i, screen := range screens {
		combinedScreen += screen
		if i < len(screens)-1 {
			combinedScreen += ControlFrame
		}
	}
	c.SendScreen(combinedScreen)
}

func (c *Client) buildDatagram(screen string) string {
	var cmd string
	cmd += ControlHead

	cmd = "\x01Z00\x02A"
	cmd += "\x0fETAA" // store to RAM

	cmd += screen
	cmd += ControlEnd

	return cmd
}
