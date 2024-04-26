package utils

import (
	"errors"

	"github.com/google/uuid"
)

func GenPicRoute(cType string) (string, string, error) {
	var fileExt string

	switch cType {
	case "image/jpeg":
		fileExt = ".jpeg"
	case "image/png":
		fileExt = ".png"
	default:
		return "", "", errors.New("invalid Type")
	}
	rand := uuid.NewString()

	return "./img/" + rand + fileExt, rand, nil
}
