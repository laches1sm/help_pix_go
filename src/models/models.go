package models

type ChampionData struct {
	Name                    string         `json:"champ_name"`
	AverageKillParticaption int            `json:"avgKP"`
	AverageCreepScore       int            `json:"avgCS"`
	AverageDamage           int            `json:"averageDamage"`
	TeamPosition            string         `json:"teamPosition"`
	ObjectiveDmg            map[string]int `json:"objectiveDamage"`
	GameID                  string         `json:"gameID"`
	WardStats               map[string]int `json:"wardStats"`
	SummonerSpells          []string       `json:"summonerSpells"`
	Items                   []string       `json:"items"`
}
type SummonerData struct {
	SummonerName        string          `json:"summonerName"`
	PUUID               string          `json:"pUUID"`
	AccountID           string          `json:"accountID"`
	ChampionsPlayed     []string        `json:"championsPlayed"`
	ChampData           []*ChampionData `json:"champData"`
	Winrate             int             `json:"winrate"`
	AmountOfGamesPlayed int             `json:"amountOfGamesPlayed"`
	BotGamesPlayed      int             `json:"botGamesPlayed"`
}
type FormResponse struct {
	SummonerName string `json:"summonerName"`
	Region       string `json:"region"`
}

type SummonerCSV struct {
	SummonerName            string         `json:"summonerName"`
	AverageKillParticaption int            `json:"avgKP"`
	AverageCreepScore       int            `json:"avgCS"`
	AverageDamage           int            `json:"averageDamage"`
	TeamPosition            string         `json:"teamPosition"`
	ObjectiveDmg            map[string]int `json:"objectiveDamage"`
	GameID                  string         `json:"gameID"`
	WardStats               map[string]int `json:"wardStats"`
	SummonerSpells          []string       `json:"summonerSpells"`
	Items                   []string       `json:"items"`
	PUUID                   string         `json:"pUUID"`
	AccountID               string         `json:"accountID"`
	ChampName               string         `json:"champName"`
	Winrate                 int            `json:"winrate"`
	AmountOfGamesPlayed     int            `json:"amountOfGamesPlayed"`
	BotGamesPlayed          int            `json:"botGamesPlayed"`
	isBot                   bool           `json:"isBot"` // N.B. this field is ONLY for training the model. 
}
