package hearthstone

import (
	"regexp"
	"sync"
)

var ignoredEvents = regexp.MustCompile("PowerTaskList|PowerProcessor")

var playerIdentifiers = regexp.MustCompile("id=([12]) Player=([\\S]+).*ChoiceType=MULLIGAN")
var gameStart = regexp.MustCompile("(\\d{2}:\\d{2}:\\d{2}.\\d{7}) .* CREATE_GAME")
var gameAccounts = regexp.MustCompile("GameAccountId=\\[(hi=[\\S]+) (lo=[\\S]+)\\]")
var gameEnd = regexp.MustCompile("TAG_CHANGE Entity=GameEntity tag=STATE value=COMPLETE")
var findWinner = regexp.MustCompile("TAG_CHANGE Entity=([\\S]+) tag=PLAYSTATE value=WON")
var findClass = regexp.MustCompile("TAG_CHANGE Entity=\\[name=.* zone=PLAY zonePos=0 cardId=(HERO_[0-9][1-9]).* player=([12])]")

var extractCardName = regexp.MustCompile("name=([a-zA-Z ]+) id=")
var cardRevealed = regexp.MustCompile("SHOW_ENTITY - Updating Entity=\\[name=UNKNOWN ENTITY \\[cardType=INVALID\\] id=([0-9]+) zone=([\\S]*).*player=([12])] CardID=([\\S]+)")
var discoverCardAddedToHand = regexp.MustCompile("TAG_CHANGE Entity=\\[.*id=([0-9]+) zone=SETASIDE .* cardId=([\\S]+) player=([12])] tag=ZONE value=HAND")
var discoverOptions = regexp.MustCompile("GameState.DebugPrintEntityChoices\\(\\).*Entities\\[.\\]=\\[.*zone=SETASIDE.*cardId=([\\S]+) player=([12])\\]")
var keptInMulligan = regexp.MustCompile("GameState.SendChoices\\(\\).*m_chosenEntities\\[.\\]=\\[.* id=([0-9]+) zone=HAND .* cardId=([\\S]+) player=([12])\\]")

