// Copyright (c) 2016-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/xzl8028/xenia-server/model"
)

// func TestCreateTask(t *testing.T) {
// 	t.Run("invalid task", func(t *testing.T) {
// 		t.Run("relative to user", func(t *testing.T) {
// 			th := Setup(t).InitBasic()
// 			defer th.TearDown()

// 			_, err := th.App.CreateTask(&model.Task{
// 				Username:    "invalid username",
// 				Description: "a task",
// 				OwnerId:     th.BasicUser.Id,
// 			})
// 			require.NotNil(t, err)
// 			require.Equal(t, "model.user.is_valid.username.app_error", err.Id)
// 		})

// 		t.Run("relative to task", func(t *testing.T) {
// 			th := Setup(t).InitBasic()
// 			defer th.TearDown()

// 			_, err := th.App.CreateTask(&model.Task{
// 				Username:    "username",
// 				Description: strings.Repeat("x", 1025),
// 				OwnerId:     th.BasicUser.Id,
// 			})
// 			require.NotNil(t, err)
// 			require.Equal(t, "model.task.is_valid.description.app_error", err.Id)
// 		})
// 	})

// 	t.Run("create task", func(t *testing.T) {
// 		th := Setup(t).InitBasic()
// 		defer th.TearDown()

// 		task, err := th.App.CreateTask(&model.Task{
// 			Username:    "username",
// 			Description: "a task",
// 			OwnerId:     th.BasicUser.Id,
// 		})
// 		require.Nil(t, err)
// 		defer th.App.PermanentDeleteTask(task.UserId)
// 		assert.Equal(t, "username", task.Username)
// 		assert.Equal(t, "a task", task.Description)
// 		assert.Equal(t, th.BasicUser.Id, task.OwnerId)

// 		// Check that a post was created to add task to team and channels
// 		channel, err := th.App.GetOrCreateDirectChannel(task.UserId, th.BasicUser.Id)
// 		require.Nil(t, err)
// 		posts, err := th.App.GetPosts(channel.Id, 0, 1)
// 		require.Nil(t, err)

// 		postArray := posts.ToSlice()
// 		assert.Len(t, postArray, 1)
// 		assert.Equal(t, postArray[0].Type, model.POST_ADD_TASK_TEAMS_CHANNELS)
// 	})

// 	t.Run("create task, username already used by a non-task user", func(t *testing.T) {
// 		th := Setup(t).InitBasic()
// 		defer th.TearDown()

// 		_, err := th.App.CreateTask(&model.Task{
// 			Username:    th.BasicUser.Username,
// 			Description: "a task",
// 			OwnerId:     th.BasicUser.Id,
// 		})
// 		require.NotNil(t, err)
// 		require.Equal(t, "store.sql_user.save.username_exists.app_error", err.Id)
// 	})
// }

// func TestPatchTask(t *testing.T) {
// 	t.Run("invalid patch for user", func(t *testing.T) {
// 		th := Setup(t).InitBasic()
// 		defer th.TearDown()

// 		task, err := th.App.CreateTask(&model.Task{
// 			Username:    "username",
// 			Description: "a task",
// 			OwnerId:     th.BasicUser.Id,
// 		})
// 		require.Nil(t, err)
// 		defer th.App.PermanentDeleteTask(task.UserId)

// 		taskPatch := &model.TaskPatch{
// 			Username:    sToP("invalid username"),
// 			DisplayName: sToP("an updated task"),
// 			Description: sToP("updated task"),
// 		}

// 		_, err = th.App.PatchTask(task.UserId, taskPatch)
// 		require.NotNil(t, err)
// 		require.Equal(t, "model.user.is_valid.username.app_error", err.Id)
// 	})

// 	t.Run("invalid patch for task", func(t *testing.T) {
// 		th := Setup(t).InitBasic()
// 		defer th.TearDown()

// 		task, err := th.App.CreateTask(&model.Task{
// 			Username:    "username",
// 			Description: "a task",
// 			OwnerId:     th.BasicUser.Id,
// 		})
// 		require.Nil(t, err)
// 		defer th.App.PermanentDeleteTask(task.UserId)

// 		taskPatch := &model.TaskPatch{
// 			Username:    sToP("username"),
// 			DisplayName: sToP("display name"),
// 			Description: sToP(strings.Repeat("x", 1025)),
// 		}

// 		_, err = th.App.PatchTask(task.UserId, taskPatch)
// 		require.NotNil(t, err)
// 		require.Equal(t, "model.task.is_valid.description.app_error", err.Id)
// 	})

