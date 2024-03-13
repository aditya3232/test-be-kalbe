package application

import (
	"test-be-kalbe/internal/domain/model"
	"test-be-kalbe/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AuthApplication struct {
	AuthService service.AuthService
	Log         *logrus.Logger
}

func NewAuthApplication(authService service.AuthService, log *logrus.Logger) *AuthApplication {
	return &AuthApplication{
		AuthService: authService,
		Log:         log,
	}
}

func (a *AuthApplication) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginRequest)
	if err := ctx.BodyParser(request); err != nil {
		a.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	response, err := a.AuthService.Login(ctx.UserContext(), request)
	if err != nil {
		a.Log.WithError(err).Error("error login")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.LoginResponse]{
		Meta: model.Meta{Code: 200, Status: "success login"},
		Data: response,
	})
}

func (a *AuthApplication) Logout(ctx *fiber.Ctx) error {
	request := &model.LogoutRequest{
		Token: ctx.Get("Authorization"),
	}

	if request.Token == "" {
		a.Log.Error("error parsing request header")
		return fiber.ErrUnauthorized
	}

	response, err := a.AuthService.Logout(ctx.UserContext(), request)
	if err != nil {
		a.Log.WithError(err).Error("error logout")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.LogoutResponse]{
		Meta: model.Meta{Code: 200, Status: "success logout"},
		Data: response,
	})
}
