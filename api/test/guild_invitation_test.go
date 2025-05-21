package test

import (
	"bytes"
	"chat_back/cmd/app"
	"chat_back/interface/types"
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

func TestCreateGuildInvitation(t *testing.T) {
	const ownerID = "user_2wtyaUmmGOlugFuP6h89SS4Pej3"
	const inviteeID = "user_2xBp8XPtjCCnc7RLAICKTpzRrfB"
	const guildID = "6b3075b7-d72b-4c41-983f-71213b16e1d7"

	ownerToken, err := utils.FetchClerkToken(ownerID)
	if err != nil {
		log.Fatal(err)
		return
	}
	inviteeToken, err := utils.FetchClerkToken(inviteeID)
	if err != nil {
		log.Fatal(err)
		return
	}

	router := app.SetupRouter()

	requestGuildInvitation := types.RequestGuildInvitation{
		GuildID: guildID,
	}
	reqBody, err := json.Marshal(requestGuildInvitation)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/invitations/guilds/", bytes.NewReader(reqBody))
	req.Header.Set("Authorization", "Bearer "+*ownerToken)
	router.ServeHTTP(w, req)

	res := w.Result()
	body, _ := io.ReadAll(res.Body)
	var responseGuildInvitation types.ResponseGuildInvitation
	if err := json.Unmarshal(body, &responseGuildInvitation); err != nil {
		log.Fatal(err)
		return
	}

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, guildID, responseGuildInvitation.GuildID)
	assert.Equal(t, ownerID, responseGuildInvitation.OwnerID)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", fmt.Sprintf("/invitations/guilds/%s", responseGuildInvitation.ID), nil)
	req.Header.Set("Authorization", "Bearer "+*inviteeToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
}
