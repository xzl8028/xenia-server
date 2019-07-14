// Copyright (c) 2017-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package api4

import (
	"net/http"
	"strconv"
    // "github.com/xzl8028/xenia-server/mlog"
	"github.com/xzl8028/xenia-server/model"
)

func (api *API) InitTask() {
	// api.BaseRoutes.Tasks.Handle("", api.ApiSessionRequired(createTask)).Methods("POST")
	// api.BaseRoutes.Task.Handle("", api.ApiSessionRequired(patchTask)).Methods("PUT")
	api.BaseRoutes.Task.Handle("", api.ApiSessionRequired(getTask)).Methods("GET")
	api.BaseRoutes.Tasks.Handle("", api.ApiSessionRequired(getTasks)).Methods("GET")
	api.BaseRoutes.Task.Handle("/insert", api.ApiSessionRequired(insertTask)).Methods("POST")
	api.BaseRoutes.Task.Handle("/update", api.ApiSessionRequired(updateTask)).Methods("POST")
	api.BaseRoutes.Task.Handle("/update_status_quick", api.ApiSessionRequired(updateTaskStatusQuick)).Methods("POST")
	// api.BaseRoutes.Task.Handle("/disable", api.ApiSessionRequired(disableTask)).Methods("POST")
	// api.BaseRoutes.Task.Handle("/enable", api.ApiSessionRequired(enableTask)).Methods("POST")
	// api.BaseRoutes.Task.Handle("/assign/{user_id:[A-Za-z0-9]+}", api.ApiSessionRequired(assignTask)).Methods("POST")
}

// func createTask(c *Context, w http.ResponseWriter, r *http.Request) {
// 	taskPatch := model.TaskPatchFromJson(r.Body)
// 	if taskPatch == nil {
// 		c.SetInvalidParam("task")
// 		return
// 	}

// 	task := &model.Task{
// 		OwnerId: c.App.Session.UserId,
// 	}
// 	task.Patch(taskPatch)

// 	if !c.App.SessionHasPermissionTo(c.App.Session, model.PERMISSION_CREATE_TASK) {
// 		c.SetPermissionError(model.PERMISSION_CREATE_TASK)
// 		return
// 	}

// 	if user, err := c.App.GetUser(c.App.Session.UserId); err == nil {
// 		if user.IsTask {
// 			c.SetPermissionError(model.PERMISSION_CREATE_TASK)
// 			return
// 		}
// 	}

// 	if !*c.App.Config().ServiceSettings.EnableTaskAccountCreation {
// 		c.Err = model.NewAppError("createTask", "api.task.create_disabled", nil, "", http.StatusForbidden)
// 		return
// 	}

// 	createdTask, err := c.App.CreateTask(task)
// 	if err != nil {
// 		c.Err = err
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	w.Write(createdTask.ToJson())
// }

// func patchTask(c *Context, w http.ResponseWriter, r *http.Request) {
// 	c.RequireTaskId()
// 	if c.Err != nil {
// 		return
// 	}
// 	TaskId := c.Params.TaskId

// 	taskPatch := model.TaskPatchFromJson(r.Body)
// 	if taskPatch == nil {
// 		c.SetInvalidParam("task")
// 		return
// 	}

// 	if err := c.App.SessionHasPermissionToManageTask(c.App.Session, TaskId); err != nil {
// 		c.Err = err
// 		return
// 	}

// 	updatedTask, err := c.App.PatchTask(TaskId, taskPatch)
// 	if err != nil {
// 		c.Err = err
// 		return
// 	}

// 	w.Write(updatedTask.ToJson())
// }

func getTask(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireTaskId()
	if c.Err != nil {
		return
	}
	taskId := c.Params.TaskId

	// includeDeleted := r.URL.Query().Get("include_deleted") == "true"

	task, err := c.App.GetTask(taskId)
	if err != nil {
		c.Err = err
		return
	}

	// if c.App.SessionHasPermissionTo(c.App.Session, model.PERMISSION_READ_OTHERS_TASKS) {
	// 	// Allow access to any task.
	// } else if task.OwnerId == c.App.Session.UserId {
	// 	if !c.App.SessionHasPermissionTo(c.App.Session, model.PERMISSION_READ_TASKS) {
	// 		// Pretend like the task doesn't exist at all to avoid revealing that the
	// 		// user is a task. It's kind of silly in this case, sine we created the task,
	// 		// but we don't have read task permissions.
	// 		c.Err = model.MakeTaskNotFoundError(TaskId)
	// 		return
	// 	}
	// } else {
	// 	// Pretend like the task doesn't exist at all, to avoid revealing that the
	// 	// user is a task.
	// 	c.Err = model.MakeTaskNotFoundError(TaskId)
	// 	return
	// }

	// if c.HandleEtag(task.Etag(), "Get Task", w, r) {
	// 	return
	// }

	w.Write(task.ToJson())
}

