package converter

import (
	"test-be-kalbe/internal/domain/model"
	"time"
)

func AuthToResponse(token string, expires time.Time) *model.LoginResponse {
	return &model.LoginResponse{
		Token:     token,
		ExpiresIn: expires.Format("2006-01-02 15:04:05"),
	}
}
