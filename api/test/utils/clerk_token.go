package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/caarlos0/env/v11"
)

const clerkURL = "https://api.clerk.com/v1"

type EnvVar struct {
	ClerkSecKey string `env:"CLERK_SECRET_KEY,notEmpty"`
}

type sessionBody struct {
	UserID string `json:"user_id"`
}
type sessionResponse struct {
	Id string `json:"id"`
}
type tokenResponse struct {
	JWT string `json:"jwt"`
}

func FetchClerkToken(userID string) (*string, error) {
	var envVar EnvVar
	if err := env.Parse(&envVar); err != nil {
		panic(err)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	reqBody, err := json.Marshal(sessionBody{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", clerkURL+"/sessions", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+envVar.ClerkSecKey)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var sessionRes sessionResponse
	if err = json.Unmarshal(resBody, &sessionRes); err != nil {
		return nil, err
	}
	res.Body.Close()

	req, err = http.NewRequest("POST", clerkURL+"/sessions/"+sessionRes.Id+"/tokens", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+envVar.ClerkSecKey)
	res, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	resBody, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var tokenRes tokenResponse
	if err = json.Unmarshal(resBody, &tokenRes); err != nil {
		return nil, err
	}
	res.Body.Close()
	return &tokenRes.JWT, nil
}
