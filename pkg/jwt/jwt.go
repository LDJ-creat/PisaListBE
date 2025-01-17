package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func GenerateToken(userID uint) (string, error) {
	nowTime := time.Now()
	expireHours := viper.GetInt("jwt.expire")
	if expireHours <= 0 {
		expireHours = 24 // 默认24小时
	}
	expireTime := nowTime.Add(time.Duration(expireHours) * time.Hour)

	fmt.Printf("Token generation details:\n")
	fmt.Printf("Current time: %v\n", nowTime)
	fmt.Printf("Expire hours: %d\n", expireHours)
	fmt.Printf("Expire time: %v\n", expireTime)

	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  nowTime.Unix(),
			NotBefore: nowTime.Unix(), // 添加 NotBefore 声明
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(viper.GetString("jwt.secret")))

	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		return "", err
	}

	fmt.Printf("Token generated successfully\n")
	return token, nil
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("jwt.secret")), nil
	})

	if err != nil {
		fmt.Printf("\nToken validation details:\n")
		fmt.Printf("Current time: %v\n", time.Now())
		if ve, ok := err.(*jwt.ValidationError); ok {
			fmt.Printf("Validation error type: %v\n", ve.Errors)
			claims, _ := tokenClaims.Claims.(*Claims)
			if claims != nil {
				fmt.Printf("Token expire time: %v\n", time.Unix(claims.ExpiresAt, 0))
				fmt.Printf("Token issue time: %v\n", time.Unix(claims.IssuedAt, 0))
			}
		}
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
