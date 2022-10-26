package database

import (
	"fmt"
	"github.com/ablancas22/messenger-backend/models"
	"testing"
)

/*
	Id:         "000000000000000000000033",
	SenderID:   "000000000000000000000012",
	ReceiverID: "000000000000000000000013",
	Group:      false,
	Content:    "Message",
*/
func Test_TaskCreate(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string          // The name of the test
		want    *models.Message // What out instance we want our function to return.
		wantErr bool            // whether we want an error.
		task    *models.Message // The input of the test
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"success",
			&models.Message{Id: "000000000000000000000033"},
			false,
			&models.Message{
				Id:         "000000000000000000000033",
				SenderID:   "000000000000000000000012",
				ReceiverID: "000000000000000000000013",
				Group:      false,
				Content:    "Message",
			},
		},
		{
			"missing content",
			&models.Message{Id: "000000000000000000000022"},
			true,
			&models.Message{
				Id:         "000000000000000000000033",
				SenderID:   "000000000000000000000012",
				ReceiverID: "000000000000000000000013",
				Group:      false,
			},
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testService := initTestMessageService()
			//fmt.Println("\n\nPRE CREATE: ", tt.task)
			got, err := testService.MessageCreate(tt.task)
			//fmt.Println("\nPOST CREATE: ", got)
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskService.TaskCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if !got.CreatedAt.IsZero() && !got.LastModified.IsZero() {
					tt.want.CreatedAt = got.CreatedAt
					tt.want.LastModified = got.LastModified
				}
			}
			var failMsg string
			switch tt.name {
			case "success":
				if got.Id != tt.want.Id || got.CreatedAt.IsZero() { // Asserting whether we get the correct wanted value
					failMsg = fmt.Sprintf("TaskService.TaskCreate() = %v, want %v", got, tt.want)
				}
			}
			if failMsg != "" { // Asserting whether we get the correct wanted value
				t.Errorf(failMsg)
			}
		})
	}
}

func Test_TasksFind(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string // The name of the test
		want    int    // What out instance we want our function to return.
		wantErr bool   // whether we want an error.
		task    *models.Message
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"user tasks success",
			1,
			false,
			&models.Message{
				SenderID: "000000000000000000000012",
			},
		},
		{
			"group tasks success",
			2,
			false,
			&models.Message{
				ReceiverID: "000000000000000000000002",
				Group:      true,
			},
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testService := setupTestMessages()
			got, err := testService.MessagesFind(tt.task)
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskService.TasksFind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want { // Asserting whether we get the correct wanted value
				t.Errorf("TaskService.TasksFind() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func Test_MessageFind(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string          // The name of the test
		want    *models.Message // What out instance we want our function to return.
		wantErr bool            // whether we want an error.
		task    *models.Message
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"find by id",
			&models.Message{Id: "000000000000000000000033"},
			false,
			&models.Message{Id: "000000000000000000000033"},
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testService := setupTestMessages()
			got, err := testService.MessageFind(tt.task)
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskService.TaskFind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var failMsg string
			switch tt.name {
			case "success":
				if got.Id != tt.want.Id { // Asserting whether we get the correct wanted value
					failMsg = fmt.Sprintf("TaskService.TaskFind() = %v, want %v", got.Id, tt.want.Id)
				}
			}
			if failMsg != "" {
				t.Errorf(failMsg)
			}

		})
	}
}

func Test_MessageDelete(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string          // The name of the test
		want    *models.Message // What out instance we want our function to return.
		wantErr bool            // whether we want an error.
		task    *models.Message
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"success",
			&models.Message{Id: "000000000000000000000033"},
			false,
			&models.Message{Id: "000000000000000000000033"},
		},
		{
			"task not found",
			nil,
			true,
			&models.Message{Id: "000000000000000000000025"},
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testService := setupTestMessages()
			got, err := testService.MessageDelete(tt.task)
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskService.TaskDelete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var failMsg string
			switch tt.name {
			case "success":
				if got.Id != tt.want.Id { // Asserting whether we get the correct wanted value
					failMsg = fmt.Sprintf("TaskService.TaskDelete() = %v, want %v", got.Id, tt.want.Id)
				}
			case "task not found":
				if got != tt.want { // Asserting whether we get the correct wanted value
					failMsg = fmt.Sprintf("TaskService.TaskDelete() = %v, want %v", got, tt.want)
				}
			}
			if failMsg != "" {
				t.Errorf(failMsg)
			}
		})
	}
}
