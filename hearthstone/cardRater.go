package hearthstone

import (
	"encoding/json"
	"fmt"
	"github.com/ChrisHines/GoSkills/skills"
	"github.com/ChrisHines/GoSkills/skills/trueskill"
	"os"
	"path"
	"sort"
)

type Quality struct {
	Rating      float64
	Uncertainty float64
}

func (quality *Quality) ToRating() skills.Rating {
	return skills.NewRating(quality.Rating, quality.Uncertainty)
}

type CardRater struct {
	RatingDB map[string]Quality
}

func raterPath() string {
	return path.Join(".", "ratingsDB.json")
}

func trueskillGame() *skills.GameInfo {
	return skills.DefaultGameInfo
}

func MakeCardRater() *CardRater {
	rater, err := loadRatings()
	if err != nil {
		println("WARNING: Couldn't load rating data. Starting with a new database. The error was: " + err.Error())
		rater = new(CardRater)
		rater.RatingDB = make(map[string]Quality)
	}
	return rater
}

func loadRatings() (*CardRater, error) {
	rater := new(CardRater)
	db := new(map[string]Quality)
	file, err := os.Open(raterPath())
	if err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		err = decoder.Decode(db)
		check(err)
		rater.RatingDB = *db
	}
	return rater, err
}

func (rater *CardRater) saveRatings() error {
	file, err := os.Create(raterPath())
	if err == nil {
		defer file.Close()
		encoder := json.NewEncoder(file)
		err = encoder.Encode(rater.RatingDB)
	}
	return err
}

func (rater *CardRater) UpdateRatings(winnerCards []string, loserCards []string) {
	winningTeam := skills.NewTeam()
	for _, card := range winnerCards {
		winningTeam.AddPlayer(card, rater.getCardRating(card))
	}

	losingTeam := skills.NewTeam()
	for _, card := range loserCards {
		losingTeam.AddPlayer(card, rater.getCardRating(card))
	}

	teams := []skills.Team{winningTeam, losingTeam}
	calculator := trueskill.TwoTeamCalc{}
	newRatings := calculator.CalcNewRatings(trueskillGame(), teams, 1, 2)
	for cardId, rating := range newRatings {
		rater.RatingDB[cardId.(string)] = Quality{Rating: rating.Mean(), Uncertainty: rating.Stddev()}
	}
	err := rater.saveRatings()

	check(err)
}

func (rater *CardRater) ResetRatings(cards []string) {
	for _, card := range cards {
		delete(rater.RatingDB, card)
	}
	rater.saveRatings()
}

func (rater *CardRater) RateCardGroup(cardNames []string) []string {
	ratings := make([]string, len(cardNames)+2)
	var totalSkill, totalUncertainty float64
	for i, cardName := range cardNames {
		ratings[i] = rater.GetCardQualityAsString(cardName) + " " + cardName
		skill, uncertainty := rater.GetCardQualityAsValues(cardName)
		totalSkill += skill
		totalUncertainty += uncertainty
	}
	sort.Sort(sort.Reverse(sort.StringSlice(ratings)))
	averageSkill := totalSkill / float64(len(cardNames))
	averageUncertainty := totalUncertainty / float64(len(cardNames))
	ratings[len(cardNames)] = fmt.Sprintf("Total rating: %.2f", totalSkill)
	ratings[len(cardNames)+1] = "Group average: " + rater.cardQualityString(averageSkill, averageUncertainty)
	return ratings
}

func (rater *CardRater) GetCardQualityAsString(cardName string) string {
	var rating, uncertainty float64
	if quality, ok := rater.RatingDB[cardName]; ok {
		rating = quality.Rating
		uncertainty = quality.Uncertainty
	} else {
		rating = trueskillGame().DefaultRating().Mean()
		uncertainty = trueskillGame().DefaultRating().Stddev()
	}
	return rater.cardQualityString(rating, uncertainty)
}

func (rater *CardRater) cardQualityString(mean float64, stddev float64) string {
	return fmt.Sprintf("(rating: %05.2f; uncertainty: %.2f)", mean, stddev)
}

func (rater *CardRater) GetCardQualityAsValues(cardName string) (rating, uncertainty float64) {
	if quality, ok := rater.RatingDB[cardName]; ok {
		return quality.Rating, quality.Uncertainty
	}
	defaultRating := trueskillGame().DefaultRating()
	return defaultRating.Mean(), defaultRating.Stddev()
}

func (rater *CardRater) getCardRating(cardName string) skills.Rating {
	if quality, ok := rater.RatingDB[cardName]; ok {
		return quality.ToRating()
	}
	return trueskillGame().DefaultRating()
}