// 	t.Run("patch task", func(t *testing.T) {
// 		th := Setup(t).InitBasic()
// 		defer th.TearDown()

// 		task := &model.Task{
// 			Username:    "username",
// 			DisplayName: "task",
// 			Description: "a task",
// 			OwnerId:     th.BasicUser.Id,
// 		}

// 		createdTask, err := th.App.CreateTask(task)
// 		require.Nil(t, err)
// 		defer th.App.PermanentDeleteTask(createdTask.UserId)

// 		taskPatch := &model.TaskPatch{
// 			Username:    sToP("username2"),
// 			DisplayName: sToP("updated task"),
// 			Description: sToP("an updated task"),
// 		}

// 		patchedTask, err := th.App.PatchTask(createdTask.UserId, taskPatch)
// 		require.Nil(t, err)

// 		createdTask.Username = "username2"
// 		createdTask.DisplayName = "updated task"
// 		createdTask.Description = "an updated task"
// 		createdTask.UpdateAt = patchedTask.UpdateAt
// 		require.Equal(t, createdTask, patchedTask)
// 	})

// 	t.Run("patch task, username already used by a non-task user", func(t *testing.T) {
// 		th := Setup(t).InitBasic()
// 		defer th.TearDown()

// 		task, err := th.App.CreateTask(&model.Task{
// 			Username:    "username",
// 			DisplayName: "task",
// 			Description: "a task",
// 			OwnerId:     th.BasicUser.Id,
// 		})
// 		require.Nil(t, err)
// 		defer th.App.PermanentDeleteTask(task.UserId)

// 		taskPatch := &model.TaskPatch{
// 			Username: sToP(th.BasicUser2.Username),
// 		}

// 		_, err = th.App.PatchTask(task.UserId, taskPatch)
// 		require.NotNil(t, err)
// 		require.Equal(t, "store.sql_user.update.username_taken.app_error", err.Id)
// 	})
// }

// func TestGetTask(t *testing.T) {
// 	th := Setup(t).InitBasic()
// 	defer th.TearDown()

// 	task1, err := th.App.CreateTask(&model.Task{
// 		Username:    "username",
// 		Description: "a task",
// 		OwnerId:     th.BasicUser.Id,
// 	})
// 	require.Nil(t, err)
// 	defer th.App.PermanentDeleteTask(task1.UserId)

// 	task2, err := th.App.CreateTask(&model.Task{
// 		Username:    "username2",
// 		Description: "a second task",
// 		OwnerId:     th.BasicUser.Id,
// 	})
// 	require.Nil(t, err)
// 	defer th.App.PermanentDeleteTask(task2.UserId)

// 	deletedTask, err := th.App.CreateTask(&model.Task{
// 		Username:    "username3",
// 		Description: "a deleted task",
// 		OwnerId:     th.BasicUser.Id,
// 	})
// 	require.Nil(t, err)
// 	deletedTask, err = th.App.UpdateTaskActive(deletedTask.UserId, false)
// 	require.Nil(t, err)
// 	defer th.App.PermanentDeleteTask(deletedTask.UserId)

// 	t.Run("get unknown task", func(t *testing.T) {
// 		_, err := th.App.GetTask(model.NewId(), false)
// 		require.NotNil(t, err)
// 		require.Equal(t, "store.sql_task.get.missing.app_error", err.Id)
// 	})

// 	t.Run("get task1", func(t *testing.T) {
// 		task, err := th.App.GetTask(task1.UserId, false)
// 		require.Nil(t, err)
// 		assert.Equal(t, task1, task)
// 	})

// 	t.Run("get task2", func(t *testing.T) {
// 		task, err := th.App.GetTask(task2.UserId, false)
// 		require.Nil(t, err)
// 		assert.Equal(t, task2, task)
// 	})

// 	t.Run("get deleted task", func(t *testing.T) {
// 		_, err := th.App.GetTask(deletedTask.UserId, false)
// 		require.NotNil(t, err)
// 		require.Equal(t, "store.sql_task.get.missing.app_error", err.Id)
// 	})

// 	t.Run("get deleted task, include deleted", func(t *testing.T) {
// 		task, err := th.App.GetTask(deletedTask.UserId, true)
// 		require.Nil(t, err)
// 		assert.Equal(t, deletedTask, task)
// 	})
// }

