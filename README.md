Tracks Hearthstone games you play and rates your cards relative to your opponents.

## Building

The current version assumes you're already generating power.log and hardcodes the path to that file on Windows. You might have to adjust that depending on what you're running. For help see [this page](https://github.com/jleclanche/fireplace/wiki/How-to-enable-logging)

To build, install Go, then run `go build` in the main directory. That will generate the `HSTeamPlay` executable that you can run from a command line.

## Usage

With no arguments, HSTeamPlay will read your power.log line by line, look for games you're playing, rate each player's cards, and save them to a DB. Leave it running while you play and it will figure things out.

To **show the ratings** for a set of cards, run `HSTeamPlay --rate=<cards.txt>`. Card names must be in the format `friendly/DRUID/Living Mana` or `opponent/WARRIOR/Patches The Pirate`. HSTeamPlay will show ratings for the cards and exit.

To **reset the ratings** for a set of cards, run `HSTeamPlay --reset=<cards.txt>`. Card names must be in the format `friendly/DRUID/Living Mana` or `opponent/WARRIOR/Patches The Pirate`. HSTeamPlay will reset the ratings for all cards in the file with no warning and then exit. Hope you were sure.

Uses TrueSkill implemented by @ChrisHines. https://github.com/ChrisHines/GoSkills
