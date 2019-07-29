// Copyright (c) 2015-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package sqlstore

import (
	"database/sql"
	"net/http"
	// "strings"
	// "time"
	"strconv"

	"github.com/xzl8028/xenia-server/einterfaces"
	"github.com/xzl8028/xenia-server/model"
	"github.com/xzl8028/xenia-server/store"
)

// task is a subset of the model.Task type, omitting the model.User fields.
type task struct {
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

func taskFromModel(t *model.Task) *task {
	return &task{
		TaskId:      t.TaskId,
		CreateAt:    t.CreateAt,
		DueAt:       t.DueAt,
		ConfirmAt:   t.ConfirmAt,
		FinishAt:    t.FinishAt,
		SendDept:    t.SendDept,
		ReceiveDept: t.ReceiveDept,
		RoomId:      t.RoomId,
		TaskType:    t.TaskType,
		Note:        t.Note,
		Status:      t.Status,
		TeamId:      t.TeamId,
		PostId:      t.PostId,

	}
}

// SqlTaskStore is a store for managing tasks in the database.
// Tasks are otherwise normal users with extra metadata record in the Tasks table. The primary key
// for a task matches the primary key value for corresponding User record.
type SqlTaskStore struct {
	SqlStore
	metrics einterfaces.MetricsInterface
}

// NewSqlTaskStore creates an instance of SqlTaskStore, registering the table schema in question.
func NewSqlTaskStore(sqlStore SqlStore, metrics einterfaces.MetricsInterface) store.TaskStore {
	us := &SqlTaskStore{
		SqlStore: sqlStore,
		metrics:  metrics,
	}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(task{}, "Tasks").SetKeys(true, "TaskId")
		table.ColMap("TaskId").SetMaxSize(200)
		table.ColMap("CreateAt").SetMaxSize(200)
		table.ColMap("DueAt").SetMaxSize(200)
		table.ColMap("ConfirmAt").SetMaxSize(200)
		table.ColMap("FinishAt").SetMaxSize(200)
		table.ColMap("SendDept").SetMaxSize(200)
		table.ColMap("ReceiveDept").SetMaxSize(200)
		table.ColMap("RoomId").SetMaxSize(200)
		table.ColMap("TaskType").SetMaxSize(200)
		table.ColMap("Note").SetMaxSize(200)
		table.ColMap("Status").SetMaxSize(200)
		table.ColMap("TeamId").SetMaxSize(200)
		table.ColMap("PostId").SetMaxSize(200)
	}

	return us
}

// func (us SqlTaskStore) CreateIndexesIfNotExists() {
// }

// traceTask is a helper function for adding to a task trace when logging.
func traceTask(task *model.Task, extra map[string]interface{}) map[string]interface{} {
	trace := make(map[string]interface{})
	for key, value := range task.Trace() {
		trace[key] = value
	}
	for key, value := range extra {
		trace[key] = value
	}

	return trace
}

// Get fetches the given task in the database.
func (us SqlTaskStore) Get(taskId string) (*model.Task, *model.AppError) {
	// var excludeDeletedSql = "AND b.DeleteAt = 0"
	// if includeDeleted {
	// 	excludeDeletedSql = ""
	// }

	query :=  `
		SELECT * FROM Tasks
		WHERE TaskId =`+taskId

	var task *model.Task
	if err := us.GetReplica().SelectOne(&task, query); err == sql.ErrNoRows {
		return nil, model.MakeTaskNotFoundError(taskId)
	} else if err != nil {
		return nil, model.NewAppError("SqlTaskStore.Get", "store.sql_task.get.app_error", map[string]interface{}{"task_id": taskId}, err.Error(), http.StatusInternalServerError)
	}

	// var task *model.Task
	// if err := us.GetReplica().SelectOne(&task, sql); /*err == sql.ErrNoRows {
	// 		result.Err = model.MakeTaskNotFoundError(taskId)
	// } else if*/ err != nil {
	// 		result.Err = model.NewAppError("SqlTaskStore.GetTask", "store.sql_task.get.app_error", $
	// } //else {
	// 		result.Data = task
	// // }

	return task, nil
}


