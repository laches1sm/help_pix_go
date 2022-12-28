package models

type ChampionData struct {
	Name                    string
	AverageKillParticaption float64
	AverageCreepScore       float64
	LastRunesUsed           map[string]string
	LastItemsBuilt          map[string]string
}
type SummonerData struct{
	SummonerName string
	PUUID string
	AccountID string
	ChampionsPlayed []string
	ChampData []*ChampionData
	Winrate                 float64
	AmountOfGamesPlayed     int
}
type FormResponse struct{
	SummonerName string
	Region string
}
