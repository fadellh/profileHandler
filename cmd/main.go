package main

import (
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/repository"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// // Create a fake authenticator. This allows us to issue tokens, and also
	// // implements a validator to check their validity.
	// fa, err := handler.NewFakeAuthenticator()
	// if err != nil {
	// 	log.Fatalln("error creating authenticator:", err)
	// }

	// Create middleware for validating tokens.
	// fmt.Println("<< HSAI")
	// var _ handler.JWSValidator = nil
	// fa := handler.FakeAuthenticator{}

	// mw, err := handler.CreateMiddleware(fa)
	// fmt.Println(&mw, err)
	// fmt.Println("<< OPP")
	// if err != nil {
	// 	log.Fatalln("error creating middleware:", err)
	// }
	// e.Use(middleware.Logger())
	// e.Use(mw...)

	var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, server)

	// We're going to print some useful things for interacting with this server.
	// This token allows access to any API's with no specific claims.
	// readerJWS, err := fa.CreateJWSWithClaims([]string{})
	// if err != nil {
	// 	log.Fatalln("error creating reader JWS:", err)
	// }
	// // This token allows access to API's with no scopes, and with the "things:w" claim.
	// writerJWS, err := fa.CreateJWSWithClaims([]string{"things:w"})
	// if err != nil {
	// 	log.Fatalln("error creating writer JWS:", err)
	// }

	// log.Println("Reader token", string(readerJWS))
	// log.Println("Writer token", string(writerJWS))

	e.Logger.Fatal(e.Start(":1323"))
}

func newServer() *handler.Server {
	// cfg, err := NewConfig()
	// if err != nil {
	// 	log.Error().Err(err).Msg("Failed to Initialize Configuration")
	// 	return nil, err
	// }
	// dbDsn := os.Getenv("DATABASE_URL")

	dbDsn := "postgres://postgres:postgres@127.0.0.1:5432/database?sslmode=disable"
	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})
	opts := handler.NewServerOptions{
		Repository: repo,
	}
	return handler.NewServer(opts)
}
