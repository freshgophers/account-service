package account

import (
	"context"

	"account-service/internal/domain/user"
	"account-service/pkg/store"
)

func (s *Service) GetOrCreateByPhone(ctx context.Context, phone string) (res user.Response, err error) {
	data, err := s.userRepository.GetByPhone(ctx, phone)
	if err != nil && err != store.ErrorNotFound {
		return
	}

	if err == store.ErrorNotFound {
		data = user.New(phone, user.CUSTOMER)

		data.ID, err = s.userRepository.Create(ctx, data)
		if err != nil {
			return
		}
	}

	data, err = s.userRepository.GetByID(ctx, data.ID)
	if err != nil {
		return
	}
	res = user.ParseFromEntity(data)

	return
}

func (s *Service) GetByID(ctx context.Context, id string) (res user.Response, err error) {
	data, err := s.userRepository.GetByID(ctx, id)
	if err != nil {
		return
	}
	res = user.ParseFromEntity(data)

	return
}

func (s *Service) Update(ctx context.Context, id string, req user.Request) (err error) {
	data := user.Entity{
		Name:      req.Name,
		Email:     req.Email,
		BirthDate: req.BirthDate,
	}
	return s.userRepository.Update(ctx, id, data)
}
