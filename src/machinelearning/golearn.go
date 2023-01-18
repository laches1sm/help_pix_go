package machinelearning

import (
	"fmt"

	"github.com/laches1sm/help_pix_go/src/models"
	dataframe "github.com/rocketlaunchr/dataframe-go"
	"github.com/sjwhitworth/golearn/base"
)

// this should return a Result obj i think
// to display what the model thinks in a more user friendly manner.
func AnalyseSummoner(csv []*models.SummonerCSV) (*models.Result, error) {
	for _, v := range csv {
		var s7, s8, s9, s10 dataframe.Series
		// might need to loop through the maps
		s1 := dataframe.NewSeriesString("summoner_name", nil, v.SummonerName)
		s2 := dataframe.NewSeriesInt64("average_kill_participation", nil, v.AverageKillParticaption)
		s3 := dataframe.NewSeriesInt64("total_games_played", nil, v.AmountOfGamesPlayed)
		s4 := dataframe.NewSeriesInt64("average_cs", nil, v.AverageCreepScore)
		s5 := dataframe.NewSeriesString("position", nil, v.TeamPosition)
		s6 := dataframe.NewSeriesInt64("average_damage", nil, v.AverageDamage)
		for _, v := range v.ObjectiveDmg {
			s7 = dataframe.NewSeriesInt64("damage_to_objectives", nil, v)
		}
		for _, v := range v.WardStats {
			s8 = dataframe.NewSeriesInt64("wards", nil, v)

		}
		for _, v := range v.SummonerSpells {
			s9 = dataframe.NewSeriesString("summoner_spells", nil, v)
		}
		for _, v := range v.Items {
			s10 = dataframe.NewSeriesString("items", nil, v)
		}
		s11 := dataframe.NewSeriesInt64("winrate", nil, v.Winrate)
		s12 := dataframe.NewSeriesInt64("bot_games", nil, v.BotGamesPlayed)

		df := dataframe.NewDataFrame(s1, s2, s3, s4, s5, s6, s7, s8, s9, s10, s11, s12)

		// now that we've created the df, convert to w/e golearn uses
		instance := base.ConvertDataFrameToInstances(df, 1)
		fmt.Println(instance)
		// do some training stuff here i guess
		// first get test csv with our training dataset

	}

	return nil, nil

}
