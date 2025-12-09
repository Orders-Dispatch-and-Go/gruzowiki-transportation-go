package services

import (
	"context"
	"fmt"
	"gruzowiki/repositories"
	"gruzowiki/rest/models"
	"gruzowiki/rest/terror"
)

type RecipientService struct {
	repo *repositories.RecipientRepo
}

func NewRecipientService(repo *repositories.RecipientRepo) *RecipientService {
	return &RecipientService{repo: repo}
}

func (s *RecipientService) CreateRecipient(ctx context.Context, firstName, secondName, thirdName, phone, email string) (*models.CreateRecipientResponse, error) {
	newID, err := s.repo.CreateRecipient(ctx, firstName, secondName, thirdName, phone, email)
	if err != nil {
		return nil, fmt.Errorf("CreateRecipient: %w", err)
	}

	resp := &models.CreateRecipientResponse{ID: newID}
	return resp, nil
}

func (s *RecipientService) GetRecipient(ctx context.Context, id int32) (*models.GetRecipientResponse, error) {
	recipient, err := s.repo.GetRecipientByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if recipient == nil {
		return nil, terror.NewNotFoundError("Recipient", fmt.Sprint(id))
	}
	return &models.GetRecipientResponse{
		ID:         recipient.ID,
		FirstName:  recipient.FirstName.String,
		SecondName: recipient.SecondName.String,
		ThirdName:  recipient.ThirdName.String,
		Phone:      recipient.Phone.String,
		Email:      recipient.Email.String,
	}, nil
}

func (s *RecipientService) ListRecipients(ctx context.Context) ([]*models.GetRecipientResponse, error) {
	recipients, err := s.repo.ListRecipients(ctx)
	if err != nil {
		return nil, err
	}

	var res []*models.GetRecipientResponse
	for _, r := range recipients {
		res = append(res, &models.GetRecipientResponse{
			ID:         r.ID,
			FirstName:  r.FirstName.String,
			SecondName: r.SecondName.String,
			ThirdName:  r.ThirdName.String,
			Phone:      r.Phone.String,
			Email:      r.Email.String,
		})
	}
	return res, nil
}

func (s *RecipientService) UpdateRecipient(ctx context.Context, id int32, firstName, secondName, thirdName, phone, email string) (*models.UpdateRecipientResponse, error) {
	updatedID, err := s.repo.UpdateRecipient(ctx, id, firstName, secondName, thirdName, phone, email)
	if err != nil {
		return nil, err
	}
	return &models.UpdateRecipientResponse{ID: updatedID}, nil
}

func (s *RecipientService) DeleteRecipient(ctx context.Context, id int32) error {
	return s.repo.DeleteRecipient(ctx, id)
}
