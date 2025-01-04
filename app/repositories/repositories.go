package repositories

import "revel-dynamodb-v2/app/models"


type Repository interface {
	GetMovieByID(id string) (*models.Movie, error)
}