package api

import (
	"os"
	"testing"
	"time"

	db "github.com/daniel/master-golang/db/sqlc"
	"github.com/daniel/master-golang/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func newTestServer(t *testing.T, store db.Store) *Server {
	config := utils.Config{
		DB: utils.DB{
			Token_symmetric_key: utils.RandString(32),
			Access_token_duration: time.Minute,
		},
	}

	server, err :=NewServer(config, store)
	require.NoError(t, err)

	return server
}