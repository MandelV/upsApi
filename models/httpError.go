package models

// HttpError est la structure de réponse JSON renvoyée en cas d'erreur HTTP.
type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
