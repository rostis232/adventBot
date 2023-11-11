package service

type Repository interface{
	GetAllCustumers()
}

type Service struct{
	Repo Repository
}

func NewService (repo Repository) *Service{
	return &Service{Repo: repo}
}