func ProcessEvents(events chan string, waitGroup sync.WaitGroup) {
	defer waitGroup.Done()
	var game *Game
	var playerMap map[string]*Player
	var cardRater *CardRater = MakeCardRater()
	var gameDB *GameDB = GetGameDB()
	var cardDB map[string]CardInfo = GetCardInfo()
	for event := range events {
		if ignoredEvents.MatchString(event) {
			continue
		}
		if game == nil || !game.InProgress() {
			gameInfo := gameStart.FindStringSubmatch(event)
			if len(gameInfo) > 1 {
				println("---STARTING A NEW GAME---")
				game = NewGame(gameInfo[1])
				playerMap = make(map[string]*Player)
				playerMap["1"] = &Player{ID: "1"}
				playerMap["2"] = &Player{ID: "2"}
			}
			accountInfo := gameAccounts.FindStringSubmatch(event)
			if len(accountInfo) > 1 {
				game.AddAccount(accountInfo[1] + ":" + accountInfo[2])
			}
		} else if gameDB.HaveLoggedGame(game) {
			println("Looks like I've already logged this game; ignoring it.")
			game = nil
		} else {
			mulliganKeeps := keptInMulligan.FindStringSubmatch(event)
			if game != nil && len(mulliganKeeps) > 1 {
				println("Mulligan choice: " + cardDB[mulliganKeeps[2]].Name)
				game.PlayerDraw(mulliganKeeps[2], mulliganKeeps[1], cardDB[mulliganKeeps[2]].Name)
			}
			playerInfo := playerIdentifiers.FindStringSubmatch(event)
			if len(playerInfo) > 1 {
				playerMap[playerInfo[1]].Name = playerInfo[2]
			}
			if (playerMap["1"] != nil && playerMap["1"].Class == "") || (playerMap["2"] != nil && playerMap["2"].Class == "") {
				classInfo := findClass.FindStringSubmatch(event)
				if len(classInfo) > 1 {
					playerMap[classInfo[2]].Class = cardDB[classInfo[1]].PlayerClass
				}
			}
			drawInfo := cardRevealed.FindStringSubmatch(event)
			if len(drawInfo) > 1 {
				if game.LocalPlayer == nil {
					if drawInfo[2] == "DECK" {
						game.LocalPlayer = playerMap[drawInfo[3]]
						game.Opponent = playerMap[oppositePlayerId(drawInfo[3])]
					} else {
						game.LocalPlayer = playerMap[oppositePlayerId(drawInfo[3])]
						game.Opponent = playerMap[drawInfo[3]]
					}
				}
				if drawInfo[2] == "HAND" || (drawInfo[2] == "DECK" && playerMap[drawInfo[3]] == game.Opponent) {
					println("Opponent played " + cardDB[drawInfo[4]].Name + " " + cardRater.GetCardQualityAsString(raterKey(cardDB[drawInfo[4]].Name, "opponent", game.Opponent.Class)))
					game.OpponentPlay(drawInfo[4], drawInfo[1], cardDB[drawInfo[4]].Name)
				} else if drawInfo[2] == "DECK" {
					println("You drew " + cardDB[drawInfo[4]].Name + " " + cardRater.GetCardQualityAsString(raterKey(cardDB[drawInfo[4]].Name, "friendly", game.LocalPlayer.Class)))
					game.PlayerDraw(drawInfo[4], drawInfo[1], cardDB[drawInfo[4]].Name)
				}
			}
			discoverInfo := discoverCardAddedToHand.FindStringSubmatch(event)
			if len(discoverInfo) > 2 {
				println("You discovered " + cardDB[discoverInfo[2]].Name + " " + cardRater.GetCardQualityAsString(raterKey(cardDB[discoverInfo[2]].Name, "friendly", game.LocalPlayer.Class)))
				game.PlayerDraw(discoverInfo[2], discoverInfo[1], cardDB[discoverInfo[2]].Name)
			}
			discoverOption := discoverOptions.FindStringSubmatch(event)
			if len(discoverOption) > 2 {
				println("Choice: " + cardDB[discoverOption[1]].Name + " " + cardRater.GetCardQualityAsString(raterKey(cardDB[discoverOption[2]].Name, "friendly", game.LocalPlayer.Class)))
			}
			winnerInfo := findWinner.FindStringSubmatch(event)
			if game.LocalPlayer != nil && len(winnerInfo) > 1 {
				game.Win = winnerInfo[1] == game.LocalPlayer.Name
			}
			if gameEnd.MatchString(event) {
				if game.LocalPlayer == nil || game.Opponent == nil {
					println("I wasn't able to figure out who the friendly player was. Not logging this game.")
				} else {
					game.GameOver()
					gameDB.AddGame(game)
					println("Game over")
					if game.Win {
						println("Congratulations!")
					} else {
						println("Better luck next time!")
					}
					if len(game.PlayerDraws) == 0 && len(game.OpponentPlays) == 0 {
						println("Someone didn't play any cards, so I'm not rating this game.")
					} else {
						playerCards := idsToNames(cardDB, game.PlayerDraws, "friendly", game.LocalPlayer.Class)
						opponentCards := idsToNames(cardDB, game.OpponentPlays, "opponent", game.Opponent.Class)
						if game.Win {
							cardRater.UpdateRatings(playerCards, opponentCards)
						} else {
							cardRater.UpdateRatings(opponentCards, playerCards)
						}
						println("Opponent's cards now rated:")
						for _, ratingString := range cardRater.RateCardGroup(opponentCards) {
							println(ratingString)
						}
						println("Player's cards now rated:")
						for _, ratingString := range cardRater.RateCardGroup(playerCards) {
							println(ratingString)
						}
					}
				}
				game = nil
			}
		}
	}
}

func idsToNames(cardDB map[string]CardInfo, cardIDs []string, playerAllegiance string, deckClass string) []string {
	names := make([]string, len(cardIDs))
	for i, cardID := range cardIDs {
		names[i] = raterKey(cardDB[cardID].Name, playerAllegiance, deckClass)
	}
	return names
}

func raterKey(cardName string, playerAllegiance string, deckClass string) string {
	return playerAllegiance + "/" + deckClass + "/" + cardName
}

func oppositePlayerId(id string) string {
	if id == "1" {
		return "2"
	}
	return "1"
}
