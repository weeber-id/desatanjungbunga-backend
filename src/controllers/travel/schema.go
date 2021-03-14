package travel

import "github.com/weeber-id/desatanjungbunga-backend/src/models"

type requestCreateUpdateTravel struct {
	Name             string                      `json:"name" binding:"required"`
	Image            string                      `json:"image" binding:"required"`
	Price            string                      `json:"price" binding:"required"`
	OperationTime    models.RequestOperationTime `json:"operation_time" binding:"required"`
	ShortDescription string                      `json:"short_description" binding:"required"`
	Description      string                      `json:"description" binding:"required"`
}

func (r *requestCreateUpdateTravel) WriteToModel(in *models.Travel) {
	in.Name = r.Name
	in.Image = r.Image
	in.Price = r.Price
	in.ShortDescription = r.ShortDescription
	in.Description = r.Description

	in.OperationTime.Monday.Open = *r.OperationTime.Monday.Open
	in.OperationTime.Monday.From = r.OperationTime.Monday.From
	in.OperationTime.Monday.To = r.OperationTime.Monday.To

	in.OperationTime.Tuesday.Open = *r.OperationTime.Tuesday.Open
	in.OperationTime.Tuesday.From = r.OperationTime.Tuesday.From
	in.OperationTime.Tuesday.To = r.OperationTime.Tuesday.To

	in.OperationTime.Wednesday.Open = *r.OperationTime.Wednesday.Open
	in.OperationTime.Wednesday.From = r.OperationTime.Wednesday.From
	in.OperationTime.Wednesday.To = r.OperationTime.Wednesday.To

	in.OperationTime.Thursday.Open = *r.OperationTime.Thursday.Open
	in.OperationTime.Thursday.From = r.OperationTime.Thursday.From
	in.OperationTime.Thursday.To = r.OperationTime.Thursday.To

	in.OperationTime.Friday.Open = *r.OperationTime.Friday.Open
	in.OperationTime.Friday.From = r.OperationTime.Friday.From
	in.OperationTime.Friday.To = r.OperationTime.Friday.To

	in.OperationTime.Saturday.Open = *r.OperationTime.Saturday.Open
	in.OperationTime.Saturday.From = r.OperationTime.Saturday.From
	in.OperationTime.Saturday.To = r.OperationTime.Saturday.To

	in.OperationTime.Sunday.Open = *r.OperationTime.Sunday.Open
	in.OperationTime.Sunday.From = r.OperationTime.Sunday.From
	in.OperationTime.Sunday.To = r.OperationTime.Sunday.To
}
