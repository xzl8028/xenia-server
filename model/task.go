// Copyright (c) 2016-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package model

import (
	"encoding/json"
	// "fmt"
	"io"
	"net/http"
	// "strings"
	// "unicode/utf8"
)

// const (
// 	TASK_DISPLAY_NAME_MAX_RUNES = USER_FIRST_NAME_MAX_RUNES
// 	TASK_DESCRIPTION_MAX_RUNES  = 1024
// 	TASK_CREATOR_ID_MAX_RUNES   = KEY_VALUE_PLUGIN_ID_MAX_RUNES // UserId or PluginId
// )

// Task is a special type of User meant for programmatic interactions.
// Note that the primary key of a task is the UserId, and matches the primary key of the
// corresponding user.
type Task struct {

	TaskId      int  		`json:"task_id"`
	CreateAt    int64	    `json:"create_at"`
	DueAt       int64	    `json:"due_at"`
	ConfirmAt   int64	    `json:"confirm_at"`
	FinishAt    int64	    `json:"finish_at"`
	SendDept    string 		`json:"send_dept"`
	ReceiveDept string 		`json:"receive_dept"`
	RoomId      int  		`json:"room_id"`
	TaskType    string 		`json:"task_type"`
	Note        string 		`json:"note"`
	Status      int  		`json:"status"`
	TeamId      string      `json:"teamid"`
	PostId      string      `json:"postid"`
}

// // TaskPatch is a description of what fields to update on an existing task.
// type TaskPatch struct {
// 	Username    *string `json:"username"`
// 	DisplayName *string `json:"display_name"`
// 	Description *string `json:"description"`
// }

// // TaskGetOptions acts as a filter on bulk task fetching queries.
// type TaskGetOptions struct {
// 	OwnerId        string
// 	IncludeDeleted bool
// 	OnlyOrphaned   bool
// 	Page           int
// 	PerPage        int
// }

// TaskList is a list of tasks.
type TaskList []*Task

// Trace describes the minimum information required to identify a task for the purpose of logging.
func (t *Task) Trace() map[string]interface{} {
	return map[string]interface{}{"task_id": t.TaskId}
}

// Clone returns a shallow copy of the task.
func (t *Task) Clone() *Task {
	copy := *t
	return &copy
}

// IsValid validates the task and returns an error if it isn't configured correctly.
func (t *Task) IsValid() *AppError {
	// if !IsValidId(t.TaskId) {
	// 	return NewAppError("Task.IsValid", "model.task.is_valid.user_id.app_error", t.Trace(), "", http.StatusBadRequest)
	// }

	// if !IsValidUsername(t.Username) {
	// 	return NewAppError("Task.IsValid", "model.task.is_valid.username.app_error", t.Trace(), "", http.StatusBadRequest)
	// }

	// if utf8.RuneCountInString(t.DisplayName) > TASK_DISPLAY_NAME_MAX_RUNES {
	// 	return NewAppError("Task.IsValid", "model.task.is_valid.user_id.app_error", t.Trace(), "", http.StatusBadRequest)
	// }

	// if utf8.RuneCountInString(t.Description) > TASK_DESCRIPTION_MAX_RUNES {
	// 	return NewAppError("Task.IsValid", "model.task.is_valid.description.app_error", t.Trace(), "", http.StatusBadRequest)
	// }

	// if len(t.OwnerId) == 0 || utf8.RuneCountInString(t.OwnerId) > TASK_CREATOR_ID_MAX_RUNES {
	// 	return NewAppError("Task.IsValid", "model.task.is_valid.creator_id.app_error", t.Trace(), "", http.StatusBadRequest)
	// }

	// if t.CreateAt == 0 {
	// 	return NewAppError("Task.IsValid", "model.task.is_valid.create_at.app_error", t.Trace(), "", http.StatusBadRequest)
	// }

	// if t.UpdateAt == 0 {
	// 	return NewAppError("Task.IsValid", "model.task.is_valid.update_at.app_error", t.Trace(), "", http.StatusBadRequest)
	// }

	return nil
}

