package api

import (
	"database/sql"

	"github.com/delapaska/cadKeeperAuth/configs"
	_ "github.com/delapaska/cadKeeperAuth/docs"
	"github.com/delapaska/cadKeeperAuth/service/user"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type APIServer struct {
	engine *gin.Engine
}

// @title CadKeeperAuth API
// @version 1.0
// @description This is a sample API for CadKeeperAuth which includes registration, authentication, and token validation.
func NewAPIServer(db *sql.DB) *APIServer {
	engine := gin.New()

	engine.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	userStore := user.NewStore(db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(engine)

	return &APIServer{
		engine: engine,
	}
}

func (s *APIServer) Run() {

	s.engine.Run(configs.Envs.Port)
}
