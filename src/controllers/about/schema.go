package about

import "github.com/weeber-id/desatanjungbunga-backend/src/models"

type requestAdminUpdate struct {
	Name           string `json:"name"`
	ProfilePicture string `json:"profile_picture"`
	Position       string `json:"position"`
	Body           string `json:"body"`
}

func (r *requestAdminUpdate) Write2Model(in *models.About) {
	in.Name = r.Name
	in.ProfilePicture = r.ProfilePicture
	in.Position = r.Position
	in.Body = r.Body
}
