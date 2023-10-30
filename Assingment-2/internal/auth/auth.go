package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey int

const AuthKey ctxKey = 1

/*
*rsa.PublicKey is a type that represents a pointer to an RSA public key, and
it's used in Go when working with RSA public keys, especially when you need
to pass the key by reference or modify it in-place, rather than copying it.
*/

/*
	 By making these fields unexported, you can control access to them and ensure 
	 that only the methods or functions provided by the same package can interact 
	 with or modify these fields. This is a common practice in Go to achieve 
	 data encapsulation*/

type Auth struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

/* 
	This is a constructor function for creating instances of the Auth struct. 
	It is a common practice in Go to provide constructor functions for creating 
	instances of structs to ensure that they are properly initialized. */

func NewAuth(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) (*Auth, error) {
	/*
		It first checks if either privateKey or publicKey is nil. If either one is nil,
		it creates an error using the errors.New function and returns nil, err. */
	if privateKey == nil || publicKey == nil {
		err := errors.New("private/public key is not present")
		return nil, err
	}
	/*
		If both keys are non-nil, it proceeds to create a new instance of the Auth 
		struct and initializes its privateKey and publicKey fields with the values 
		passed as arguments. It returns the new Auth struct and a nil error to 
		indicate that the creation was successful.*/
	return &Auth{privateKey: privateKey,
		publicKey: publicKey}, nil
}

func (a *Auth) GenerateToken(claims jwt.RegisteredClaims) (string, error) {
	//NewWithClaims creates a new Token with the specified signing method and claims.
	tkn := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Signing our token with our private key.
	tokenStr, err := tkn.SignedString(a.privateKey)
	if err != nil {
		return "", fmt.Errorf("signing token %w", err)
	}

	return tokenStr, nil
}

func (a *Auth) ValidateToken(token string) (jwt.RegisteredClaims, error) {
	var c jwt.RegisteredClaims
	// Parse the token with the registered claims.
	tkn, err := jwt.ParseWithClaims(token, &c, func(token *jwt.Token) (interface{}, error) {
		return a.publicKey, nil
	})
	if err != nil {
		return jwt.RegisteredClaims{}, fmt.Errorf("parsing token %w", err)
	}
	// Check if the parsed token is valid.
	if !tkn.Valid {
		return jwt.RegisteredClaims{}, errors.New("invalid token")
	}
	return c, nil
}
