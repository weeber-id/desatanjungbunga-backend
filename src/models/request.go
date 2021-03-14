package models

// RequestOperationTime model
type RequestOperationTime struct {
	Monday struct {
		Open *bool  `json:"open" binding:"required"`
		From string `json:"from" binding:"required"`
		To   string `json:"to" binding:"required"`
	} `json:"monday" binding:"required"`
	Tuesday struct {
		Open *bool  `json:"open" binding:"required"`
		From string `json:"from" binding:"required"`
		To   string `json:"to" binding:"required"`
	} `json:"tuesday" binding:"required"`
	Wednesday struct {
		Open *bool  `json:"open" binding:"required"`
		From string `json:"from" binding:"required"`
		To   string `json:"to" binding:"required"`
	} `json:"wednesday" binding:"required"`
	Thursday struct {
		Open *bool  `json:"open" binding:"required"`
		From string `json:"from" binding:"required"`
		To   string `json:"to" binding:"required"`
	} `json:"thursday" binding:"required"`
	Friday struct {
		Open *bool  `json:"open" binding:"required"`
		From string `json:"from" binding:"required"`
		To   string `json:"to" binding:"required"`
	} `json:"friday" binding:"required"`
	Saturday struct {
		Open *bool  `json:"open" binding:"required"`
		From string `json:"from" binding:"required"`
		To   string `json:"to" binding:"required"`
	} `json:"saturday" binding:"required"`
	Sunday struct {
		Open *bool  `json:"open" binding:"required"`
		From string `json:"from" binding:"required"`
		To   string `json:"to" binding:"required"`
	} `json:"sunday" binding:"required"`
}
