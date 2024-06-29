package utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/iamrz1/rutils"
)

const (
	accessTokenKey   = "jdskhhiuewfhosfkaskfhajksfeiwhfuiowehfiwejdkfewudhuiewhjfdiu"
	refreshTokenKey  = "fdshfjdshfjhdsjlfhuoashfuherifherhfuqheruifhiquwhfukwjnfjiwhl"
	RoleUser         = "user"
	RoleManger       = "manager"
	AuthorizationKey = "authorization"
	UsernameKey      = "username"
	RoleKey          = "role"
)

type claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func GenerateTokens(username, role string) (string, string) {
	accessTokenValidity := 60
	refreshTokenValidity := 7 * 24 * 60

	now := time.Now().UTC()
	accessExpTime := now.Add(time.Minute * time.Duration(accessTokenValidity))
	refreshExpTime := now.Add(time.Minute * time.Duration(refreshTokenValidity))

	accessClaims := &claims{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: accessExpTime.Unix(),
		},
	}

	refreshClaims := &claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: refreshExpTime.Unix(),
		},
	}

	jwtAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	jwtRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessTokenString, _ := jwtAccessToken.SignedString([]byte(accessTokenKey))
	refreshTokenString, _ := jwtRefreshToken.SignedString([]byte(refreshTokenKey))

	return accessTokenString, refreshTokenString

}

func verifyToken(token string, isRefresh bool) (*claims, error) {
	tokenClaims := &claims{}
	secretKey := accessTokenKey
	if isRefresh {
		secretKey = refreshTokenKey
	}

	tkn, err := jwt.ParseWithClaims(token, tokenClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		return nil, fmt.Errorf("%s", "Invalid token")
	}

	return tokenClaims, nil
}

func ManagerOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtTkn := r.Header.Get(AuthorizationKey)
		if jwtTkn == "" {
			rutils.HandleObjectError(w, rutils.NewGenericError(http.StatusUnauthorized, "Missing access token"))
			return
		}

		jwtTkn = stripBearerFromToken(jwtTkn)

		jwtClaims, err := verifyToken(jwtTkn, false)
		if err != nil {
			rutils.HandleObjectError(w, rutils.NewGenericError(http.StatusUnauthorized, err.Error()))
			return
		}

		if jwtClaims.Role != RoleManger {
			rutils.HandleObjectError(w, rutils.NewGenericError(http.StatusForbidden, "Manager only"))
			return
		}

		setTokenClaimsInHeader(r, jwtClaims.Username, jwtClaims.Role)

		next.ServeHTTP(w, r)
	})
}

func UserOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtTkn := r.Header.Get(AuthorizationKey)
		if jwtTkn == "" {
			rutils.HandleObjectError(w, rutils.NewGenericError(http.StatusUnauthorized, "Missing access token"))
			return
		}

		jwtTkn = stripBearerFromToken(jwtTkn)

		jwtClaims, err := verifyToken(jwtTkn, false)
		if err != nil {
			rutils.HandleObjectError(w, rutils.NewGenericError(http.StatusUnauthorized, err.Error()))
			return
		}

		if jwtClaims.Role != RoleUser {
			rutils.HandleObjectError(w, rutils.NewGenericError(http.StatusForbidden, "User only"))
			return
		}

		setTokenClaimsInHeader(r, jwtClaims.Username, jwtClaims.Role)

		next.ServeHTTP(w, r)
	})
}

func stripBearerFromToken(token string) string {
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
		//log.Println("token contains Bearer")
	}

	if strings.HasPrefix(token, "bearer ") {
		token = strings.TrimPrefix(token, "bearer ")
		//log.Println("token contains bearer")
	}

	return token
}

func setTokenClaimsInHeader(r *http.Request, username, role string) {
	r.Header.Set(UsernameKey, username)
	r.Header.Set(RoleKey, role)
}