func TestGetTasks(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()

	// OwnerId1 := model.NewId()
	// OwnerId2 := model.NewId()

	// task1, err := th.App.CreateTask(&model.Task{
	// 	Username:    "username",
	// 	Description: "a task",
	// 	OwnerId:     OwnerId1,
	// })
	// require.Nil(t, err)
	// defer th.App.PermanentDeleteTask(task1.UserId)

	// deletedTask1, err := th.App.CreateTask(&model.Task{
	// 	Username:    "username4",
	// 	Description: "a deleted task",
	// 	OwnerId:     OwnerId1,
	// })
	// require.Nil(t, err)
	// deletedTask1, err = th.App.UpdateTaskActive(deletedTask1.UserId, false)
	// require.Nil(t, err)
	// defer th.App.PermanentDeleteTask(deletedTask1.UserId)

	// task2, err := th.App.CreateTask(&model.Task{
	// 	Username:    "username2",
	// 	Description: "a second task",
	// 	OwnerId:     OwnerId1,
	// })
	// require.Nil(t, err)
	// defer th.App.PermanentDeleteTask(task2.UserId)

	// task3, err := th.App.CreateTask(&model.Task{
	// 	Username:    "username3",
	// 	Description: "a third task",
	// 	OwnerId:     OwnerId1,
	// })
	// require.Nil(t, err)
	// defer th.App.PermanentDeleteTask(task3.UserId)

	// task4, err := th.App.CreateTask(&model.Task{
	// 	Username:    "username5",
	// 	Description: "a fourth task",
	// 	OwnerId:     OwnerId2,
	// })
	// require.Nil(t, err)
	// defer th.App.PermanentDeleteTask(task4.UserId)

	// deletedTask2, err := th.App.CreateTask(&model.Task{
	// 	Username:    "username6",
	// 	Description: "a deleted task",
	// 	OwnerId:     OwnerId2,
	// })
	// require.Nil(t, err)
	// deletedTask2, err = th.App.UpdateTaskActive(deletedTask2.UserId, false)
	// require.Nil(t, err)
	// defer th.App.PermanentDeleteTask(deletedTask2.UserId)

	t.Run("get tasks", func(t *testing.T) {
		fmt.Println("hahaha")
		tasks, err := th.App.GetTasks()
		fmt.Println(tasks)
		fmt.Println(err)
		require.Nil(t, err)

		// assert.Equal(t, model.TaskList{task1, task2, task3, task4}, tasks)
	})

	// t.Run("get tasks, page=0, perPage=10", func(t *testing.T) {
	// 	tasks, err := th.App.GetTasks(&model.TaskGetOptions{
	// 		Page:           0,
	// 		PerPage:        10,
	// 		OwnerId:        "",
	// 		IncludeDeleted: false,
	// 	})
	// 	require.Nil(t, err)
	// 	assert.Equal(t, model.TaskList{task1, task2, task3, task4}, tasks)
	// })

	// t.Run("get tasks, page=0, perPage=1", func(t *testing.T) {
	// 	tasks, err := th.App.GetTasks(&model.TaskGetOptions{
	// 		Page:           0,
	// 		PerPage:        1,
	// 		OwnerId:        "",
	// 		IncludeDeleted: false,
	// 	})
	// 	require.Nil(t, err)
	// 	assert.Equal(t, model.TaskList{task1}, tasks)
	// })

	// t.Run("get tasks, page=1, perPage=2", func(t *testing.T) {
	// 	tasks, err := th.App.GetTasks(&model.TaskGetOptions{
	// 		Page:           1,
	// 		PerPage:        2,
	// 		OwnerId:        "",
	// 		IncludeDeleted: false,
	// 	})
	// 	require.Nil(t, err)
	// 	assert.Equal(t, model.TaskList{task3, task4}, tasks)
	// })

	// t.Run("get tasks, page=2, perPage=2", func(t *testing.T) {
	// 	tasks, err := th.App.GetTasks(&model.TaskGetOptions{
	// 		Page:           2,
	// 		PerPage:        2,
	// 		OwnerId:        "",
	// 		IncludeDeleted: false,
	// 	})
	// 	require.Nil(t, err)
	// 	assert.Equal(t, model.TaskList{}, tasks)
	// })

	// t.Run("get tasks, page=0, perPage=10, include deleted", func(t *testing.T) {
	// 	tasks, err := th.App.GetTasks(&model.TaskGetOptions{
	// 		Page:           0,
	// 		PerPage:        10,
	// 		OwnerId:        "",
	// 		IncludeDeleted: true,
	// 	})
	// 	require.Nil(t, err)
	// 	assert.Equal(t, model.TaskList{task1, deletedTask1, task2, task3, task4, deletedTask2}, tasks)
	// })

	// t.Run("get tasks, page=0, perPage=1, include deleted", func(t *testing.T) {
	// 	tasks, err := th.App.GetTasks(&model.TaskGetOptions{
	// 		Page:           0,
	// 		PerPage:        1,
	// 		OwnerId:        "",
	// 		IncludeDeleted: true,
	// 	})
	// 	require.Nil(t, err)
	// 	assert.Equal(t, model.TaskList{task1}, tasks)
	// })

	// t.Run("get tasks, page=1, perPage=2, include deleted", func(t *testing.T) {
	// 	tasks, err := th.App.GetTasks(&model.TaskGetOptions{
	// 		Page:           1,
	// 		PerPage:        2,
	// 		OwnerId:        "",
	// 		IncludeDeleted: true,
	// 	})
	// 	require.Nil(t, err)
	// 	assert.Equal(t, model.TaskList{task2, task3}, tasks)
	// })

	// t.Run("get tasks, page=2, perPage=2, include deleted", func(t *testing.T) {
	// 	tasks, err := th.App.GetTasks(&model.TaskGetOptions{
	// 		Page:           2,
	// 		PerPage:        2,
	// 		OwnerId:        "",
	// 		IncludeDeleted: true,
	// 	})
	// 	require.Nil(t, err)
	// 	assert.Equal(t, model.TaskList{task4, deletedTask2}, tasks)
	// })

	// t.Run("get offset=0, limit=10, creator id 1", func(t *testing.T) {
	// 	tasks, err := th.App.GetTasks(&model.TaskGetOptions{
	// 		Page:           0,
	// 		PerPage:        10,
	// 		OwnerId:        OwnerId1,
	// 		IncludeDeleted: false,
	// 	})
	// 	require.Nil(t, err)
	// 	require.Equal(t, model.TaskList{task1, task2, task3}, tasks)
	// })

	// t.Run("get offset=0, limit=10, creator id 2", func(t *testing.T) {
	// 	tasks, err := th.App.GetTasks(&model.TaskGetOptions{
	// 		Page:           0,
	// 		PerPage:        10,
	// 		OwnerId:        OwnerId2,
	// 		IncludeDeleted: false,
	// 	})
	// 	require.Nil(t, err)
	// 	require.Equal(t, model.TaskList{task4}, tasks)
	// })

	// t.Run("get offset=0, limit=10, include deleted, creator id 1", func(t *testing.T) {
	// 	tasks, err := th.App.GetTasks(&model.TaskGetOptions{
	// 		Page:           0,
	// 		PerPage:        10,
	// 		OwnerId:        OwnerId1,
	// 		IncludeDeleted: true,
	// 	})
	// 	require.Nil(t, err)
	// 	require.Equal(t, model.TaskList{task1, deletedTask1, task2, task3}, tasks)
	// })

	// t.Run("get offset=0, limit=10, include deleted, creator id 2", func(t *testing.T) {
	// 	tasks, err := th.App.GetTasks(&model.TaskGetOptions{
	// 		Page:           0,
	// 		PerPage:        10,
	// 		OwnerId:        OwnerId2,
	// 		IncludeDeleted: true,
	// 	})
	// 	require.Nil(t, err)
	// 	require.Equal(t, model.TaskList{task4, deletedTask2}, tasks)
	// })
}