// // PreSave should be run before saving a new task to the database.
// func (t *Task) PreSave() {
// 	t.CreateAt = GetMillis()
// 	t.UpdateAt = t.CreateAt
// 	t.DeleteAt = 0
// }

// PreUpdate should be run before saving an updated task to the database.
func (t *Task) PreUpdate() {
	// t.UpdateAt = GetMillis()
}

// // Etag generates an etag for caching.
// func (t *Task) Etag() string {
// 	return Etag(t.UserId, t.UpdateAt)
// }

// ToJson serializes the task to json.
func (t *Task) ToJson() []byte {
	data, _ := json.Marshal(t)
	return data
}

// taskFromJson deserializes a task from json.
func TaskFromJson(data io.Reader) *Task {
	var task *Task
	json.NewDecoder(data).Decode(&task)
	return task
}

// // Patch modifies an existing task with optional fields from the given patch.
// func (t *Task) Patch(patch *TaskPatch) {
// 	if patch.Username != nil {
// 		t.Username = *patch.Username
// 	}

// 	if patch.DisplayName != nil {
// 		t.DisplayName = *patch.DisplayName
// 	}

// 	if patch.Description != nil {
// 		t.Description = *patch.Description
// 	}
// }

// // ToJson serializes the task patch to json.
// func (t *TaskPatch) ToJson() []byte {
// 	data, err := json.Marshal(t)
// 	if err != nil {
// 		return nil
// 	}

// 	return data
// }

// // TaskPatchFromJson deserializes a task patch from json.
// func TaskPatchFromJson(data io.Reader) *TaskPatch {
// 	decoder := json.NewDecoder(data)
// 	var taskPatch TaskPatch
// 	err := decoder.Decode(&taskPatch)
// 	if err != nil {
// 		return nil
// 	}

// 	return &taskPatch
// }

// // UserFromTask returns a user model describing the task fields stored in the User store.
// func UserFromTask(t *Task) *User {
// 	return &User{
// 		Id:        t.UserId,
// 		Username:  t.Username,
// 		Email:     fmt.Sprintf("%s@localhost", strings.ToLower(t.Username)),
// 		FirstName: t.DisplayName,
// 		Roles:     SYSTEM_USER_ROLE_ID,
// 	}
// }

// // TaskFromUser returns a task model given a user model
// func TaskFromUser(u *User) *Task {
// 	return &Task{
// 		OwnerId:     u.Id,
// 		UserId:      u.Id,
// 		Username:    u.Username,
// 		DisplayName: u.GetDisplayName(SHOW_USERNAME),
// 	}
// }

// taskListFromJson deserializes a list of tasks from json.
func TaskListFromJson(data io.Reader) TaskList {
	var tasks TaskList
	json.NewDecoder(data).Decode(&tasks)
	return tasks
}

// ToJson serializes a list of tasks to json.
func (l *TaskList) ToJson() []byte {
	t, _ := json.Marshal(l)
	return t
}

// // Etag computes the etag for a list of tasks.
// func (l *TaskList) Etag() string {
// 	id := "0"
// 	var t int64 = 0
// 	var delta int64 = 0

// 	for _, v := range *l {
// 		if v.UpdateAt > t {
// 			t = v.UpdateAt
// 			id = v.UserId
// 		}

// 	}

// 	return Etag(id, t, delta, len(*l))
// }

// MakeTaskNotFoundError creates the error returned when a task does not exist, or when the user isn't allowed to query the task.
// The errors must the same in taskh cases to avoid leaking that a user is a task.
func MakeTaskNotFoundError(taskId string) *AppError {
	return NewAppError("SqlTaskStore.Get", "store.sql_task.get.missing.app_error", map[string]interface{}{"task_id": taskId}, "", http.StatusNotFound)
}
