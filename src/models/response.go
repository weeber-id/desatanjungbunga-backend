package models

// Meta structure for response
type Meta struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Code    uint   `json:"code"`
}

// Response structure
type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}
