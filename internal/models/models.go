package models

type Costumer struct{
	CostumerID *int
	ChatID *int
	Name *string
	Status *int
}

type Message struct{
	MessageID int
	DateTime string
	Text string
	Sent int
}