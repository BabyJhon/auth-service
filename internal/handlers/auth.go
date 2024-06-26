package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) recieve(c *gin.Context) {
	q := c.Request.URL.Query()
	guid := q.Get("guid")
	if guid == "" {
		newErrorResponse(c, http.StatusBadRequest, "empty guid query param")
	}

	accessToken, refreshToken, err := h.services.Auth.CreateTokens(context.Background(), guid)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "error occured while generate jwt tokens")
	}

	//токены в httpOnly cookie для запрета воздействия со стороны пользователя
	c.SetCookie(
		"access_token",
		accessToken,
		12341000, //через сколько кука станет недействительой в секундах
		"/",      //path
		"",       //domain
		true,     //secure
		true,     //httponly
	)
	c.SetCookie(
		"refresh_token",
		refreshToken,
		12341000,
		"/",
		"",
		true,
		true,
	)
}

func (h *Handler) refresh(c *gin.Context) {
	refreshCookie, err := c.Request.Cookie("refresh_token")
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}
	base64RefreshToken := refreshCookie.Value

	accessCookie, err := c.Request.Cookie("access_token")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	accessToken := accessCookie.Value

	//generate new tokens
	accessToken, refreshToken, err := h.services.RefreshTokens(context.Background(), accessToken, base64RefreshToken)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	c.SetCookie(
		"access_token",
		accessToken,
		12341, //через сколько кука станет недействительой в секундах
		"/",   //path
		"",    //domain
		true,  //secure
		true,  //httponly
	)
	c.SetCookie(
		"refresh_token",
		refreshToken,
		12341,
		"/",
		"",
		true,
		true,
	)
}
