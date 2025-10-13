package consigner

type Service interface {
	get(id int) (GetConsignerResponse, error)
	save(PostConsignerRequest) error
}

type Repository interface {
}
