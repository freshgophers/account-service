package secret

import (
	"errors"
	"time"

	"github.com/xlzd/gotp"
)

const (
	ACTIVE    = "active"
	EXPIRED   = "expired"
	DISABLED  = "disabled"
	CONFIRMED = "confirmed"
)

type Entity struct {
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	ID        string    `db:"id"`
	Secret    string    `db:"secret"`
	Phone     string    `db:"phone"`
	Attempts  *int      `db:"attempts"`
	Status    *string   `db:"status"`
}

func New(phone string) (dest Entity, otp string) {
	dest = Entity{
		CreatedAt: time.Now(),
		Secret:    gotp.RandomSecret(16),
		Phone:     phone,
		Status:    &[]string{ACTIVE}[0],
		Attempts:  &[]int{0}[0],
	}
	otp = gotp.NewTOTP(dest.Secret, 4, 60, nil).Now()

	return
}

func (e *Entity) Validate(interval int64, attempts int, otp string) (err error) {
	// check if secret is active
	if *e.Status != ACTIVE {
		return errors.New(e.GetText())
	}

	// check if secret has expired
	duration := time.Now().Unix() - e.CreatedAt.Unix()
	if duration > interval {
		e.Status = &[]string{EXPIRED}[0]
		return
	}

	// check if secret is exceeded attempts
	*e.Attempts += 1
	if *e.Attempts > attempts {
		e.Status = &[]string{DISABLED}[0]
		return
	}

	// check otp code
	valid := gotp.NewTOTP(e.Secret, 4, 60, nil).Verify(otp, e.CreatedAt.Unix())
	if valid {
		e.Status = &[]string{CONFIRMED}[0]
		return
	}

	return
}

func (e *Entity) GetText() (message string) {
	switch *e.Status {
	case ACTIVE:
		message = "otp is invalid"
	case EXPIRED:
		message = "otp has expired"
	case DISABLED:
		message = "otp is reached maximum attempts"
	case CONFIRMED:
		message = "otp is already confirmed"
	}

	return
}
