package model

import (
	"encoding/json"
	"strconv"
)

type FlexInt int

func (fi *FlexInt) UnmarshalJSON(b []byte) error {
	// Try normal number first
	var n int
	if err := json.Unmarshal(b, &n); err == nil {
		*fi = FlexInt(n)
		return nil
	}

	// Fall back to string
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*fi = FlexInt(n)
	return nil
}

type Item struct {
	SlNo     int     `json:"slNo"`
	Item     string  `json:"item" binding:"required"`
	MRPOrNet float64 `json:"mrpOrNet" binding:"required,gt=0"`
	Quantity FlexInt `json:"quantity" binding:"required,gt=0"`
	Discount string  `json:"discount"`
	SubTotal float64 `json:"subTotal" binding:"required,gt=0"`
}