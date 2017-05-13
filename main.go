package main

import (
	"github.com/frogstack/HSTeamPlay/hearthstone"
	"github.com/frogstack/HSTeamPlay/tail"
	"bufio"
	"flag"
	"os"
	"path"
	"runtime"
	"sync"
)

func main() {
	//reset := flag.String("reset", "", "Reset the ratings for cards that look like the provided value and exit")
	rateFile := flag.String("rate", "", "Show the ratings for the cards in the file provided and exit")
	resetFile := flag.String("reset", "", "Reset ratings for the cards in the file provided and exit")
	flag.Parse()

	if *rateFile != "" {
		cards := getCardsFromFile(rateFile)
		rater := hearthstone.MakeCardRater()
		ratings := rater.RateCardGroup(cards)
		for _, rating := range ratings {
			println(rating)
		}
	} else if *resetFile != "" {
		cards := getCardsFromFile(resetFile)
		rater := hearthstone.MakeCardRater()
		rater.ResetRatings(cards)
		println("Hope you were sure, because I reset the ratings for all cards in " + *resetFile)
	} else {
		var powerlog string
		if runtime.GOOS == "windows" {
			programfiles := os.Getenv("programfiles(x86)")
			powerlog = path.Join(programfiles, "Hearthstone", "Logs", "Power.log")
		} else if runtime.GOOS == "darwin" {
			powerlog = path.Join("Applications", "Hearthstone", "Logs", "Power.log")
		} else {
			panic("I don't recognize your OS. Sorry!")
		}
		tail, err := tail.TailFile(powerlog)
		check(err)
		defer tail.Close()
		var waitGroup sync.WaitGroup
		waitGroup.Add(1)
		go hearthstone.ProcessEvents(tail.Lines, waitGroup)
		waitGroup.Wait()
	}
}

func getCardsFromFile(filename *string) []string {
	file, err := os.Open(*filename)
	check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	cards := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		cards = append(cards, line)
	}
	return cards
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
