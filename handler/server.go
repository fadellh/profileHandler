package handler

import (
	"sync"

	"github.com/SawitProRecruitment/UserService/repository"
)

type Server struct {
	Repository repository.RepositoryInterface
	sync.RWMutex
}

type NewServerOptions struct {
	Repository repository.RepositoryInterface
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		Repository: opts.Repository,
	}
}
