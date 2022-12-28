package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/laches1sm/help_pix_go/src/models"
	"google.golang.org/api/option"
)

type HelpPixInfra interface {
	CreateSummoner(*models.SummonerData) (*models.SummonerData, error)
	GetSummonerByName(string) (*models.SummonerData, error)
}


type firebaseRepo struct {
	firestoreClient *firestore.Client
}

func NewFirebaseRepo() (*firebaseRepo, error) {

	opt := option.WithCredentialsFile("../../lolno")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}
	client, err := app.Firestore(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error creating firestore client: %v", err)
	}
	repo := &firebaseRepo{firestoreClient: client}
	return repo, nil
}

func (repo *firebaseRepo) CreateSummoner(model *models.SummonerData) (*models.SummonerData, error) {
	_, err := repo.firestoreClient.Collection("botstobeinvestigated").Doc(model.SummonerName).Set(context.Background(), model)
	if err != nil {
		return nil, errors.New("Error while creating summoner")
	}
	return model, nil
}

func (repo *firebaseRepo) GetSummonerByName(puuid string) (*models.SummonerData, error) {
	_, err := repo.firestoreClient.Collection("botstobeinvestigated").Doc(puuid).Get(context.Background())
	if err != nil {
		return nil, errors.New("Summoner not found")
	}
	//marshal the Wrtie result back into summonerdata model
	return nil, nil
}
