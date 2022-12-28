package adapters

import (
	"encoding/json"
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
	Logger log.Logger
	Infra  infrastructure.HelpPixInfra
}

func NewHelpPixAdapter(logger log.Logger, infra infrastructure.HelpPixInfra) *HelpPixHTTPAdapter {
	return &HelpPixHTTPAdapter{
		Logger: logger,
		Infra:  infra,
	}
}

var (
	regionToRealmRegion = map[api.Region]string{
		api.RegionEuropeWest:        "euw",
		api.RegionEuropeNorthEast:   "eun",
		api.RegionJapan:             "jp",
		api.RegionKorea:             "kr",
		api.RegionLatinAmericaNorth: "lan",
		api.RegionLatinAmericaSouth: "las",
		api.RegionNorthAmerica:      "na",
		api.RegionOceania:           "oce",
		api.RegionPBE:               "pbe",
		api.RegionRussia:            "ru",
		api.RegionTurkey:            "tr",
		api.RegionBrasil:            "br",
	}
)

// GibeParrot accepts only GET HTTP method. You get a parrot, you get a parrot, everyone gets a parrot!
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
	for k, v := range regionToRealmRegion {
		if strings.ToLower(formResp.Region) == v {
			adapter.Logger.Printf("Valid region found")
			region = k
		} else {
			adapter.Logger.Print(`Invalid region provided.`)
			_ = marshalAndWriteErrorResponse(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
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
	matchHistoryStream := client.Riot.LoL.Match.ListStream(summoner.PUUID, &lol.MatchListOptions{Type: "MATCHED_GAME"})
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
	// var wardStats interface{} // sort out this type later
	// while going through match history, filter for only these particular champs: Yuumi, Janna, Sona, Soraka, Ammumu, Taric, Morgana
	for _, v := range matchHistoryActual {
		// get match info here
		// make sure we get partipatnt info for the correct user
		if v.Info.GameMode == "RANKED_SOLO_5x5" {
			//var id string
			//id = string(v.Info.GameID)
			for _, v := range v.Info.Participants {
				if v.PUUID == summoner.PUUID {
					// We only want data if it's for the champions listed above...
					if v.ChampionName == "Yuumi" || v.ChampionName == "Janna" || v.ChampionName == "Sona" || v.ChampionName == "Soraka" || v.ChampionName == "Ammumu" || v.ChampionName == "Taric" || v.ChampionName == "Morgana" {
						champsPlayed = append(champsPlayed, v.ChampionName)
						// figure out winrate??
						// WIn is a boolean, so go through each true and false I guess?
						if v.Win == true {
							wins = append(wins, v.Win)
						}
						win := len(wins)
						winrate = win / len(matchHistoryActual) * 100

						// Need to fill in the ChampData struct...
						// What do we need in champ data??
						// Runes, items, avg KP and CS...
						cs = append(cs, v.TotalMinionsKilled)
						// figure out avg here
						// Riot APIs store kills and assists, just add them together i guess?
						killPart := v.Kills + v.Assists
						kp = append(kp, killPart)
						// get items, another thing came to mind, should we get info for total dmg dealt?
						dmg = append(dmg, v.TotalDamageDealt)

						// And if they're in the support role, check for amount of wards bought/placed??
						// if v.TeamPosition == "UTILITY" {
						// 	// do something with these i guess
						// 	wardStats[id]["wardsKilled"] = v.WardsKilled
						// 	wardStats[id]["wardsPlaced"] = v.WardsPlaced
						// 	wardStats[id]["pinkWards"] = v.DetectorWardsPlaced

						// }

						// TODO: do something with items here
						// rito pls i beg make getting items a lot nicer than having to perform a look up on ids
						// just lookup the most common items these bots buy ig

						// We need champdata for each champ played?
						// something like len champname create a new champ data struct? could be a nice new func
						// place this somewhere else i guess idk i'm writing this at 3am and i'm off my meds so fuck me i guess
						champData := createNewChampData(v.ChampionName, cs, kp, dmg)
						champsData = append(champsData, champData)

					}
				}
			}

		}

	}
	// start constructing our JSON doc for storage in firebase
	summonerData := &models.SummonerData{
		SummonerName:        summoner.Name,
		AccountID:           summoner.AccountID,
		PUUID:               summoner.PUUID,
		ChampionsPlayed:     champsPlayed,
		Winrate:             winrate,
		AmountOfGamesPlayed: len(matchHistoryActual),
		ChampData:           champsData,
	}

	// add infra stuff here... set up firebase in infrafolder
	parrotGet, err := json.Marshal(resp)
	if err != nil {
		return
	}
	writeResponse(w, parrotGet, http.StatusOK)
}

func createNewChampData(champName string, cs []int, killPart []int, dmg []int) *models.ChampionData {
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
	}
}
