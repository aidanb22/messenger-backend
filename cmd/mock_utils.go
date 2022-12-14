package cmd

import (
	"bytes"
	"encoding/json"
	"github.com/ablancas22/messenger-backend/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Execute test an http request
func executeRequest(ta App, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	ta.server.Router.ServeHTTP(rr, req)
	return rr
}

// Check response code returned from a test http request
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

// signIn
func signIn(ta App, email string, password string) *httptest.ResponseRecorder {
	payload := []byte(`{"email":"` + email + `","password":"` + password + `"}`)
	req, _ := http.NewRequest("POST", "/auth", bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")
	response := executeRequest(ta, req)
	return response
}

// CreateTestGroup creates a group doc for test setup
func createTestGroup(ta App, groupType int) *models.Group {
	group := models.Group{}
	if groupType == 1 {
		group.Id = "000000000000000000000002"
		group.Name = "test2"
		group.LastModified = time.Now().UTC()
		group.CreatedAt = time.Now().UTC()
	} else {
		group.Id = "000000000000000000000003"
		group.Name = "test3"
		group.LastModified = time.Now().UTC()
		group.CreatedAt = time.Now().UTC()
	}
	_, err := ta.server.GroupService.GroupDocInsert(&group)
	if err != nil {
		panic(err)
	}
	return &group
}

// createTestUser creates a user doc for test setup
func createTestUser(ta App, userType int) *models.User {
	user := models.User{}
	if userType == 1 {
		user.Id = "000000000000000000000012"
		user.Username = "test_user"
		user.Password = "abc123"
		user.Email = "test2@email.com"
		user.RootAdmin = false
		user.LastActive = time.Now().UTC()
		user.CreatedAt = time.Now().UTC()
	} else {
		user.Id = "000000000000000000000013"
		user.Username = "test_user2"
		user.Password = "abc123"
		user.Email = "test3@email.com.com"
		user.RootAdmin = false
		user.LastActive = time.Now().UTC()
		user.CreatedAt = time.Now().UTC()
	}
	_, err := ta.server.UserService.UserDocInsert(&user)
	if err != nil {
		panic(err)
	}
	return &user
}

// createTestMessage creates a message doc for test setup
func createTestMessage(ta App, messageType int) *models.Message {
	message := models.Message{}
	now := time.Now()
	if messageType == 1 {

		message.ReceiverID = "000000000000000000000012"
		message.SenderID = "000000000000000000000002"
		message.Id = "0000000000000000000000021"
		message.Group = false
		message.CreatedAt = time.Now()
		message.Content = "Content"
		message.UpdatedAt = now.UTC()
		message.CreatedAt = now.UTC()
	} else {
		message.Id = "000000000000000000000022"
		message.ReceiverID = "000000000000000000000012"
		message.SenderID = "000000000000000000000002"
		message.Id = "0000000000000000000000021"
		message.Group = false
		message.CreatedAt = time.Now()
		message.Content = "Content"
		message.UpdatedAt = now.UTC()
	}
	_, err := ta.server.MessageService.MessageDocInsert(&message)
	if err != nil {
		panic(err)
	}
	return &message
}

// getTestUserPayload
func getTestUserPayload(tCase string) []byte {
	switch tCase {
	case "CREATE":
		return []byte(`{"username":"test_user","password":"abc123","firstname":"test","lastname":"user","email":"test2@email.com"}`)
	case "UPDATE":
		return []byte(`{"username":"newUserName","password":"newUserPass","email":"new_test@email.com"}`)
	}
	return nil
}

// getTestPasswordPayload
func getTestPasswordPayload(tCase string) []byte {
	switch tCase {
	case "UPDATE_PASSWORD_ERROR":
		return []byte(`{"current_password":"789test122","new_password":"789test124"}`)
	case "UPDATE_PASSWORD_SUCCESS":
		return []byte(`{"current_password":"abc123","new_password":"789test124"}`)
	}
	return nil
}

// getTestGroupPayload
func getTestGroupPayload(tCase string) []byte {
	switch tCase {
	case "CREATE":
		return []byte(`{"name":"testingGroup"}`)
	case "UPDATE":
		return []byte(`{"name":"newTestingGroup"}`)
	}
	return nil
}

// getTestTaskPayload
func getTestTaskPayload(tCase string) []byte {
	var tMessage models.Message
	now := time.Now()
	switch tCase {
	case "CREATE":
		/*
			tTask.Name = "testTask"
			tTask.Completed = false
			tTask.Due = now.Add(time.Hour * 24).UTC()
			tTask.Description = "Updated Task to complete"
			tTask.UserId = "000000000000000000000012"
			tTask.GroupId = "000000000000000000000002"
			b, _ := json.Marshal(tTask)
		*/

		tMessage.Id = "000000000000000000000022"
		tMessage.ReceiverID = "000000000000000000000012"
		tMessage.SenderID = "000000000000000000000002"
		tMessage.Id = "0000000000000000000000021"
		tMessage.Group = false
		tMessage.CreatedAt = time.Now()
		tMessage.Content = "Content"
		tMessage.UpdatedAt = now.UTC()
		b, _ := json.Marshal(tMessage)
		return b
	case "UPDATE":

		tMessage.Id = "000000000000000000000022"
		tMessage.ReceiverID = "000000000000000000000013"
		tMessage.SenderID = "000000000000000000000003"
		tMessage.Id = "0000000000000000000000041"
		tMessage.Group = false
		tMessage.CreatedAt = time.Now()
		tMessage.Content = "Content"
		tMessage.UpdatedAt = now.UTC()
		b, _ := json.Marshal(tMessage)
		return b

	}
	return nil
}
