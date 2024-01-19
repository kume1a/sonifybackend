package shared

type HttpErrorDto struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
