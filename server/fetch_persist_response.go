package server

type ResponseStatus string

const (
	StatusSuccess ResponseStatus = "success"
	StatusError   ResponseStatus = "error"
)

type FetchPersistResponse struct {
	Status  ResponseStatus
	Message string `json:"message,omitempty"`
}
