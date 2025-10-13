package consigner

type Service interface {
	Get(id int) (GetConsignerResponse, error)
	Save(PostConsignerRequest) error
}
