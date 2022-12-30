package api

import (
	db "github.com/daniel/master-golang/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	router *gin.Engine
	store db.Store
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}

	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	// group裡面不能/會有問題
	accountGroup :=  router.Group("/accounts")
	{
		accountGroup.GET("", server.getAccounts)
		accountGroup.POST("", server.createAccount)
		accountGroup.GET("/:id", server.getAccount)
	}

	transferGroup := router.Group("/transfers")
	{
		transferGroup.POST("/", server.createTransfer)
	}

	// router.POST("/users", server.createUser)

	userGroup := router.Group("/users")
	{
		userGroup.POST("", server.createUser)
	}

	server.router = router
	return server
}

func (s *Server) Run(address ...string) error {
	return s.router.Run(address...)
}