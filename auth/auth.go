package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"leggett.dev/devmarks/api/app"
	"leggett.dev/devmarks/api/model"
)

func AuthenticateToken(ctx context.Context, app app.App) (*model.User, error) {
	var token, ok = ctx.Value("token").(string)
	if !ok {
		return nil, errors.New("bookmarks: no authenticated user in context")
	}

	user, err := getUserFromToken(token, app)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func getUserFromToken(tokenString string, app app.App) (*model.User, error) {
	// decode token with the secret if was encoded with
	tokenObj, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return app.Config.SecretKey, nil
	})
	if err != nil {
		return nil, err
	}
	// get user ID from the map we encoded in the token
	userID, ok := tokenObj.Claims.(jwt.MapClaims)["ID"].(float64)
	if !ok {
		return nil, errors.New("GetUserIDFromToken error: type conversion in claims")
	}

	user, err := app.Database.GetUserById(uint(userID))

	if err != nil {
		return nil, errors.New("No user with ID " + fmt.Sprint(userID))
	}

	return user, nil
}