func getTasks(c *Context, w http.ResponseWriter, r *http.Request) {
	// includeDeleted := r.URL.Query().Get("include_deleted") == "true"
	// onlyOrphaned := r.URL.Query().Get("only_orphaned") == "true"

	// var OwnerId string
	// if c.App.SessionHasPermissionTo(c.App.Session, model.PERMISSION_READ_OTHERS_TASKS) {
	// 	// Get tasks created by any user.
	// 	OwnerId = ""
	// } else if c.App.SessionHasPermissionTo(c.App.Session, model.PERMISSION_READ_TASKS) {
	// 	// Only get tasks created by this user.
	// 	OwnerId = c.App.Session.UserId
	// } else {
	// 	c.SetPermissionError(model.PERMISSION_READ_TASKS)
	// 	return
	// }

	tasks, err := c.App.GetTasks()
	if err != nil {
		c.Err = err
		return
	}

	// if c.HandleEtag(tasks.Etag(), "Get Tasks", w, r) {
	// 	return
	// }

	w.Write(tasks.ToJson())
}

// func disableTask(c *Context, w http.ResponseWriter, r *http.Request) {
// 	updateTaskActive(c, w, r, false)
// }

// func enableTask(c *Context, w http.ResponseWriter, r *http.Request) {
// 	updateTaskActive(c, w, r, true)
// }

func updateTask(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireTaskId()
	if c.Err != nil {
		return
	}
	taskId := c.Params.TaskId

	old_task, err := c.App.GetTask(taskId)
	if err != nil {
		c.Err = err
		return
	}

	new_task := old_task

	// if r.URL.Query().Get("task_id") != "" {
	// 	i, err1 := strconv.Atoi(r.URL.Query().Get("task_id"))
	// 	err1 = err1
	// 	new_task.TaskId = i
	// }

	if r.URL.Query().Get("create_at") != "" {
		i, err1 := strconv.ParseInt(r.URL.Query().Get("create_at"), 10, 64)
		err1 = err1
		new_task.CreateAt = i
	}

	if r.URL.Query().Get("due_at") != "" {
		i, err1 := strconv.ParseInt(r.URL.Query().Get("due_at"), 10, 64)
		err1 = err1
		new_task.DueAt = i
	}

	if r.URL.Query().Get("confirm_at") != "" {
		i, err1 := strconv.ParseInt(r.URL.Query().Get("confirm_at"), 10, 64)
		err1 = err1
		new_task.ConfirmAt = i
	}

	if r.URL.Query().Get("finish_at") != "" {
		i, err1 := strconv.ParseInt(r.URL.Query().Get("finish_at"), 10, 64)
		err1 = err1
		new_task.FinishAt = i
	}

	if r.URL.Query().Get("send_dept") != "" {
		new_task.SendDept = r.URL.Query().Get("send_dept")
	}

	if r.URL.Query().Get("receive_dept") != "" {
		new_task.ReceiveDept = r.URL.Query().Get("receive_dept")
	}

	if r.URL.Query().Get("room_id") != "" {
		i, err1 := strconv.Atoi(r.URL.Query().Get("room_id"))
		err1 = err1
		new_task.RoomId = i
	}

	if r.URL.Query().Get("task_type") != "" {
		new_task.TaskType = r.URL.Query().Get("task_type")
	}

	if r.URL.Query().Get("note") != "" {
		new_task.Note = r.URL.Query().Get("note")
	}

	if r.URL.Query().Get("status") != "" {
		i, err1 := strconv.Atoi(r.URL.Query().Get("status"))
		err1 = err1
		new_task.Status = i
	}


	// if err := c.App.SessionHasPermissionToManageTask(c.App.Session, TaskId); err != nil {
	// 	c.Err = err
	// 	return
	// }

	// task, err := c.App.UpdateTaskActive(TaskId, active)
	// if err != nil {
	// 	c.Err = err
	// 	return
	// }

	task, err := c.App.UpdateTask(new_task)
	if err != nil {
		c.Err = err
		return
	}

	w.Write(task.ToJson())
}

