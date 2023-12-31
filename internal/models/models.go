package models

type Costumer struct{
	ChatID *int
	Name *string
	IsActivated *int
	WaitingFor *int
}

type Message struct{
	MessageID int
	DateTime string
	Text string
	Sent int
}

type SecretKey struct{
	SkID *int
	SecretKey *int
	ChatID *int
}