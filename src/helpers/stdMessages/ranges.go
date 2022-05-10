package stdMessages

// This file defines standard results for responses of a slice of a resource

// generic range data
type RangeData struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Amount  int         `json:"amount"`
	Offset  int         `json:"offset"`
	Limmit  int         `json:"limmit"`
	Next    string      `json:"next"`
	Data    interface{} `json:"data"`
}
