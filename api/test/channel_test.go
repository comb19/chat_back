package test

import (
	"bytes"
	"chat_back/cmd/app"
	"chat_back/test/utils"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	userID          = "user_2wtyaUmmGOlugFuP6h89SS4Pej3"
	channelID       = "590b06b0-d5b9-4ab8-b878-f90f24883b4c"
	joinedChannelID = "0a96949f-cf72-43e5-aba0-cfd085ed016c"
)

type RequestChannel struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
}

type RequestUsers struct {
	UserIDs []string `json:"user_ids"`
}

// func TestPostChannel(t *testing.T) {
// 	token, err := utils.FetchClerkToken(userID)
// 	if err != nil {
// 		log.Fatal(err)
// 		return
// 	}

// 	router := app.SetupRouter()
// 	w := httptest.NewRecorder()

// 	req, _ := http.NewRequest("POST", "/channels", nil)
// 	req.Header.Set("Authorization", "Bearer "+*token)
// 	router.ServeHTTP(w, req)
// 	assert.Equal(t, 400, w.Code)

// 	w = httptest.NewRecorder()
// 	requestChannel := RequestChannel{
// 		Name:        "test channel",
// 		Description: "test description",
// 		Private:     false,
// 	}
// 	body, err := json.Marshal(requestChannel)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	req, _ = http.NewRequest("POST", "/channels", bytes.NewReader(body))
// 	req.Header.Set("Authorization", "Bearer "+*token)
// 	router.ServeHTTP(w, req)
// 	assert.Equal(t, 201, w.Code)
// }

func TestPostChannelUsers(t *testing.T) {
	token, err := utils.FetchClerkToken(userID)
	if err != nil {
		log.Fatal(err)
		return
	}

	router := app.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/channels/"+channelID+"/users", nil)
	req.Header.Set("Authorization", "Bearer "+*token)
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	requestUsers := RequestUsers{
		UserIDs: []string{userID},
	}
	body, err := json.Marshal(requestUsers)
	if err != nil {
		log.Fatal(err)
	}
	req, _ = http.NewRequest("POST", "/channels/"+channelID+"/users", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+*token)
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
}
