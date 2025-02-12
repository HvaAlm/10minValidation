package guard

import (
	"encoding/base64"
	"strings"
)

func CreateToken(userId string) string {
	token := base64.StdEncoding.EncodeToString([]byte(userId))
	return userId + "-Bale" + token
}
func ExtractToken(token string) string {
	userToken := strings.Split(token, "-Bale")
	userId := userToken[0]
	return userId
}
