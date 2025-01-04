package models



type Movie struct {
	ID string `json:"id"`
	Year string `json:"year"`
	Title string `json:"title"`
	Plot string `json:"plot"`
	Rating float32 `json:"rating"`
}