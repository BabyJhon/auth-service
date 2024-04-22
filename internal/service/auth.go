package service

import (
	"context"
	"encoding/base64"
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/BabyJhon/auth-service/internal/entity"
	"github.com/BabyJhon/auth-service/internal/repos"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	signingKey = "something_env_secret_key" //TODO: вынести в переменные окружеия
	bcryptCost = 10
)

type AuthService struct {
	repo repos.Auth

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

type tokenClaims struct {
	jwt.RegisteredClaims
	Guid string `json:"guid"`
}

func NewAuthService(repo repos.Auth) *AuthService {
	return &AuthService{
		repo:            repo,
		accessTokenTTL:  15 * time.Minute,
		refreshTokenTTL: 48 * time.Hour,
	}
}

func (a *AuthService) CreateTokens(ctx context.Context, guid string) (string, string, error) { //acess, refresh, err
	accessToken, err := a.newAccessToken(guid, a.accessTokenTTL)
	if err != nil {
		return "", "", err
	}

	//возьмем часть от access токена для вставки в refresh токен и их связи
	linkPart := accessToken[len(accessToken)-5:]

	//для связи токенов возьмем последние
	refreshToken, err := a.newRefreshToken(ctx, guid, linkPart)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (a *AuthService) newAccessToken(guid string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "lev_osipov@auth_service",
			//ID:        guid,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Guid: guid,
	})

	return token.SignedString([]byte(signingKey))
}

func (a *AuthService) newRefreshToken(ctx context.Context, guid, linkPart string) (string, error) {
	//создание рандомного refresh token, тип- строка
	token, err := a.generateRefreshToken(linkPart)
	if err != nil {
		return "", err
	}

	//кодирование токена bcrypt и передача в монгу

	tokenHash, err := a.bcryptEncodeRefreshToken(token)
	if err != nil {
		return "", err
	}
	var session entity.Session = entity.Session{
		GUID:             guid,
		RefreshTokenHash: tokenHash,
		ExpiresAt:        time.Now().Add(a.refreshTokenTTL), //время когда истекает токен
	}

	err = a.repo.AddSession(ctx, session)
	if err != nil {
		return "", err
	}

	//кодирование токена в base64 и возврат для передачи в куки
	encodeToken := base64.StdEncoding.EncodeToString(token)

	return encodeToken, nil
}

func (a *AuthService) bcryptEncodeRefreshToken(token []byte) (string, error) {
	tokenHash, err := bcrypt.GenerateFromPassword(token, bcryptCost)
	if err != nil {
		return "", err
	}
	return string(tokenHash), nil
}

func (a *AuthService) generateRefreshToken(linkPart string) ([]byte, error) {
	token := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(token); err != nil {
		return nil, err
	}

	for i := 0; i < len(linkPart); i++ {
		token = append(token, byte(linkPart[i]))
	}

	return token, nil
}

func (a *AuthService) RefreshTokens(ctx context.Context, accessToken, base64RefreshToken string) (string, string, error) {
	/*
		1.парсим access токен
		2.декодируем из base64 refresh токен
		3.сравниваем последние 5 символов
		если равны, то

		4.по guid находим в базе сессии с хэшами рефреш токенов
		5.bcrypt.CompareHashAndPassword на рефреш токен и хэши сессий этого пользователя
		6.если нашли совпадающий то проверяем, не истек ли он
		если нет то
		7.старый рефреш удаляем из базы

		8.генерируем новый access и refresh токены
	*/

	//1.
	guid, err := a.Parsetoken(accessToken)
	if err != nil {
		return "", "", errors.New("error while parsing access token")
	}

	//2.
	//TODO: решить это
	//cookie для синтаксиса используют служебные символы, которые при передаче записываются в виде unicode сигнатуры
	normalToken1 := strings.ReplaceAll(base64RefreshToken, "%3D", "=")
	normalToken2 := strings.ReplaceAll(normalToken1, "%2F", "/")
	normalToken3 := strings.ReplaceAll(normalToken2, "%2B", "+")
	refreshTokenBytes, err := base64.StdEncoding.DecodeString(normalToken3)

	if err != nil {
		return "", "", errors.New("error while decode refresh token from base64")
	}
	refreshToken := string(refreshTokenBytes)

	//3.сравнение последних символов токенов
	if strings.Compare(refreshToken[len(refreshToken)-5:], accessToken[len(accessToken)-5:]) != 0 {
		return "", "", errors.New("access and refresh tokens are not linked")
	} //иначе последние символы access и рефреш токенов совпадают и значит они связаны

	//4.
	sessions, err := a.repo.FindSessionsByGUID(ctx, guid)
	if err != nil {
		return "", "", err
	}
	if len(sessions) == 0 {
		return "", "", errors.New("no sessions - need to auth")
	}

	//5.сверяем, могут ли быть хэши из базы получены от refresh токена
	for i := 0; i < len(sessions); i++ {
		if err := bcrypt.CompareHashAndPassword([]byte(sessions[i].RefreshTokenHash), refreshTokenBytes); err == nil {
			//6.сверяем истек ли токен
			if time.Now().Before(sessions[i].ExpiresAt) { //еще можно использовать
				//7.удаляем сессия с этим токеном из базы
				a.repo.DeleteSession(ctx, sessions[i])
				break
			} else {
				return "", "", errors.New("refresh token is expire")
			}
		}
	}

	//8.
	newAccessToken, newRefreshToken, err := a.CreateTokens(ctx, guid)
	if err != nil {
		return newAccessToken, newRefreshToken, err
	}

	return newAccessToken, newRefreshToken, nil
}

func (a *AuthService) Parsetoken(accessToken string) (string, error) { //вернет guid
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return "", nil
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", errors.New("token claims are not of type *tokenClaims")
	}

	return claims.Guid, nil
}
