package hearthstone

import (
	"encoding/json"
	"os"
	"path"
)

type GameDB struct {
	Games map[string]Game
}

func GetGameDB() *GameDB {
	gameDB, err := loadGameDB()
	if err != nil {
		println("WARNING: Couldn't load game history. Starting a new log. The error was: " + err.Error())
		gameDB = new(GameDB)
		gameDB.Games = make(map[string]Game)
	}
	return gameDB
}

func loadGameDB() (*GameDB, error) {
	gameDB := new(GameDB)
	file, err := os.Open(gameDBPath())
	if err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		err = decoder.Decode(gameDB)
		check(err)
	}
	return gameDB, err
}

func (gameDB *GameDB) saveGameDB() error {
	file, err := os.Create(gameDBPath())
	if err == nil {
		defer file.Close()
		encoder := json.NewEncoder(file)
		err = encoder.Encode(gameDB)
	}
	return err
}

func gameDBPath() string {
	return path.Join(".", "gameDB.json")
}

func (gameDB *GameDB) AddGame(game *Game) {
	if !gameDB.HaveLoggedGame(game) && game.Signature != "" {
		gameDB.Games[game.Signature] = *game
		err := gameDB.saveGameDB()
		check(err)
	}
}

func (gameDB *GameDB) HaveLoggedGame(game *Game) bool {
	_, ok := gameDB.Games[game.Signature]
	return ok
}
