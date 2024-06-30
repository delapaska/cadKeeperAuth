package user

import (
	"fmt"
	"net/http"

	"github.com/delapaska/cadKeeperAuth/configs"
	"github.com/delapaska/cadKeeperAuth/models"
	"github.com/delapaska/cadKeeperAuth/service/auth"
	"github.com/delapaska/cadKeeperAuth/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	store models.UserStore
}

func NewHandler(store models.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {

	api := router.Group("/api/v1")
	{
		api.POST("/login", h.handleLogin)
		api.POST("/register", h.handleRegister)
		api.GET("/locked", auth.WithJWTAuth(h.locked, configs.Envs.JWTSecret))
	}

}
func (h *Handler) locked(c *gin.Context) {
	c.JSON(200, "unlocked")
}
func (h *Handler) handleLogin(c *gin.Context) {
	var payload models.LoginUserPayload
	if err := utils.ParseJSON(c, &payload); err != nil {

		utils.WriteError(c, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(c, http.StatusBadRequest, fmt.Errorf("ivalid payload %v", errors))
		return
	}
	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(c, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}
	if !auth.ComparePasswords(u.Password, []byte(payload.Password)) {
		utils.WriteError(c, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}
	secret := []byte(configs.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		utils.WriteError(c, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(c, http.StatusOK,
		gin.H{"token": token})
}

func (h *Handler) handleRegister(c *gin.Context) {

	var payload models.RegisterUserPayload
	if err := utils.ParseJSON(c, &payload); err != nil {

		utils.WriteError(c, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(c, http.StatusBadRequest, fmt.Errorf("ivalid payload %v", errors))
		return
	}

	_, err := h.store.GetUserByEmail(payload.Email)

	if err == nil {
		utils.WriteError(c, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	err = h.store.CreateUser(models.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: hashedPassword,
	})
	if err != nil {
		utils.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(c, http.StatusCreated, gin.H{
		"status": "OK",
		"code":   http.StatusCreated,
	})
}
