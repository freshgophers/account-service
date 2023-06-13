package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/xlzd/gotp"

	"account-service/internal/domain/secret"
)

func (s *Service) SendOTP(ctx context.Context, phone string) (res secret.Response, err error) {
	data := secret.Entity{
		CreatedAt: time.Now(),
		Secret:    gotp.RandomSecret(16),
		Phone:     phone,
		Status:    &[]string{secret.ACTIVE}[0],
		Attempts:  &[]int{0}[0],
	}
	otp := gotp.NewTOTP(data.Secret, 4, 60, nil).Now()

	data.ID, err = s.secretRepository.Create(ctx, data)
	if err != nil {
		return
	}
	res = secret.ParseFromEntity(data)

	if os.Getenv("DEBUG") != "" {
		res.OTP = otp
	}

	return
}

func (s *Service) CheckOTP(ctx context.Context, req secret.Request) (err error) {
	data, err := s.secretRepository.Get(ctx, req.Key)
	if err != nil {
		return
	}

	if err = data.Validate(60, 3, req.OTP); err != nil {
		return
	}

	if err = s.secretRepository.Update(ctx, data.ID, data); err != nil {
		return
	}

	if *data.Status != secret.CONFIRMED {
		err = errors.New(data.GetText())
	}

	return
}
