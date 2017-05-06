package hearthstone

import (
	"time"
)

type Game struct {
	PlayerDraws   []string
	OpponentPlays []string
	LocalPlayer   *Player
	Opponent      *Player
	Win           bool
	Signature     string
	Date          string

	startTime     string
	accountIds    [2]string
	localDraws    map[Card]bool
	opponentPlays map[Card]bool
}

func NewGame(startTime string) *Game {
	game := &Game{startTime: startTime}
	game.localDraws = make(map[Card]bool)
	game.opponentPlays = make(map[Card]bool)
	game.Date = time.Now().Format("2006-01-02")
	return game
}

func (game *Game) InProgress() bool {
	return game.Signature != ""
}

func (game *Game) AddAccount(account string) {
	if game.accountIds[0] == "" {
		game.accountIds[0] = account
	} else if game.accountIds[1] == "" {
		game.accountIds[1] = account
		game.Signature = game.startTime + "/" + game.accountIds[0] + "/" + game.accountIds[1]
	} else {
		panic("Tried to define more than two accounts for a game.")
	}
}

func (game *Game) PlayerDraw(cardID string, deckID string, name string) {
	game.localDraws[Card{CardID: cardID, DeckID: deckID, Name: name}] = true
}

func (game *Game) OpponentPlay(cardID string, deckID string, name string) {
	game.opponentPlays[Card{CardID: cardID, DeckID: deckID, Name: name}] = true
}

func (game *Game) GameOver() {
	game.PlayerDraws = make([]string, len(game.localDraws))
	i := 0
	for card := range game.localDraws {
		game.PlayerDraws[i] = card.CardID
		i++
	}
	game.OpponentPlays = make([]string, len(game.opponentPlays))
	i = 0
	for card := range game.opponentPlays {
		game.OpponentPlays[i] = card.CardID
		i++
	}
}

type Card struct {
	CardID string
	DeckID string
	Name   string
}

type Player struct {
	Name  string
	ID    string
	Class string
}
