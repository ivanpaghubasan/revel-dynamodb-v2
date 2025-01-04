package services

import (
	"revel-dynamodb-v2/app/models"
	"revel-dynamodb-v2/app/repositories"
)

type MovieService struct {
	Repo repositories.Repository
}

func NewMovieService(repo repositories.Repository) *MovieService {
	return &MovieService{Repo: repo}
}

func (c *MovieService) GetMovieByID(id string) (*models.Movie, error) {
	return c.Repo.GetMovieByID(id)
}