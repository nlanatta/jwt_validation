package main

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
)

type TokenResponse struct {
	Token string `json:"token"`
}

func CreateTokenResponse(uid int64) (*TokenResponse, error) {
	token, err := GenerateToken(strconv.FormatUint(uint64(uid), 10), []string{"USER"}, reflect.TypeOf(TokenResponse{}))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("generating token: %v", err))
	}

	toReturn := &TokenResponse{
		Token: token,
	}
	return toReturn, nil
}
