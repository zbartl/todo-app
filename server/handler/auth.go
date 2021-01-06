package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gobridge-kr/todo-app/server/utils"
	"io/ioutil"
	"log"
	"net/http"
)

type AuthHandler struct{
	jwt *jwtea.Provider
}

type AuthRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type AuthPayload struct {
	ClientId string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience string `json:"audience"`
	GrantType string `json:"grant_type"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	var body AuthRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	payload := &AuthPayload{
		ClientId: h.jwt.Config.ThirdPartyConfig.ClientId,
		ClientSecret: h.jwt.Config.ThirdPartyConfig.ClientSecret,
		Audience: h.jwt.Config.ThirdPartyConfig.ThirdPartyAudience,
		GrantType:	"password",
		UserName: 	body.UserName,
		Password: 	body.Password,
	}

	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", h.jwt.Config.ThirdPartyConfig.Url, bytes.NewBuffer(jsonPayload))
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	if res.StatusCode != 200 {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(bodyBytes))
		
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}
	
	token := h.jwt.Generate(body.UserName)
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	json.NewEncoder(w).Encode(token)
}

func Auth(jwt *jwtea.Provider) *AuthHandler {
	return &AuthHandler{
		jwt: jwt,
	}
}
