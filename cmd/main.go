package main

import (
	"fmt"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/repository"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, server)
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
	fmt.Println(dbDsn)
	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})
	opts := handler.NewServerOptions{
		Repository: repo,
	}
	return handler.NewServer(opts)
}
