package models

import "net/http"

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

// ErrorBadRequest form
func (r *Response) ErrorBadRequest(message string) *Response {
	r.Meta.Message = message
	r.Meta.Status = "bad request"
	r.Meta.Code = http.StatusBadRequest

	r.Data = nil
	return r
}

// ErrorForbidden form
func (r *Response) ErrorForbidden() *Response {
	r.Meta.Message = "Forbidden access"
	r.Meta.Status = "forbidden"
	r.Meta.Code = http.StatusForbidden

	r.Data = nil
	return r
}

// ErrorDataNotFound form
func (r *Response) ErrorDataNotFound() *Response {
	r.Meta.Message = "data tidak ditemukan"
	r.Meta.Status = "not found"
	r.Meta.Code = http.StatusNotFound

	r.Data = nil
	return r
}

// ErrorInternalServer form
func (r *Response) ErrorInternalServer(err error) *Response {
	r.Meta.Message = err.Error()
	r.Meta.Status = "internal server error"
	r.Meta.Code = http.StatusInternalServerError

	r.Data = nil
	return r
}

// SuccessDataList form
func (r *Response) SuccessDataList(data interface{}) *Response {
	r.Meta.Message = "data berhasil di proses"
	r.Meta.Status = "ok"
	r.Meta.Code = http.StatusOK

	// handle empty array became [], not null
	r.Data = data
	return r
}

// SuccessData form
func (r *Response) SuccessData(data interface{}) *Response {
	r.Meta.Message = "data berhasil di proses"
	r.Meta.Status = "ok"
	r.Meta.Code = http.StatusOK

	r.Data = data
	return r
}

// SuccessDataCreated form
func (r *Response) SuccessDataCreated(data interface{}) *Response {
	r.Meta.Message = "data berhasil di buat"
	r.Meta.Status = "created"
	r.Meta.Code = http.StatusCreated

	r.Data = data
	return r
}
