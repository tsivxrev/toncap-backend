package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"toncap-backend/types"
)

func GenerateToken(tokenData types.TokenData) (token string, err error) {
	hmacHash := hmac.New(sha256.New, []byte(os.Getenv("SECRET")))
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
	hmacHash := hmac.New(sha256.New, []byte(os.Getenv("SECRET")))
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

		now := time.Now()
		expiresIn := time.Unix(tokenData.ExpiresIn, 0)
		if now.After(expiresIn) {
			return tokenData, false
		}

		return tokenData, true
	}

	return nil, false
}
