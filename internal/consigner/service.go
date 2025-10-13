package consigner

type service struct {
	repo Repository
}

func (s service) get(id int) (GetConsignerResponse, error) {
	return GetConsignerResponse{}, nil
}

func (s service) save(request PostConsignerRequest) int {
	return 0
}
