package models

type ChampionData struct {
	Name                    string `json:"champ_name"`
	AverageKillParticaption int `json:"avgKP"`
	AverageCreepScore       int `json:"avgCS"`
	AverageDamage           int `json:"averageDamage"`
}
type SummonerData struct {
	SummonerName        string `json:"summonerName"`
	PUUID               string `json:"pUUID"`
	AccountID           string `json:"accountID"`
	ChampionsPlayed     []string `json:"championsPlayed"`
	ChampData           []*ChampionData `json:"champData"`
	Winrate             int `json:"winrate"`
	AmountOfGamesPlayed int `json:"amountOfGamesPlayed"`
}
type FormResponse struct {
	SummonerName string `json:"summonerName"`
	Region       string `json:"region"`
}
