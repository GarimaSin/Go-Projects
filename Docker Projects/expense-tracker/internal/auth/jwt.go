package auth

import (
    "time"
    "errors"

    "github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
    secret string
    expiry time.Duration
}

type Claims struct {
    UserID int64 `json:"user_id"`
    jwt.RegisteredClaims
}

func NewJWTManager(secret string, expiry time.Duration) *JWTManager {
    return &JWTManager{secret: secret, expiry: expiry}
}

func (j *JWTManager) Generate(userID int64) (string, error) {
    claims := Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiry)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(j.secret))
}

func (j *JWTManager) Verify(tokenStr string) (*Claims, error) {
    tkn, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(j.secret), nil
    })
    if err != nil {
        return nil, err
    }
    if claims, ok := tkn.Claims.(*Claims); ok && tkn.Valid {
        return claims, nil
    }
    return nil, errors.New("invalid token")
}
