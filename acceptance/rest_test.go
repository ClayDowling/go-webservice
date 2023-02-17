package acceptance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"example.com/webservice/data"
)

func Test_CreatedUser_CanBeQueried(t *testing.T) {

	jsonBody := []byte(`{"FirstName": "Leroy", "LastName": "Jenkins", "Email": "tpk@wow.org"}`)
	bodyReader := bytes.NewReader(jsonBody)

	resp, err := http.Post("http://localhost:8080/user/", "application/json", bodyReader)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected response code 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	var returnedUser data.User
	err = json.Unmarshal(body, &returnedUser)
	if err != nil {
		t.Error(err)
	}

	if returnedUser.Id == 0 {
		t.Error("Expected user with Id in response, but not found")
	}

	getResp, err := http.Get(fmt.Sprintf("http://localhost:8080/user/%d", returnedUser.Id))
	if err != nil {
		t.Fatal(err)
	}
	if getResp.StatusCode != 200 {
		t.Errorf("Expected to get user %d with response code 200, but got %d", returnedUser.Id, getResp.StatusCode)
	}
}
