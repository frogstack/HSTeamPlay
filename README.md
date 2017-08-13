Tracks Hearthstone games you play and rates your cards relative to your opponents.

## Building

The current version assumes you're already generating power.log, which you are if you already have a tracker installed. If you've never installed a tracker, the easiest way to get this working is to install one. HSTeamPlay should be able to find your power.log file whether you're on Windows or Mac, but I'm on Windows and haven't tested it on Mac. If you find a problem, feel free to send a pull request.

To build, install Go, then:

* Clone this repo to `$GOROOT/src/github.com/frogstack/HSTeamPlay/`
* Run `go get github.com/ChrisHines/GoSkills/skills`. If this does not work, clone it to `$GOROOT/src/ChrisHines/GoSkills/skills/`
* Run `go build` in `$GOROOT/src/github.com/frogstack/HSTeamPlay/`

Go will generate the `HSTeamPlay` executable that you can run from a command line.

Note: HSTeamPlay tries to locate your Hearthstone install path automatically, if that fails, you can set the path manually in main.go like this: `programfiles := "C:/Path/to/your/Hearthstone/"`.

## Usage

With no arguments, HSTeamPlay will read your power.log line by line, look for games you're playing, rate each player's cards, and save them to a DB. Leave it running while you play and it will figure things out.

To **show the ratings** for a set of cards, run `HSTeamPlay --rate=<cards.txt>`. Card names must be in the format `friendly/DRUID/Living Mana` or `opponent/WARRIOR/Patches The Pirate`. HSTeamPlay will show ratings for the cards and exit. You can find an example cards.txt for a Pirate Warrior deck in the source.

To **reset the ratings** for a set of cards, run `HSTeamPlay --reset=<cards.txt>`. Card names must be in the format `friendly/DRUID/Living Mana` or `opponent/WARRIOR/Patches The Pirate`. HSTeamPlay will reset the ratings for all cards in the file with no warning and then exit. Hope you were sure.
