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

func (h *Handler) zalupa(c *gin.Context) {
	// refreshCookie, err := c.Request.Cookie("refresh_token")
	// if err != nil {
	// 	newErrorResponse(c, http.StatusUnauthorized, err.Error())
	// }

	// value := refreshCookie.Value
	// fmt.Printf("refresh token from cookie value is: %s\n", value)

	// replaceString := strings.ReplaceAll(value, "%3D", "=")
	// //замена части строки
	// fmt.Printf("norm string is: %s\n", replaceString)
	c.SetCookie(
		"3d",
		"=",
		12341000, //через сколько кука станет недействительой в секундах
		"/",      //path
		"",       //domain
		true,     //secure
		true,     //httponly
	)
	c.SetCookie(
		"2f",
		"/",
		12341000, //через сколько кука станет недействительой в секундах
		"/",      //path
		"",       //domain
		true,     //secure
		true,     //httponly
	)
}

func (h *Handler) refresh(c *gin.Context) {
	refreshCookie, err := c.Request.Cookie("refresh_token")
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}
	refreshTokenHash := refreshCookie.Value

	accessCookie, err := c.Request.Cookie("access_token")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	accessToken := accessCookie.Value

	//generate new tokens
	accessToken, refreshToken, err := h.services.RefreshTokens(context.Background(), accessToken, refreshTokenHash)
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
