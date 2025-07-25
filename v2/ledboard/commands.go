package ledboard

// Control commands
const (
	ControlHead             = "\x01\x5A\x30\x30\x02\x41\x58"
	ControlHalfSpace        = "\x82"
	ControlEnd              = "\x04"
	ControlFlash            = "\x07"
	ControlPatternIn        = "\x06\x0aI"
	ControlPatternOut       = "\x06\x0aO"
	ControlSpecial          = "\x0B"
	ControlFrame            = "\x0C"
	ControlSpeed            = "\x0F"
	ControlLineFeed         = "\x0D"
	ControlFontColor        = "\x1C"
	ControlBackgroundColor  = "\x1D"
	ControlAlignHorizontal  = "\x1E"
	ControlAlignVertical    = "\x1F"
)

// Flash commands
const (
	FlashOff = "\x30"
	FlashOn  = "\x31"
)

// Special commands
const (
	SpecialMMDDYYSLA = "\x20"
	SpecialDDMMYYSLA = "\x21"
	SpecialMMDDYYDSH = "\x22"
	SpecialDDMMYYDSH = "\x23"
	SpecialMMDDYYYYDOT = "\x24"
	SpecialYY        = "\x25"
	SpecialYYYY      = "\x26"
	SpecialMM        = "\x27"
	SpecialMMM       = "\x28"
	SpecialDD        = "\x29"
	SpecialDDOfWeek  = "\x2A"
	SpecialDDDOfWeek = "\x2B"
	SpecialHH        = "\x2C"
	SpecialMIN       = "\x2D"
	SpecialSEC       = "\x2E"
	SpecialHHMin24   = "\x2F"
	SpecialHHMin12   = "\x30"
)

// Pattern commands
const (
	PatternRandom                  = "\x2F"
	PatternJumpOut                 = "\x30"
	PatternMoveLeft                = "\x31"
	PatternMoveRight               = "\x32"
	PatternScrollLeft              = "\x33"
	PatternScrollRight             = "\x34"
	PatternMoveUp                  = "\x35"
	PatternMoveDown                = "\x36"
	PatternScrollLR                = "\x37"
	PatternScrollUp                = "\x38"
	PatternScrollDown              = "\x39"
	PatternFoldLR                  = "\x3A"
	PatternFoldUD                  = "\x3B"
	PatternScrollUD                = "\x3C"
	PatternShuttleLR               = "\x3D"
	PatternShuttleUD               = "\x3E"
	PatternPeelOffL                = "\x3F"
	PatternPeelOffR                = "\x40"
	PatternShutterUD               = "\x41"
	PatternShutterLR               = "\x42"
	PatternRaindrops               = "\x43"
	PatternRandomMosaic            = "\x44"
	PatternTwinklingStar           = "\x45"
	PatternHipHop                  = "\x46"
	PatternRadarScan               = "\x47"
	PatternFanOut                  = "\x48"
	PatternFanIn                   = "\x49"
	PatternSpiralR                 = "\x4A"
	PatternSpiralL                 = "\x4B"
	PatternToFourCorners           = "\x4C"
	PatternFromFourCorners         = "\x4D"
	PatternToFourSides             = "\x4E"
	PatternFromFourSides           = "\x4F"
	PatternScrollOutFromFourBlocks = "\x50"
)

// Pause commands
const (
	PauseSecond2      = "\x0E\x30"
	PauseSecond4      = "\x0E\x32"
	PauseMillisecond2 = "\x0E\x31"
	PauseMillisecond4 = "\x0E\x33"
)

// Align commands
const (
	AlignHorizontalCenter    = "\x30"
	AlignHorizontalLeft      = "\x31"
	AlignHorizontalRight     = "\x32"
	AlignHorizontalReserved  = "\x33"
	AlignVerticalCenter      = "\x30"
	AlignVerticalTop         = "\x31"
	AlignVerticalBottom      = "\x32"
	AlignVerticalReserved    = "\x33"
)

// BackgroundColor commands
const (
	BackgroundColorBlack  = "\x30"
	BackgroundColorRed    = "\x31"
	BackgroundColorGreen  = "\x32"
	BackgroundColorYellow = "\x33"
)

// Font commands
const (
	FontNormal5x5   = "\x1A\x30"
	FontNormal7x6   = "\x1A\x31"
	FontNormal14x8  = "\x1A\x32"
	FontNormal11x9  = "\x1A\x3A"
	FontNormal15x9  = "\x1A\x33"
	FontNormal16x9  = "\x1A\x34"
	FontNormal24x16 = "\x1A\x36"
	FontNormal22x18 = "\x1A\x3C"
	FontNormal30x18 = "\x1A\x3D"
	FontNormal32x18 = "\x1A\x38"
	FontNormal40x21 = "\x1A\x3E"
	FontBold5x7     = "\x1A\x4D"
	FontBold14x10   = "\x1A\x4E"
	FontBold15x10   = "\x1A\x4F"
	FontBold16x12   = "\x1A\x50"
	FontCustom3     = "\x1A\x63"
	FontCustom4     = "\x1A\x64"
	FontCustom5     = "\x1A\x65"
	FontCustom6     = "\x1A\x66"
	FontCustom7     = "\x1A\x67"
	FontCustom8     = "\x1A\x68"
	FontCustom9     = "\x1A\x69"
)

// FontColor commands
const (
	FontColorBlack        = "\x30"
	FontColorRed          = "\x31"
	FontColorGreen        = "\x32"
	FontColorYellow       = "\x33"
	FontColorYGRCharacter = "\x34"
	FontColorYGRHorizontal = "\x35"
	FontColorYGRWave      = "\x36"
	FontColorYGRDiagonal  = "\x37"
)

// Speed commands
const (
	SpeedVeryFast    = "\x30"
	SpeedFast        = "\x31"
	SpeedMediumFast  = "\x32"
	SpeedMedium      = "\x33"
	SpeedMediumSlow  = "\x34"
	SpeedSlow        = "\x35"
	SpeedVerySlow    = "\x36"
)