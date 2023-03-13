package gapi

import (
	"fmt"

	db "github.com/daniel/master-golang/db/sqlc"
	"github.com/daniel/master-golang/pb"
	"github.com/daniel/master-golang/token"
	"github.com/daniel/master-golang/utils"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config utils.Config
	tokenMaker token.Marker
	store db.Store
}

func NewServer(config utils.Config, store db.Store)  (*Server, error) {
	tokenMake, err := token.NewJwMarkt(config.Token_symmetric_key)
	if err != nil {
		return nil, fmt.Errorf("cannet create token maker failed: %w", err)
	}

	server := &Server{
		store: store,
		config: config,
		tokenMaker: tokenMake,
	}

	return server, nil
}

