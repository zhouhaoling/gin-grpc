package jwts

import (
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/golang-jwt/jwt/v5"
	"test.com/project-user/config"
)

type JwtToken struct {
	AccessToken  string
	RefreshToken string
	AccessExp    int64
	RefreshExp   int64
}

type ProjectClaims struct {
	MemberId int64
	jwt.RegisteredClaims
}

// CreateToken 创建token
func CreateToken(mid int64) (JwtToken, error) {
	c := ProjectClaims{
		MemberId: mid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, config.AppConf.Jwt.AccessExp)),
			Issuer:    config.AppConf.Jwt.Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, c)
	aToken, err := token.SignedString([]byte(config.AppConf.Jwt.Secret))
	if err != nil {
		return JwtToken{}, err
	}
	rToken, err := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, config.AppConf.Jwt.RefreshExp)),
		Issuer:    config.AppConf.Jwt.Issuer,
	}).SignedString([]byte(config.AppConf.Jwt.Secret))
	if err != nil {
		return JwtToken{}, err
	}
	//两种方法的值是一样的
	//fmt.Println("addDate:", time.Now().AddDate(0, 0, config.AppConf.Jwt.AccessExp).Unix())
	//fmt.Println("addDate2", c.ExpiresAt.Unix())
	jt := JwtToken{
		AccessToken:  aToken,
		RefreshToken: rToken,
		AccessExp:    c.ExpiresAt.Unix(),
		RefreshExp:   c.ExpiresAt.Unix(),
	}
	return jt, nil
}

// ParseToken 解析token
func ParseToken(tokenStr string) (*ProjectClaims, error) {
	claims := &ProjectClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConf.Jwt.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		zap.L().Error("token is invalid")
		return nil, errors.New("token is invalid")
	}
	return claims, nil
}
