package jwtea

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
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

func (jwt *Provider) Validate(token string) error {
	bearerPrefix := "Bearer "
	if len(token) == 0 || !strings.HasPrefix(token, bearerPrefix) {
		return errors.New("invalid token")
	}
	tokenParts := strings.Split(strings.Split(token, bearerPrefix)[1], ".")
	if len(tokenParts) != 3 {
		return errors.New("invalid token")
	}
	
	jwtHeaderEncoded := tokenParts[0]
	jwtBodyEncoded := tokenParts[1]
	jwtSignatureEncoded := tokenParts[2]
	
	signatureFromToken, _ := base64.RawURLEncoding.DecodeString(jwtSignatureEncoded)
	signatureFromValidation := jwt.sign(jwtHeaderEncoded, jwtBodyEncoded)
	validSignature := hmac.Equal(signatureFromToken, signatureFromValidation)
	if !validSignature {
		return errors.New("invalid signature for token")
	}
	
	jwtBodyDecoded, _ := base64.RawURLEncoding.DecodeString(jwtBodyEncoded)
	var body *Body
	json.Unmarshal(jwtBodyDecoded, &body)
	validIssuer := body.Iss == jwt.Config.Issuer
	if !validIssuer {
		return errors.New("invalid iss for token")
	}
	validAudience := body.Aud == jwt.Config.Audience
	if !validAudience {
		return errors.New("invalid aud for token")
	}
	tokenNotBefore := time.Unix(body.Nbf, 0).UTC()
	tokenExpiration := time.Unix(body.Exp, 0).UTC()
	now := time.Now().UTC()
	if now.Before(tokenNotBefore) {
		return errors.New(fmt.Sprintf("token not allowed before %s", tokenNotBefore))
	}
	if now.After(tokenExpiration) {
		return errors.New(fmt.Sprintf("token expired at %s", tokenExpiration))
	}

	return nil
}

func (jwt *Provider) Generate(userName string) string {
	now := time.Now()
	body := &Body{
		User: userName,
		Role: "everyone",
		Jti: uuid.New().String(),
		Iss: jwt.Config.Issuer,
		Aud: jwt.Config.Audience,
		Exp: now.Add(time.Hour * 1).Unix(),
		Nbf: now.Unix(),
		Iat: now.Unix(),
	}
	return jwt.generate(body)
}

func (jwt *Provider) generate(body *Body) string {
	headerBytes, _ := json.Marshal(jwt.jwtHeader)
	headerBase64 := base64.RawURLEncoding.EncodeToString(headerBytes)

	bodyBytes, _ := json.Marshal(body)
	bodyBase64 := base64.RawURLEncoding.EncodeToString(bodyBytes)

	signatureBytes := jwt.sign(headerBase64, bodyBase64)
	signatureBase64 := base64.RawURLEncoding.EncodeToString(signatureBytes)
	return fmt.Sprintf("%s.%s.%s", headerBase64, bodyBase64, signatureBase64)
}

func (jwt *Provider) sign(headerBase64 string, bodyBase64 string) []byte {
	toSign := fmt.Sprintf("%s.%s", headerBase64, bodyBase64)
	hasher := hmac.New(sha256.New, []byte(jwt.Config.Secret))
	log.Println(jwt.Config.Secret)
	log.Println([]byte(jwt.Config.Secret))
	hasher.Write([]byte(toSign))
	return hasher.Sum(nil)
}
