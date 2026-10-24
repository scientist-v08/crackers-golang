package dto

import "github.com/scientist-v08/crackers/db"

type Routes struct {
	Id      uint   `json:"id"`
	Route   string `json:"route"`
	Heading string `json:"heading"`
	Role    string `json:"role"`
}

func ToRoute(r db.Route) Routes {
	return Routes{
		Id:      uint(r.ID),
		Route:   r.Route,
		Heading: r.Heading,
		Role:    r.Role,
	}
}

func ToRoutes(routes []db.Route) []Routes {
	result := make([]Routes, len(routes))
	for i, r := range routes {
		result[i] = ToRoute(r)
	}
	return result
}