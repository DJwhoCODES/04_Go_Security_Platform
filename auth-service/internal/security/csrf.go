package security

import "time"

type CSRFToken struct {
	Token     string
	ExpiresAt time.Time
}

func GenerateCSRFToken() (*CSRFToken, error) {

	token, err := GenerateRandomToken(32)
	if err != nil {
		return nil, err
	}

	return &CSRFToken{
		Token:     token,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}, nil
}
