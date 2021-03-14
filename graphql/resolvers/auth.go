package resolvers

import (
	"github.com/dgrijalva/jwt-go"
	"leggett.dev/devmarks/api/model"
)

type AuthPayload struct {
	Token *string
	User *model.User
}

type AuthResolver struct {
	AuthPayload AuthPayload
}

func GenerateToken(user *model.User, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID": user.ID,
	})
	tokenString, errToken := token.SignedString(secret)
	if errToken != nil {
		return "", errToken
	}
	return tokenString, nil
}

func (r *AuthResolver) Token() *string {
	return r.AuthPayload.Token
}

func (r *AuthResolver) User() *UserResolver {
	return &UserResolver{*r.AuthPayload.User}
}