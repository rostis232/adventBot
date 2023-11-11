package handler

type Service interface{}

type Handler struct {
	Service Service
}

func NewHandler (service Service) *Handler {
	return &Handler{Service: service}
}