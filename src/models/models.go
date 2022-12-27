package models

type ChampionData struct {
	Name                    string
	AverageKillParticaption float64
	AverageCreepScore       float64
	LastRunesUsed           map[string]string
	LastItemsBuilt          map[string]string
	Winrate                 float64
	AmountOfGamesPlayed     int
}
