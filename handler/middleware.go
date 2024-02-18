package handler

import (
	"fmt"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	middleware "github.com/oapi-codegen/echo-middleware"
)

func CreateMiddleware(v JWSValidator) ([]echo.MiddlewareFunc, error) {
	spec, err := generated.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("loading spec: %w", err)
	}

	validator := middleware.OapiRequestValidatorWithOptions(spec,
		&middleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: NewAuthenticator(v),
			},
			// SilenceServersWarning: true,
		})

	spec.Servers = nil

	return []echo.MiddlewareFunc{validator}, nil
}
