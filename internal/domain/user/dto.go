package user

import (
	"net/http"
	"time"
)

type Request struct {
	Name      *string    `json:"name"`
	Email     *string    `json:"email"`
	BirthDate *time.Time `json:"birth_date"`
}

func (s *Request) Bind(r *http.Request) error {

	return nil
}

type Response struct {
	ID        string     `json:"id"`
	Phone     string     `json:"phone"`
	Type      string     `json:"type"`
	Name      *string    `json:"name"`
	Email     *string    `json:"email"`
	BirthDate *time.Time `json:"birth_date"`
}

func ParseFromEntity(data Entity) (res Response) {
	res = Response{
		ID:        data.ID,
		Phone:     data.Phone,
		Type:      data.Type,
		Name:      data.Name,
		Email:     data.Email,
		BirthDate: data.BirthDate,
	}
	return
}
