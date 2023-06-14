package otp

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/xlzd/gotp"

	"account-service/internal/domain/secret"
	"account-service/internal/domain/user"
)

func (s *Service) Send(ctx context.Context, phone string) (res secret.Response, err error) {
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

	if os.Getenv("DEBUG") == "true" {
		res.OTP = otp
	} else {
		message := fmt.Sprintf("%s код подтверждения. Никому не говорите код!", otp)

		err = s.smsClient.Send(phone, message)
		if err != nil {
			return
		}
	}

	return
}

func (s *Service) Check(ctx context.Context, req secret.Request) (res user.Response, err error) {
	data, err := s.secretRepository.GetByID(ctx, req.Key)
	if err != nil {
		return
	}

	if err = data.Validate(60, 3, req.OTP); err != nil {
		return
	}

	if err = s.secretRepository.Update(ctx, data.ID, data); err != nil {
		return
	}

	switch *data.Status {
	case secret.CONFIRMED:
		res, err = s.accountService.GetOrCreateByPhone(ctx, data.Phone)
	default:
		err = errors.New(data.GetText())
	}

	return
}
