package jwtea

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type ThirdPartyConfiguration struct {
	Url string
	ClientId string
	ClientSecret string
	ThirdPartyAudience string
}

type Configuration struct {
	ThirdPartyConfig ThirdPartyConfiguration
	Secret string
	Issuer string
	Audience string
}

type Provider struct {
	jwtHeader header
	Config    *Configuration
}

type header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type Body struct {
	User string
	Role string
	Jti string `json:"jti"`
	Iss string `json:"iss"`
	Aud string `json:"aud"`
	Exp int64 `json:"exp"`
	Nbf int64 `json:"nbf"`
	Iat int64 `json:"iat"`
}

func NewProvider(config *Configuration) *Provider {
	jwt := &Provider{
		jwtHeader: header{"HS256", "JWT"},
		Config:    config,
	}
	return jwt
}

func (jwt *Provider) Validate(token string) bool {
	tokenParts := strings.Split(strings.Split(token, "Bearer ")[1], ".")
	jwtHeaderEncoded := tokenParts[0]
	jwtBodyEncoded := tokenParts[1]
	jwtSignatureEncoded := tokenParts[2]
	
	signatureFromToken, _ := base64.RawURLEncoding.DecodeString(jwtSignatureEncoded)
	signatureFromValidation := jwt.sign(jwtHeaderEncoded, jwtBodyEncoded)
	return hmac.Equal(signatureFromToken, signatureFromValidation)
}

func (jwt *Provider) Generate(userName string) string {
	headerBytes, _ := json.Marshal(jwt.jwtHeader)
	headerBase64 := base64.RawURLEncoding.EncodeToString(headerBytes)

	now := time.Now()
	body := &Body{
		User: userName,
		Role: "everyone",
		Jti: "",
		Iss: jwt.Config.Issuer,
		Aud: jwt.Config.Audience,
		Exp: now.Add(time.Hour * 1).Unix(),
		Nbf: now.Unix(),
		Iat: now.Unix(),
	}
	bodyBytes, _ := json.Marshal(body)
	bodyBase64 := base64.RawURLEncoding.EncodeToString(bodyBytes)

	signatureBytes := jwt.sign(headerBase64, bodyBase64)
	signatureBase64 := base64.RawURLEncoding.EncodeToString(signatureBytes)
	return fmt.Sprintf("%s.%s.%s", headerBase64, bodyBase64, signatureBase64)
}

func (jwt *Provider) sign(headerBase64 string, bodyBase64 string) []byte {
	toSign := fmt.Sprintf("%s.%s", headerBase64, bodyBase64)
	hasher := hmac.New(sha256.New, []byte(jwt.Config.Secret))
	hasher.Write([]byte(toSign))
	return hasher.Sum(nil)
}
