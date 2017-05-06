package hearthstone

import (
	"encoding/json"
	"io/ioutil"
)

type CardInfo struct {
	ID          string
	Name        string
	PlayerClass string
}

var cardMap map[string]CardInfo

func GetCardInfo() map[string]CardInfo {
	if cardMap == nil {
		cardMap = make(map[string]CardInfo)
		file, err := ioutil.ReadFile("./cardDB.json")
		check(err)

		var rawCardList []map[string]*json.RawMessage
		err = json.Unmarshal(file, &rawCardList)
		check(err)

		for _, card := range rawCardList {
			cardID := unmarshalPossiblyEmptyValue(card, "id")
			cardName := unmarshalPossiblyEmptyValue(card, "name")
			cardClass := unmarshalPossiblyEmptyValue(card, "playerClass")
			cardMap[cardID] = CardInfo{ID: cardID, Name: cardName, PlayerClass: cardClass}
		}
	}
	return cardMap
}

func unmarshalPossiblyEmptyValue(rawCard map[string]*json.RawMessage, key string) string {
	value := "Unknown"
	if _, contains := rawCard[key]; contains {
		err := json.Unmarshal(*rawCard[key], &value)
		check(err)
	}
	return value
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