func updateTaskStatusQuick(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireTaskId()
	if c.Err != nil {
			return
	}
	taskId := c.Params.TaskId

	old_task, err := c.App.GetTask(taskId)
	if err != nil {
			c.Err = err
			return
	}

	new_task := old_task

	if r.URL.Query().Get("status") == "1" && old_task.Status == 0 {
			new_task.ConfirmAt = model.GetMillis()
			new_task.Status = 1
	} else if r.URL.Query().Get("status") == "2" && old_task.Status == 1 {
			new_task.FinishAt = model.GetMillis()
			new_task.Status = 2
	} else if r.URL.Query().Get("status") == "3" && old_task.Status != 3 && old_task.Status != 2 {
			new_task.Status = 3
	}

	task, err := c.App.UpdateTask(new_task)
	if err != nil {
			c.Err = err
			return
	}

	w.Write(task.ToJson())
}

func insertTask(c *Context, w http.ResponseWriter, r *http.Request) {
	// new_task := model.TaskFromJson(r.Body)
	var new_task model.Task

	new_task.CreateAt = model.GetMillis()

	if r.URL.Query().Get("due_at") != "" {
		i, err1 := strconv.ParseInt(r.URL.Query().Get("due_at"), 10, 64)
		err1 = err1
		new_task.DueAt = i
	} else {
		return
	}

	new_task.ConfirmAt = -1

	new_task.FinishAt = -1

	if r.URL.Query().Get("send_dept") != "" {
		new_task.SendDept = r.URL.Query().Get("send_dept")
	} else {
		return
	}

	if r.URL.Query().Get("receive_dept") != "" {
		new_task.ReceiveDept = r.URL.Query().Get("receive_dept")
	} else {
		return
	}

	if r.URL.Query().Get("room_id") != "" {
		i, err1 := strconv.Atoi(r.URL.Query().Get("room_id"))
		err1 = err1
		new_task.RoomId = i
	} else {
		return
	}

	if r.URL.Query().Get("task_type") != "" {
		new_task.TaskType = r.URL.Query().Get("task_type")
	} else {
		return
	}

	if r.URL.Query().Get("note") != "" {
		new_task.Note = r.URL.Query().Get("note")
	} else {
		return
	}

	if r.URL.Query().Get("status") != "" {
		i, err1 := strconv.Atoi(r.URL.Query().Get("status"))
		err1 = err1
		new_task.Status = i
	} else {
		return
	}

	// mlog.Debug("passed")

	// if err := c.App.SessionHasPermissionToManageTask(c.App.Session, taskUserId); err != nil {
	// 	c.Err = err
	// 	return
	// }

	// task, err := c.App.UpdateTaskActive(taskUserId, active)
	// if err != nil {
	// 	c.Err = err
	// 	return
	// }

	task, err := c.App.InsertTask(&new_task)
	if err != nil {
		c.Err = err
		return
	}

	w.Write(task.ToJson())
}


// func assignTask(c *Context, w http.ResponseWriter, r *http.Request) {
// 	c.RequireUserId()
// 	c.RequireTaskId()
// 	if c.Err != nil {
// 		return
// 	}
// 	TaskId := c.Params.TaskId
// 	userId := c.Params.UserId

// 	if err := c.App.SessionHasPermissionToManageTask(c.App.Session, TaskId); err != nil {
// 		c.Err = err
// 		return
// 	}

// 	if user, err := c.App.GetUser(userId); err == nil {
// 		if user.IsTask {
// 			c.SetPermissionError(model.PERMISSION_ASSIGN_TASK)
// 			return
// 		}
// 	}

// 	task, err := c.App.UpdateTaskOwner(TaskId, userId)
// 	if err != nil {
// 		c.Err = err
// 		return
// 	}

// 	w.Write(task.ToJson())
// }