// func TestUpdateTaskActive(t *testing.T) {
// 	t.Run("unknown task", func(t *testing.T) {
// 		th := Setup(t).InitBasic()
// 		defer th.TearDown()

// 		_, err := th.App.UpdateTaskActive(model.NewId(), false)
// 		require.NotNil(t, err)
// 		require.Equal(t, "store.sql_user.missing_account.const", err.Id)
// 	})

// 	t.Run("disable/enable task", func(t *testing.T) {
// 		th := Setup(t).InitBasic()
// 		defer th.TearDown()

// 		task, err := th.App.CreateTask(&model.Task{
// 			Username:    "username",
// 			Description: "a task",
// 			OwnerId:     th.BasicUser.Id,
// 		})
// 		require.Nil(t, err)
// 		defer th.App.PermanentDeleteTask(task.UserId)

// 		disabledTask, err := th.App.UpdateTaskActive(task.UserId, false)
// 		require.Nil(t, err)
// 		require.NotEqual(t, 0, disabledTask.DeleteAt)

// 		// Disabling should be idempotent
// 		disabledTaskAgain, err := th.App.UpdateTaskActive(task.UserId, false)
// 		require.Nil(t, err)
// 		require.Equal(t, disabledTask.DeleteAt, disabledTaskAgain.DeleteAt)

// 		reenabledTask, err := th.App.UpdateTaskActive(task.UserId, true)
// 		require.Nil(t, err)
// 		require.EqualValues(t, 0, reenabledTask.DeleteAt)

