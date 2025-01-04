package controllers

import (
	"revel-dynamodb-v2/app"

	"github.com/revel/revel"
)

type MovieController struct {
	*revel.Controller
}

func (c MovieController) GetMovie(id string) revel.Result {
	movie, err := app.Repository.GetMovie(id)
	if err != nil {
		return c.RenderJSON(err)
	}

	return c.RenderJSON(movie)
}