package controllers

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"toncap-backend/config"
	"toncap-backend/types"

	"github.com/gin-gonic/gin"
)

func GenerateToken(tokenData types.TokenData) (token string, err error) {
	hmacHash := hmac.New(sha256.New, []byte(config.HASH_SECRET))
	md5Hash := md5.New()

	tokenDataJSON, err := json.Marshal(tokenData)
	if err != nil {
		return "", err
	}

	tokenDataEncoded := base64.RawStdEncoding.EncodeToString(tokenDataJSON)

	hmacHash.Write(tokenDataJSON)
	hmacSum := hmacHash.Sum(nil)

	md5Hash.Write(hmacSum)
	md5Sum := md5Hash.Sum(nil)

	sign := hex.EncodeToString(md5Sum)

	return fmt.Sprintf("%s.%s", tokenDataEncoded, sign), nil
}

func ValidateToken(token string) (tokenData *types.TokenData, valid bool) {
	hmacHash := hmac.New(sha256.New, []byte(config.HASH_SECRET))
	md5Hash := md5.New()

	splittedToken := strings.Split(token, ".")

	if len(splittedToken) < 2 {
		return nil, false
	}

	encodedData := splittedToken[0]
	sign := splittedToken[1]

	decodedData, err := base64.RawStdEncoding.DecodeString(encodedData)
	if err != nil {
		return nil, false
	}

	hmacHash.Write(decodedData)
	md5Hash.Write(hmacHash.Sum(nil))

	if sign == hex.EncodeToString(md5Hash.Sum(nil)) {
		err = json.Unmarshal([]byte(decodedData), &tokenData)
		if err != nil {
			return nil, false
		}

		return tokenData, true
	}

	return nil, false
}

func GenerateTokenController(c *gin.Context) {
	if c.GetString("auth_token_type") != "service" {
		NewError(c, 403, errors.New("access denied"))
		return
	}

	tokenData := types.TokenData{}

	err := c.ShouldBindJSON(&tokenData)
	if err != nil {
		NewError(c, 400, err)
		return
	}

	token, err := GenerateToken(tokenData)
	if err != nil {
		NewError(c, 400, err)
	}

	c.JSON(200, gin.H{
		"data":  tokenData,
		"token": token,
	})
}

func ValidateTokenController(c *gin.Context) {
	if c.GetString("auth_token_type") != "service" {
		NewError(c, 403, errors.New("access denied"))
		return
	}

	token := c.Param("token")

	tokenData, valid := ValidateToken(token)
	if !valid {
		c.JSON(200, gin.H{
			"valid": valid,
		})
		return
	}

	c.JSON(200, gin.H{
		"valid": valid,
		"data":  tokenData,
	})
}
