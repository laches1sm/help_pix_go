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
	if r.Method != http.MethodGet {
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
	formResp := &models.FormResponse{
		Name: body["summoner_name"],
		Region: body["region"],
	}
	// Make sure that the region is a valid Riot region.
	var region api.Region
	for k, v := range regionToRealmRegion{
		if strings.ToLower(formResp.Region) == v{
			adapter.Logger.Printf("Valid region found")
			region = k
		}else{
			adapter.Logger.Print(`Invalid region provided.`)
			_ = marshalAndWriteErrorResponse(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		    return
		}
	}
    api_key := os.Getenv(`RIOT_API_KEY`)
	client := golio.NewClient(api_key, golio.WithRegion(region), golio.WithLogger(logrus.New()))
	summoner, err := client.Riot.LoL.Summoner.GetByName(formResp.SummonerName)
	if err != nil{
		adapter.Logger.Printf(`error while getting summoner: %s`, err.Error())
		_ = marshalAndWriteErrorResponse(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	matchHistoryStream := client.Riot.LoL.Match.ListStream(summoner.PUUID, &lol.MatchListOptions{Type: "MATCHED_GAME"})
	var matchIDs []string
	for k := range matchHistoryStream{
		 matchIDs = append(k.MatchID)
	}
	var matchHistoryActual []*lol.Match
	for k, v := range matchIDs{
		matchHistoryActual = append(client.Riot.LoL.Match.Get(v))
	}
	var champsPlayed []string

	var wins []bool
	var winrate float64
	var cs []int
	var kp[] int
	// while going through match history, filter for only these particular champs: Yuumi, Janna, Sona, Soraka, Ammumu, Taric, Morgana
	for _, v := range matchHistoryActual{
		// get match info here
		// make sure we get partipatnt info for the correct user
		if v.Info.QueueID == "RANKED_SOLO_5x5"{
			for _, v := range v.Info.Participants{
				if v.PUUID == summoner.PUUID{
					// let's get some data boiiiiii
					// We only want data if it's for the champions listed above...
					if v.ChampionName == "Yuumi" || "Janna" || "Sona" || "Soraka" || "Ammumu" || "Taric" || "Morgana"{
						champsPlayed = append(v.ChampionName)
						// figure out winrate??
						// WIn is a boolean, so go through each true and false I guess?
						if v.Win == true{
							wins = append(wins, v.Win)
						} 
						win := len(wins)
						winrate = win / len(matchHistoryActual) * 100

						// Need to fill in the ChampData struct...
						// We need champdata for each champ played?
						// something like len champname create a new champ data struct? could be a nice new func
						// place this somewhere else i guess idk i'm writing this at 3am and i'm off my meds so fuck me i guess
						champData := createNewChampData()
						// What do we need in champ data??
						// Runes, items, avg KP and CS...
						cs = append(v.TotalMinionsKilled) 
						// figure out avg here
						// Riot APIs store kills and assists, just add them together i guess?
						killPart := v.Kills + v.Assists
						kp = append(kp, killPart)
						// get items, another thing came to mind, should we get info for total dmg dealt?
						// And if they're in the support role, check for amount of wards bought/placed??
						if v.TeamPosition == "UTILITY"{
							// do something with these i guess
							wardsKilled := v.WardsKilled
							wardsPlaced := v.WardsPlaced
							pinkWards := v.DetectorWardsPlaced
						}

						
						
					
						
					}
				}
			}

		}

	}
	// filter only for ranked games.
	// loop through match history
	// need to get last 10 matches
	// What needs to be collected here?
	// TODO: calulate Winrate, and how to attach several champs played to champ data
	summonerData := &models.SummonerData{
		SummonerName: summoner.Name,
		AccountID: summoner.AccountID,
		PUUID: summoner.PUUID,
		ChampionsPlayed: champsPlayed,
		Winrate: winrate,
		AmountOfGamesPlayed: len(matchHistory),
		ChampData: [&models.ChampionData{}],
	}
	// After getting summoner, we need to get their match history 
	parrotGet, err := json.Marshal(resp)
	if err != nil {
		return
	}
	writeResponse(w, parrotGet, http.StatusOK)
}

func (adapter *ParrotHTTPAdapter) AddParrot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		adapter.Logger.Printf(`lol_nope`)
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
	resp, err := adapter.Infra.AddParrot(body)
	if err != nil {
		return
	}
	parrotAdd, err := json.Marshal(resp)
	if err != nil {
		return
	}
	writeResponse(w, parrotAdd, http.StatusCreated)
}

func createNewChampData() *models.ChampionData{
	return nil
}