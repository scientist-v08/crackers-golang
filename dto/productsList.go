package dto

type ProductsList struct {
	Brand string     `json:"brand"`
	List  []Products `json:"list"`
}