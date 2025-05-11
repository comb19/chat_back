package test

import (
	"bytes"
	"chat_back/cmd/app"
	"chat_back/test/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type RequestGuild struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ResponseGuild struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func TestPostGuild(t *testing.T) {
	const userID = "user_2wtyaUmmGOlugFuP6h89SS4Pej3"

	token, err := utils.FetchClerkToken(userID)
	if err != nil {
		log.Fatal(err)
		return
	}

	router := app.SetupRouter()

	// リクエストボディが空
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/guilds", nil)
	req.Header.Set("Authorization", "Bearer "+*token)
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)

	// 正例
	fmt.Println("正例")
	w = httptest.NewRecorder()

	requestGuild := RequestGuild{
		Name:        "test guild",
		Description: "test description",
	}
	body, err := json.Marshal(requestGuild)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))

	req, _ = http.NewRequest("POST", "/guilds", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+*token)

	router.ServeHTTP(w, req)

	res := w.Result()
	body, _ = io.ReadAll(res.Body)

	var responseGuild RequestGuild
	if err := json.Unmarshal(body, &responseGuild); err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, 201, w.Code)
	assert.Equal(t, requestGuild.Name, responseGuild.Name)
	assert.Equal(t, requestGuild.Description, responseGuild.Description)
}
