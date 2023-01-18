package adapters

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/KnutZuidema/golio/riot/lol"
	"github.com/laches1sm/help_pix_go/src/models"
	"github.com/sirupsen/logrus"
)

type HelpPixHTTPAdapter struct {
	*log.Logger
}

func NewHelpPixAdapter(logger *log.Logger) *HelpPixHTTPAdapter {
	return &HelpPixHTTPAdapter{
		logger
	}
}


// GetSummonerInfo accepts a POST request.
func (adapter *HelpPixHTTPAdapter) GetSummonerInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		adapter.Logger.Printf(`nope`)
		_ = marshalAndWriteErrorResponse(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		adapter.Logger.Printf(`whoops there's an error while reading request body`)
		_ = marshalAndWriteErrorResponse(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// read resp body...
	// TODO: Make this all its own function?
	// Marshal the form response into its own object.
	formResp := &models.FormResponse{}
	json.Unmarshal(body, formResp)
	// Make sure that the region is a valid Riot region.
	var region api.Region
	switch strings.ToUpper(formResp.Region) {
	case "EUW":
		region = api.RegionEuropeWest
	case "EUNE":
		region = api.RegionEuropeNorthEast
	case "NA":
		region = api.RegionNorthAmerica
	case "BR":
		region = api.RegionBrasil
	case "KR":
		region = api.RegionKorea
	case "JP":
		region = api.RegionJapan
	case "LAN":
		region = api.RegionLatinAmericaNorth
	case "LAS":
		region = api.RegionLatinAmericaSouth
	case "OCE":
		region = api.RegionOceania
	case "RU":
		region = api.RegionRussia
	default:
		adapter.Logger.Print(`Invalid region provided.`)
		_ = marshalAndWriteErrorResponse(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return

	}

	api_key := os.Getenv(`RIOT_API_KEY`)
	client := golio.NewClient(api_key, golio.WithRegion(region), golio.WithLogger(logrus.New()))
	summoner, err := client.Riot.LoL.Summoner.GetByName(formResp.SummonerName)
	if err != nil {
		adapter.Logger.Printf(`error while getting summoner: %s`, err.Error())
		_ = marshalAndWriteErrorResponse(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// the following code is a product of both me operating on zero sleep and some mank from accessing the data from riots side. anyway to make this a lot nicer would be gr8
	matchHistoryStream := client.Riot.LoL.Match.ListStream(summoner.PUUID)
	var matchIDs []string
	for k := range matchHistoryStream {
		matchIDs = append(matchIDs, k.MatchID)
	}
	var matchHistoryActual []*lol.Match
	for _, v := range matchIDs {
		match, _ := client.Riot.LoL.Match.Get(v)

		matchHistoryActual = append(matchHistoryActual, match)
	}
	var champsPlayed []string

	var wins []bool
	var winrate int
	var cs []int
	var kp []int
	var dmg []int
	var champsData []*models.ChampionData
	var botIntroGamesPlayed []string
	// while going through match history, filter for only these particular champs: Yuumi, Janna, Sona, Soraka, Ammumu, Taric, Morgana
	for _, v := range matchHistoryActual {
		// also need to check for large amount of bot games
		if v.Info.QueueID == 830 { // bot intro games
			// just check for the amount of them tbh...
			botIntroGamesPlayed = append(botIntroGamesPlayed, v.Metadata.MatchID)
		}
		// get match info here
		// make sure we get partipatnt info for the correct user
		if v.Info.QueueID == 420 { // I see what you did there riot, 420 is the id for ranked games
			id := fmt.Sprint(v.Info.GameID)
			for _, v := range v.Info.Participants {
				if v.PUUID == summoner.PUUID {
					// We only want data if it's for the champions listed above...
					if v.ChampionName == "Yuumi" || v.ChampionName == "Janna" || v.ChampionName == "Sona" || v.ChampionName == "Soraka" || v.ChampionName == "Ammumu" || v.ChampionName == "Taric" || v.ChampionName == "Morgana" {
						champsPlayed = append(champsPlayed, v.ChampionName)

						if v.Win {
							wins = append(wins, v.Win)
						}
						win := len(wins)
						winrate = win / len(matchHistoryActual) * 100

						cs = append(cs, v.TotalMinionsKilled)

						killPart := v.Kills + v.Assists
						kp = append(kp, killPart)

						dmg = append(dmg, v.TotalDamageDealt)

						wardStats := make(map[string]int)
						// And if they're in the support role, check for amount of wards bought/placed??
						if v.TeamPosition == "UTILITY" {
							wardStats["wardsKilled"] = v.WardsKilled
							wardStats["wardsPlaced"] = v.WardsPlaced
							wardStats["pinkWards"] = v.DetectorWardsPlaced
						}

						var summonerSpells []string
						summonerSpell1, _ := client.DataDragon.GetSummonerSpell(fmt.Sprint(v.Summoner1ID))
						summonerSpell2, _ := client.DataDragon.GetSummonerSpell(fmt.Sprint((v.Summoner2ID)))
						summonerSpells = append(summonerSpells, summonerSpell1.Name)
						summonerSpells = append(summonerSpells, summonerSpell2.Name)

						// TODO: do something with items here
						// get items with datadragon
						var items []string
						i1, _ := client.DataDragon.GetItem(fmt.Sprint(v.Item1))
						i2, _ := client.DataDragon.GetItem(fmt.Sprint(v.Item2))
						i3, _ := client.DataDragon.GetItem(fmt.Sprint(v.Item3))
						i4, _ := client.DataDragon.GetItem(fmt.Sprint(v.Item4))
						i5, _ := client.DataDragon.GetItem(fmt.Sprint(v.Item5))
						i6, _ := client.DataDragon.GetItem(fmt.Sprint(v.Item6))
						items = append(items, i1.Name)
						items = append(items, i2.Name)
						items = append(items, i3.Name)
						items = append(items, i4.Name)
						items = append(items, i5.Name)
						items = append(items, i6.Name)

						// For the new jungler bots, check how many objectives have been taken. This isn't perfect - I've had games as jg where I've never been able to get a single objective (thanks botlane!!)
						objectives := make(map[string]int)
						if v.TeamPosition == "JUNGLE" {
							baronKill := v.BaronKills
							drag := v.DragonKills
							dmgObj := v.DamageDealtToObjectives
							objectives["BaronKills"] = baronKill
							objectives["DragonKills"] = drag
							objectives["DamageDealtToObjectives"] = dmgObj
						}

						// place this somewhere else i guess idk i'm writing this at 3am and i'm off my meds so fuck me i guess
						champData := createNewChampData(v.ChampionName, cs, kp, dmg, v.TeamPosition, objectives, id, wardStats, summonerSpells, items)
						champsData = append(champsData, champData)

					}
				}
			}

		}

	}

	summonerData := &models.SummonerData{
		SummonerName:        summoner.Name,
		AccountID:           summoner.AccountID,
		PUUID:               summoner.PUUID,
		ChampionsPlayed:     champsPlayed,
		Winrate:             winrate,
		AmountOfGamesPlayed: len(matchHistoryActual),
		ChampData:           champsData,
		BotGamesPlayed:      len(botIntroGamesPlayed),
	}

 
	adapter.Logger.Printf(`Creating summoner %s`, summonerData.SummonerName)
	

	summonerGet, err := json.Marshal(summonerData)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	writeResponse(w, summonerGet, http.StatusOK)

	// This is some mank on my part.
	// Looking back at the JSON blob we create vs an actual CSV, I think the best way to send data to the model is one row per game.
	// for _, v := range summonerData.ChampData {
	// 	// Create a new model
	// 	_ = &models.SummonerCSV{
	// 		SummonerName:            summonerData.SummonerName,
	// 		AverageKillParticaption: v.AverageKillParticaption,
	// 		AverageCreepScore:       v.AverageCreepScore,
	// 		AverageDamage:           v.AverageDamage,
	// 		TeamPosition:            v.TeamPosition,
	// 		ObjectiveDmg:            v.ObjectiveDmg,
	// 		GameID:                  v.GameID,
	// 		WardStats:               v.WardStats,
	// 		SummonerSpells:          v.SummonerSpells,
	// 		Items:                   v.Items,
	// 		PUUID:                   summonerData.PUUID,
	// 		AccountID:               summonerData.AccountID,
	// 		ChampName:               v.Name,
	// 		Winrate:                 summonerData.Winrate,
	// 		AmountOfGamesPlayed:     summonerData.AmountOfGamesPlayed,
	// 		BotGamesPlayed:          summonerData.BotGamesPlayed,
	// 	}

		// After we successfully create our summoner, we need to convert the data into a CSV so I can feed it into the model.
		//  _, _ = json2csv.JSON2CSV(csvModel)

		// create csv file....
		// send it to go learn
		// give user csv file???
		// how to handle this idk

		// Do some stuff with GoLearn here.
		// TODO: make seperate folder for Golearn

	//}
	

	// send csv to model
	// write a different response based on what the model says

}

func createNewChampData(champName string, cs []int, killPart []int, dmg []int, teamPos string, objectives map[string]int, gameID string, wardStats map[string]int, summoners []string, items []string) *models.ChampionData {
	// calcuate avg cs
	var sumCS int
	for v := range cs {
		sumCS += v
	}
	avgCS := sumCS / len(cs)

	// calc avg KP

	var sumKP int
	for v := range killPart {
		sumKP += v
	}
	avgKP := sumKP / len(killPart)

	// calc avg dmg
	var sumDMG int
	for v := range dmg {
		sumDMG += v
	}
	avgDMG := sumDMG / len(dmg)

	return &models.ChampionData{
		Name:                    champName,
		AverageKillParticaption: avgKP,
		AverageCreepScore:       avgCS,
		AverageDamage:           avgDMG,
		TeamPosition:            teamPos,
		ObjectiveDmg:            objectives,
		GameID:                  gameID,
		WardStats:               wardStats,
		SummonerSpells:          summoners,
		Items:                   items,
	}
}
