package consigner

type GetConsignerResponse struct {
	id             int
	driverCategory string
}

type PostConsignerRequest struct {
}

type PostConsignerResponse struct {
	id int
}
