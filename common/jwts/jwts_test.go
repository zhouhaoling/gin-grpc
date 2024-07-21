package jwts

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"testing"
)

func TestParseToken(t *testing.T) {
	tokenStr := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJNZW1iZXJJZCI6MjE1MDE5MjE2NTEwNjg5MjgsImlzcyI6Inpob3UiLCJleHAiOjE3MTg2OTUzNTl9.tiH1KlsxK225WusxbNu2Z2Bk7y_1iI3OdQakjO2p2VF1jxGwG8Zl3BBOmx9IixF004MgdFQfLxoSwhXh8jJ7XA"
	claims, err := ParseTokenTest(tokenStr)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("mid:", claims.MemberId, "issuer:", claims.Issuer)

}

// ParseToken 解析token
func ParseTokenTest(tokenStr string) (*ProjectClaims, error) {
	claims := &ProjectClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("llhn5KCCIaSih1wJc90Xu3aRHVgqsUuB"), nil
	})
	fmt.Println("claims:", claims)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