// GetAll fetches from all tasks in the database.
func (us SqlTaskStore) GetAll() ([]*model.Task, *model.AppError) {
	// params := map[string]interface{}{
	// 	"offset": options.Page * options.PerPage,
	// 	"limit":  options.PerPage,
	// }

	// var conditions []string
	// var conditionsSql string
	// var additionalJoin string

	// if !options.IncludeDeleted {
	// 	conditions = append(conditions, "b.DeleteAt = 0")
	// }
	// if options.OwnerId != "" {
	// 	conditions = append(conditions, "b.OwnerId = :creator_id")
	// 	params["creator_id"] = options.OwnerId
	// }
	// if options.OnlyOrphaned {
	// 	additionalJoin = "JOIN Users o ON (o.Id = b.OwnerId)"
	// 	conditions = append(conditions, "o.DeleteAt != 0")
	// }

	// if len(conditions) > 0 {
	// 	conditionsSql = "WHERE " + strings.Join(conditions, " AND ")
	// }

	query := `
		SELECT * FROM Tasks
	`

	var tasks []*model.Task
	if _, err := us.GetReplica().Select(&tasks, query); err != nil {
		return nil, model.NewAppError("SqlTaskStore.GetAll", "store.sql_task.get_all.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return tasks, nil
}

// // Save persists a new task to the database.
// // It assumes the corresponding user was saved via the user store.
// func (us SqlTaskStore) Save(task *model.Task) (*model.Task, *model.AppError) {
// 	task = task.Clone()
// 	task.PreSave()

// 	if err := task.IsValid(); err != nil {
// 		return nil, err
// 	}

// 	if err := us.GetMaster().Insert(taskFromModel(task)); err != nil {
// 		return nil, model.NewAppError("SqlTaskStore.Save", "store.sql_task.save.app_error", task.Trace(), err.Error(), http.StatusInternalServerError)
// 	}

// 	return task, nil
// }

// Update persists an updated task to the database.
// It assumes the corresponding user was updated via the user store.
func (us SqlTaskStore) Update(task *model.Task) (*model.Task, *model.AppError) {
	// task = task.Clone()

	// task.PreUpdate()
	// if err := task.IsValid(); err != nil {
	// 	return nil, err
	// }

	// oldTask, err := us.Get(task.TaskId, true)
	// if err != nil {
	// 	return nil, err
	// }

	// oldTask.CreateAt = task.CreateAt
	// oldTask.DueAt = task.DueAt
	// oldTask.ConfirmAt = task.ConfirmAt
	// oldTask.FinishAt = task.FinishAt
	// oldTask.SendDept = task.SendDept
	// oldTask.ReceiveDept = task.ReceiveDept
	// oldTask.RoomId = task.RoomId
	// oldTask.TaskType = task.TaskType
	// oldTask.Note = task.Note
	// oldTask.Status = task.Status
	// task = oldTask

	// if count, err := us.GetMaster().Update(taskFromModel(task)); err != nil {
	// 	return nil, model.NewAppError("SqlTaskStore.Update", "store.sql_task.update.updating.app_error", task.Trace(), err.Error(), http.StatusInternalServerError)
	// } else if count != 1 {
	// 	return nil, model.NewAppError("SqlTaskStore.Update", "store.sql_task.update.app_error", traceTask(task, map[string]interface{}{"count": count}), "", http.StatusInternalServerError)
	// }

	if _, err := us.GetMaster().Exec("UPDATE Tasks SET CreateAt = :CreateAt, DueAt = :DueAt, ConfirmAt = :ConfirmAt, FinishAt = :FinishAt, SendDept = :SendDept, ReceiveDept = :ReceiveDept, RoomId = :RoomId, TaskType = :TaskType, Note = :Note, Status = :Status, TeamId = :TeamId, PostId = :PostId WHERE TaskId = :TaskId", map[string]interface{}{"TaskId": task.TaskId, "CreateAt": task.CreateAt, "DueAt": task.DueAt, "ConfirmAt": task.ConfirmAt, "FinishAt": task.FinishAt, "SendDept": task.SendDept, "ReceiveDept": task.ReceiveDept, "RoomId": task.RoomId, "TaskType": task.TaskType, "Note": task.Note, "Status": task.Status, "TeamId": task.TeamId, "PostId": task.PostId}); err != nil {
		return nil, model.NewAppError("SqlTaskStore.Update", "store.sql_task.update.app_error", nil, "task_id="+strconv.Itoa(task.TaskId), http.StatusInternalServerError)
	}

	return task, nil
}

func (us SqlTaskStore) Insert(task *model.Task) (*model.Task, *model.AppError) {

	if res, err := us.GetMaster().Exec("INSERT INTO Tasks (CreateAt, DueAt, ConfirmAt, FinishAt, SendDept, ReceiveDept, RoomId, TaskType, Note, Status,TeamId,PostId) VALUES (:CreateAt, :DueAt, :ConfirmAt, :FinishAt, :SendDept, :ReceiveDept, :RoomId, :TaskType, :Note, :Status, :TeamId, :PostId);", map[string]interface{}{"CreateAt": task.CreateAt, "DueAt": task.DueAt, "ConfirmAt": task.ConfirmAt, "FinishAt": task.FinishAt, "SendDept": task.SendDept, "ReceiveDept": task.ReceiveDept, "RoomId": task.RoomId, "TaskType": task.TaskType, "Note": task.Note, "Status": task.Status,"TeamId": task.TeamId, "PostId": task.PostId}); err != nil {
		return nil, model.NewAppError("SqlTaskStore.Insert", "store.sql_task.insert.app_error", nil, "", http.StatusInternalServerError)
	} else {
		id, err := res.LastInsertId()
		if err != nil {
			println("Error!")
		} else {
			println("LastInsertId:", id)
			task.TaskId = int(id)
		}
	}

	return task, nil
}



// // PermanentDelete removes the task from the database altogether.
// // If the corresponding user is to be deleted, it must be done via the user store.
// func (us SqlTaskStore) PermanentDelete(taskUserId string) *model.AppError {
// 	query := "DELETE FROM Tasks WHERE UserId = :user_id"
// 	if _, err := us.GetMaster().Exec(query, map[string]interface{}{"user_id": taskUserId}); err != nil {
// 		return model.NewAppError("SqlTaskStore.Update", "store.sql_task.delete.app_error", map[string]interface{}{"user_id": taskUserId}, err.Error(), http.StatusBadRequest)
// 	}
// 	return nil
// }
