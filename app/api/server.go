package api

import (
	"fmt"

	db "github.com/daniel/master-golang/db/sqlc"
	"github.com/daniel/master-golang/token"
	"github.com/daniel/master-golang/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	router *gin.Engine
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


	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	
	server.setUpRouter()

	return server, nil
}

func (s *Server) Run(address ...string) error {
	return s.router.Run(address...)
}

func (s *Server) setUpRouter() {
	router := gin.Default()

	userGroup := router.Group("/users")
	{
		userGroup.POST("", s.createUser)
		userGroup.POST("/login", s.loginUser)
	}

	//ath.Join("/", "/") 等於 /
	authRouteGroup := router.Group("/")
	authRouteGroup.Use(authMiddleware(s.tokenMaker))
	{
		// group裡面不能/會有問題
		accountGroup :=  authRouteGroup.Group("/accounts")
		{
			accountGroup.GET("", s.getAccounts)
			accountGroup.POST("", s.createAccount)
			accountGroup.GET("/:id", s.getAccount)
		}

		transferGroup := authRouteGroup.Group("/transfers")
		{
			transferGroup.POST("", s.createTransfer)
		}
	}
	s.router = router
}
