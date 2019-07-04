// Copyright (c) 2016-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	// "github.com/xzl8028/xenia-server/mlog"
	"github.com/xzl8028/xenia-server/model"
	// "github.com/xzl8028/xenia-server/store"
	// "github.com/xzl8028/xenia-server/utils"
)

// CreateTask creates the given task and corresponding user.
// func (a *App) CreateTask(task *model.Task) (*model.Task, *model.AppError) {
// 	result := <-a.Srv.Store.User().Save(model.UserFromTask(task))
// 	if result.Err != nil {
// 		return nil, result.Err
// 	}
// 	task.UserId = result.Data.(*model.User).Id

// 	savedTask, err := a.Srv.Store.Task().Save(task)
// 	if err != nil {
// 		a.Srv.Store.User().PermanentDelete(task.UserId)
// 		return nil, err
// 	}

// 	// Get the owner of the task, if one exists. If not, don't send a message
// 	ownerUser, err := a.Srv.Store.User().Get(task.OwnerId)
// 	if err != nil && err.Id != store.MISSING_ACCOUNT_ERROR {
// 		mlog.Error(err.Error())
// 		return nil, err
// 	} else if ownerUser != nil {
// 		// Send a message to the task's creator to inform them that the task needs to be added
// 		// to a team and channel after it's created
// 		channel, err := a.GetOrCreateDirectChannel(savedTask.UserId, task.OwnerId)
// 		if err != nil {
// 			return nil, err
// 		}

// 		T := utils.GetUserTranslations(ownerUser.Locale)
// 		taskAddPost := &model.Post{
// 			Type:      model.POST_ADD_TASK_TEAMS_CHANNELS,
// 			UserId:    savedTask.UserId,
// 			ChannelId: channel.Id,
// 			Message:   T("api.task.teams_channels.add_message_mobile"),
// 		}

// 		if _, err := a.CreatePostAsUser(taskAddPost, a.Session.Id); err != nil {
// 			return nil, err
// 		}
// 	}

// 	return savedTask, nil
// }

// // PatchTask applies the given patch to the task and corresponding user.
// func (a *App) PatchTask(taskId string, taskPatch *model.TaskPatch) (*model.Task, *model.AppError) {
// 	task, err := a.GetTask(taskId, true)
// 	if err != nil {
// 		return nil, err
// 	}

// 	task.Patch(taskPatch)

// 	user, err := a.Srv.Store.User().Get(taskId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	patchedUser := model.UserFromTask(task)
// 	user.Id = patchedUser.Id
// 	user.Username = patchedUser.Username
// 	user.Email = patchedUser.Email
// 	user.FirstName = patchedUser.FirstName
// 	if _, err := a.Srv.Store.User().Update(user, true); err != nil {
// 		return nil, err
// 	}

// 	return a.Srv.Store.Task().Update(task)
// }

// GetTask returns the given task.
// func (a *App) GetTask(taskId string) (*model.Task, *model.AppError) {
// 	result := <-a.Srv.Store.Task().Get(taskId)
// 	if result.Err != nil {
// 		return nil, result.Err
// 	}

// 	return result.Data.(*model.Task), nil
// }

func (a *App) GetTask(taskId string) (*model.Task, *model.AppError) {
	return a.Srv.Store.Task().Get(taskId)
}

// GetTasks returns the requested page of tasks.
// func (a *App) GetTasks() (model.TaskList, *model.AppError) {
// 	result := <-a.Srv.Store.Task().GetAll()
// 	if result.Err != nil {
// 		return nil, result.Err
// 	}

// 	return result.Data.([]*model.Task), nil
// }

func (a *App) GetTasks() (model.TaskList, *model.AppError) {
	return a.Srv.Store.Task().GetAll()
}

// UpdateTaskActive marks a task as active or inactive, along with its corresponding user.
// func (a *App) UpdateTaskActive(taskId string, active bool) (*model.Task, *model.AppError) {
// 	user, err := a.Srv.Store.User().Get(taskId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if _, err = a.UpdateActive(user, active); err != nil {
// 		return nil, err
// 	}

// 	task, err := a.Srv.Store.Task().Get(taskId, true)
// 	if err != nil {
// 		return nil, err
// 	}

// 	changed := true
// 	if active && task.DeleteAt != 0 {
// 		task.DeleteAt = 0
// 	} else if !active && task.DeleteAt == 0 {
// 		task.DeleteAt = model.GetMillis()
// 	} else {
// 		changed = false
// 	}

// 	if changed {
// 		task, err = a.Srv.Store.Task().Update(task)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	return task, nil
// }

// PermanentDeleteTask permanently deletes a task and its corresponding user.
// func (a *App) PermanentDeleteTask(taskId string) *model.AppError {
// 	if err := a.Srv.Store.Task().PermanentDelete(taskId); err != nil {
// 		return err
// 	}

// 	if err := a.Srv.Store.User().PermanentDelete(taskId); err != nil {
// 		return err
// 	}

// 	return nil
// }

// UpdateTaskOwner changes a task's owner to the given value
// func (a *App) UpdateOwner(taskId, newOwnerId string) (*model.Task, *model.AppError) {
// 	task, err := a.Srv.Store.Task().Get(taskId, true)
// 	if err != nil {
// 		return nil, err
// 	}

// 	task.OwnerId = newOwnerId

// 	task, err = a.Srv.Store.Task().Update(task)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return task, nil
// }

// func (a *App) UpdateTask(task *model.Task) (*model.Task, *model.AppError) {
// 	if result := <-a.Srv.Store.Task().Update(task); result.Err != nil {
// 		return nil, result.Err
// 	}

// 	return task, nil
// }

// func (a *App) InsertTask(task *model.Task) (*model.Task, *model.AppError) {
// 	if result := <-a.Srv.Store.Task().Insert(task); result.Err != nil {
// 		return nil, result.Err
// 	}

// 	return task, nil
// }

// // disableUserTasks disables all tasks owned by the given user
// func (a *App) disableUserTasks(userId string) *model.AppError {
// 	perPage := 20
// 	for {
// 		options := &model.TaskGetOptions{
// 			OwnerId:        userId,
// 			IncludeDeleted: false,
// 			OnlyOrphaned:   false,
// 			Page:           0,
// 			PerPage:        perPage,
// 		}
// 		userTasks, err := a.GetTasks(options)
// 		if err != nil {
// 			return err
// 		}

// 		for _, task := range userTasks {
// 			_, err := a.UpdateTaskActive(task.UserId, false)
// 			if err != nil {
// 				mlog.Error("Unable to deactivate task.", mlog.String("task_user_id", task.UserId), mlog.Err(err))
// 			}
// 		}

// 		// Get next set of tasks if we got the max number of tasks
// 		if len(userTasks) == perPage {
// 			options.Page += 1
// 			continue
// 		}
// 		break
// 	}

// 	return nil
// }

// // ConvertUserToTask converts a user to task
// func (a *App) ConvertUserToTask(user *model.User) (*model.Task, *model.AppError) {
// 	return a.Srv.Store.Task().Save(model.TaskFromUser(user))
// }
