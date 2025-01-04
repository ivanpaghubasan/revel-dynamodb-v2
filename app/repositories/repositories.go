package repositories

import "revel-dynamodb-v2/app/models"


type Repository interface {
	GetMovie(id string) (*models.Movie, error)
}