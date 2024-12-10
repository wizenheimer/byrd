package models

import "time"

type Competitor struct {
	ID        int       `json:"id"`
	Domain    string    `json:"domain"`
	Name      string    `json:"name"`
	URLs      []string  `json:"urls"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CompetitorInput struct {
	Name   string   `json:"name" validate:"required"`
	Domain string   `json:"domain" validate:"required"`
	URLs   []string `json:"urls" validate:"required,min=1,dive,url"`
}

type CompetitorDTO struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Domain    string    `json:"domain"`
	URLs      []string  `json:"urls"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
