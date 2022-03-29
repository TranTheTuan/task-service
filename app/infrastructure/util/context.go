package util

import (
	"context"
	"errors"
)

func GetUserIDFromContext(ctx context.Context, key string) (uint32, error) {
	val := ctx.Value("UserID")
	userId, ok := val.(uint32)
	if !ok {
		return 0, errors.New("user id is invalide")
	}

	return userId, nil
}
