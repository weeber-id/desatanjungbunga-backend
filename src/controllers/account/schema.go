package account

import "github.com/weeber-id/desatanjungbunga-backend/src/models"

type responseAdminList struct {
	Data    []*models.Admin `json:"data"`
	MaxPage uint            `json:"max_page"`
}
