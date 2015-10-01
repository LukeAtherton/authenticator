// Copyright (c) Luke Atherton 2015

package authenticator

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Authenticator interface {
	Authenticate(email string, password string, uri string) (string, error)
	ValidateToken(tokenString string) bool
	GetTokenClaim(tokenString string, claim string) (value interface{}, err error)
}

type TokenAuthenticator struct {
	repo       Repo
	privateKey []byte
	publicKey  []byte
}

func BuildAuthenticator(repo Repo, privateKeyPath string, publicKeyPath string) Authenticator {
	publicKey, _ := readKeyFromFile(publicKeyPath)
	privateKey, _ := readKeyFromFile(privateKeyPath)

	return &TokenAuthenticator{
		repo:       repo,
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

func (auth *TokenAuthenticator) Authenticate(email string, password string, uri string) (string, error) {
	userId, dbErr := auth.repo.FindEmail(email)
	if dbErr != nil {
		fmt.Printf("ERROR: %s\n", dbErr.Error())
	}

	if userId == nil {
		return "", errors.New("Authentication Failed: Invalid username.")
	}

	credentials, _ := auth.repo.GetCredentials(userId)

	if credentials.Id == nil {
		return "", errors.New("Authentication Failed: Unable to find user credentials.")
	}

	if !MatchPassword(password, &PasswordKey{credentials.Salt, credentials.Key}) {
		return "", errors.New("Authentication Failed: Invalid password.")
	}

	token := jwt.New(jwt.GetSigningMethod("RS256"))

	token.Claims["id"] = userId.String()
	token.Claims["email"] = email
	token.Claims["isAccountVerified"] = credentials.IsEmailVerified
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString(auth.privateKey)

	if err != nil {
		fmt.Printf("ERROR: %v\n", err.Error())
		return "", err
	}

	return tokenString, nil
}

func (auth *TokenAuthenticator) ValidateToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return auth.publicKey, nil
	})

	if err == nil && token.Valid {
		return true
	}

	fmt.Printf("ERROR: %v", err)

	return false
}

func (auth *TokenAuthenticator) GetTokenClaim(tokenString string, claim string) (value interface{}, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return auth.publicKey, nil
	})

	if err == nil && token.Valid {
		return token.Claims[claim], nil
	}

	fmt.Printf("ERROR: %v", err)

	return nil, err
}

func readKeyFromFile(path string) (key []byte, err error) {
	file, err := os.Open(path)

	if err != nil {
		fmt.Printf("ERROR: %v\n", err.Error())
		return nil, err
	}

	defer file.Close()

	buf := new(bytes.Buffer)
	io.Copy(buf, file)
	key = buf.Bytes()

	return key, err
}
