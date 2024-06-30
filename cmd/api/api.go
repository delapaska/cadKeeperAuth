package api

import (
	"database/sql"

	"github.com/delapaska/cadKeeperAuth/service/user"
	"github.com/gin-gonic/gin"
)

type APIServer struct {
	engine *gin.Engine
}

func NewAPIServer(db *sql.DB) *APIServer {
	engine := gin.New()

	userStore := user.NewStore(db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(engine)

	return &APIServer{
		engine: engine,
	}
}

func (s *APIServer) Run() {

	s.engine.Run(":8000")
}