// 		// Re-enabling should be idempotent
// 		reenabledTaskAgain, err := th.App.UpdateTaskActive(task.UserId, true)
// 		require.Nil(t, err)
// 		require.Equal(t, reenabledTask.DeleteAt, reenabledTaskAgain.DeleteAt)
// 	})
// }

// func TestPermanentDeleteTask(t *testing.T) {
// 	th := Setup(t).InitBasic()
// 	defer th.TearDown()

// 	task, err := th.App.CreateTask(&model.Task{
// 		Username:    "username",
// 		Description: "a task",
// 		OwnerId:     th.BasicUser.Id,
// 	})
// 	require.Nil(t, err)

// 	require.Nil(t, th.App.PermanentDeleteTask(task.UserId))

// 	_, err = th.App.GetTask(task.UserId, false)
// 	require.NotNil(t, err)
// 	require.Equal(t, "store.sql_task.get.missing.app_error", err.Id)
// }

// func TestDisableUserTasks(t *testing.T) {
// 	th := Setup(t).InitBasic()
// 	defer th.TearDown()

// 	ownerId1 := model.NewId()
// 	ownerId2 := model.NewId()

// 	tasks := []*model.Task{}
// 	defer func() {
// 		for _, task := range tasks {
// 			th.App.PermanentDeleteTask(task.UserId)
// 		}
// 	}()

// 	for i := 0; i < 46; i++ {
// 		task, err := th.App.CreateTask(&model.Task{
// 			Username:    fmt.Sprintf("username%v", i),
// 			Description: "a task",
// 			OwnerId:     ownerId1,
// 		})
// 		require.Nil(t, err)
// 		tasks = append(tasks, task)
// 	}
// 	require.Len(t, tasks, 46)

// 	u2task1, err := th.App.CreateTask(&model.Task{
// 		Username:    "username_nodisable",
// 		Description: "a task",
// 		OwnerId:     ownerId2,
// 	})
// 	require.Nil(t, err)
// 	defer th.App.PermanentDeleteTask(u2task1.UserId)

// 	err = th.App.disableUserTasks(ownerId1)
// 	require.Nil(t, err)

// 	// Check all tasks and corrensponding users are disabled for creator 1
// 	for _, task := range tasks {
// 		rettask, err2 := th.App.GetTask(task.UserId, true)
// 		require.Nil(t, err2)
// 		require.NotZero(t, rettask.DeleteAt, task.Username)
// 	}

// 	// Check tasks and corresponding user not disabled for creator 2
// 	task, err := th.App.GetTask(u2task1.UserId, true)
// 	require.Nil(t, err)
// 	require.Zero(t, task.DeleteAt)

// 	user, err := th.App.GetUser(u2task1.UserId)
// 	require.Nil(t, err)
// 	require.Zero(t, user.DeleteAt)

// 	// Bad id doesn't do anything or break horribly
// 	err = th.App.disableUserTasks(model.NewId())
// 	require.Nil(t, err)
// }

// func TestConvertUserToTask(t *testing.T) {
// 	t.Run("invalid user", func(t *testing.T) {
// 		t.Run("invalid user id", func(t *testing.T) {
// 			th := Setup(t).InitBasic()
// 			defer th.TearDown()

// 			_, err := th.App.ConvertUserToTask(&model.User{
// 				Username: "username",
// 				Id:       "",
// 			})
// 			require.NotNil(t, err)
// 			require.Equal(t, "model.task.is_valid.user_id.app_error", err.Id)
// 		})

// 		t.Run("invalid username", func(t *testing.T) {
// 			th := Setup(t).InitBasic()
// 			defer th.TearDown()

// 			_, err := th.App.ConvertUserToTask(&model.User{
// 				Username: "invalid username",
// 				Id:       th.BasicUser.Id,
// 			})
// 			require.NotNil(t, err)
// 			require.Equal(t, "model.task.is_valid.username.app_error", err.Id)
// 		})
// 	})

// 	t.Run("valid user", func(t *testing.T) {
// 		th := Setup(t).InitBasic()
// 		defer th.TearDown()

// 		task, err := th.App.ConvertUserToTask(&model.User{
// 			Username: "username",
// 			Id:       th.BasicUser.Id,
// 		})
// 		require.Nil(t, err)
// 		defer th.App.PermanentDeleteTask(task.UserId)
// 		assert.Equal(t, "username", task.Username)
// 		assert.Equal(t, th.BasicUser.Id, task.OwnerId)
// 	})
// }

// func sToP(s string) *string {
// 	return &s
// }
