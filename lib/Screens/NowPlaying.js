'use strict';

const Commands = require('../LedBoard/Commands');
const commonUtils = require('../Utils/Common.js');

module.exports = (message) => {
    let cmd = '';

    cmd += Commands.Control.PATTERN_IN + Commands.Pattern.RADAR_SCAN;

    cmd += Commands.Font.NORMAL_7x6;

    cmd += Commands.Control.FONT_COLOR + Commands.FontColor.YELLOW;
    cmd += 'NOW PLAYING';

    cmd += Commands.Pause.SECOND_2 + '05';
    cmd += Commands.Control.FRAME;

    cmd += commonUtils.sanitizeUmlauts(message);
    cmd += Commands.Pause.SECOND_2 + '45';

    return cmd;
};
