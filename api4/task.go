//  Copyright  (c)  2017-present  Xenia,  Inc.  All  Rights  Reserved.
//  See  License.txt  for  license  information.

package  api4

import  (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	//  "github.com/xzl8028/xenia-server/mlog"
	"github.com/xzl8028/xenia-server/model"
)



type payload2 struct{

	Id string`json:"id"`
	Is_pinned bool`json:"is_pinned"`
	Message string `json:"message"`
	File_ids []string`json:"file_ids"`
	Has_reactions bool `json:"has_reactions"`
	Props props `json:"props"`


}


type  payload  struct{
	Channel_id    string  `json:"channel_id"`
	Message  string  `json:"message"`
	Root_id  string  `json:"root_id"`
	//File_ids  []string`json："file_ids"`
	Props  props  `json:"props"`
}

type  props  struct  {
	From_webhook  string  `json:"from_webhook"`
	Override_icon_url  string  `json:"override_icon_url"`
	Override_username  string  `json:"override_username"`
	Webhook_display_name  string  `json:"webhook_display_name"`
	Attachments  []attachment  `json:"attachments"`
}


type  attachment  struct  {
	Color string `json:"color"`
	Fields    []field  `json:"fields"`
	Actions  []action  `json:"actions"`
}


type  field  struct{
	Short  bool  `json:"short"`
	Title  string  `json:"title"`
	Value  string  `json:"value"`
}

type  action  struct{
	Name  string  `json:"name"`
	Integration    integration_url  `json:"integration"`
}


type  integration_url  struct{
	Url    string  `json:"url"`
}

type  Host  struct  {
	IP  string
	Name  string
}


func  (api  *API)  InitTask()  {
	//  api.BaseRoutes.Tasks.Handle("",  api.ApiSessionRequired(createTask)).Methods("POST")
	//  api.BaseRoutes.Task.Handle("",  api.ApiSessionRequired(patchTask)).Methods("PUT")
	api.BaseRoutes.Task.Handle("",  api.ApiSessionRequired(getTask)).Methods("GET")
	api.BaseRoutes.Tasks.Handle("",  api.ApiSessionRequired(getTasks)).Methods("GET")
	api.BaseRoutes.Tasks.Handle("/withteam",  api.ApiSessionRequired(getTasksWithTeam)).Methods("GET")
	api.BaseRoutes.Task.Handle("/insert",  api.ApiSessionRequired(insertTask)).Methods("POST")
	api.BaseRoutes.Task.Handle("/",  api.ApiSessionRequired(insertTask)).Methods("POST")
	api.BaseRoutes.Task.Handle("/insertpost",  api.ApiSessionRequired(insertTaskWithPost)).Methods("POST")
	api.BaseRoutes.Task.Handle("/update",  api.ApiSessionRequired(updateTask)).Methods("POST")
	api.BaseRoutes.Task.Handle("/updatepost",  api.ApiSessionRequired(updateTaskStatusQuickWithPost)).Methods("POST")
	api.BaseRoutes.Task.Handle("/update_status_quick",  api.ApiSessionRequired(updateTaskStatusQuick)).Methods("POST")

	//  api.BaseRoutes.Task.Handle("/disable",  api.ApiSessionRequired(disableTask)).Methods("POST")
	//  api.BaseRoutes.Task.Handle("/enable",  api.ApiSessionRequired(enableTask)).Methods("POST")
	//  api.BaseRoutes.Task.Handle("/assign/{user_id:[A-Za-z0-9]+}",  api.ApiSessionRequired(assignTask)).Methods("POST")
}

//  func  createTask(c  *Context,  w  http.ResponseWriter,  r  *http.Request)  {
//  taskPatch  :=  model.TaskPatchFromJson(r.Body)
//  if  taskPatch  ==  nil  {
//  c.SetInvalidParam("task")
//  return
//  }

//  task  :=  &model.Task{
//  OwnerId:  c.App.Session.UserId,
//  }
//  task.Patch(taskPatch)

//  if  !c.App.SessionHasPermissionTo(c.App.Session,  model.PERMISSION_CREATE_TASK)  {
//  c.SetPermissionError(model.PERMISSION_CREATE_TASK)
//  return
//  }

//  if  user,  err  :=  c.App.GetUser(c.App.Session.UserId);  err  ==  nil  {
//  if  user.IsTask  {
//  c.SetPermissionError(model.PERMISSION_CREATE_TASK)
//  return
//  }
//  }

//  if  !*c.App.Config().ServiceSettings.EnableTaskAccountCreation  {
//  c.Err  =  model.NewAppError("createTask",  "api.task.create_disabled",  nil,  "",  http.StatusForbidden)
//  return
//  }

//  createdTask,  err  :=  c.App.CreateTask(tas k)
// 	if err != nil {
// 	c.Err = err
// 	return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	w.Write(createdTask.ToJson())
// }

// func patchTask(c *Context, w http.ResponseWriter, r *http.Request) {
// 	c.RequireTaskId()
// 	if c.Err != nil {
// 	return
// 	}
// 	TaskId := c.Params.TaskId

// 	taskPatch := model.TaskPatchFromJson(r.Body)
// 	if taskPatch == nil {
// 	c.SetInvalidParam("task")
// 	return
// 	}

// 	if err := c.App.SessionHasPermissionToManageTask(c.App.Session, TaskId); err != nil {
// 	c.Err = err
// 	return
// 	}

// 	updatedTask, err := c.App.PatchTask(TaskId, taskPatch)
// 	if err != nil {
// 	c.Err = err
// 	return
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
	// 	// Pretend like the task doesn't exist at all to avoid revealing that the
	// 	// user is a task. It's kind of silly in this case, sine we created the task,
	// 	// but we don't have read task permissions.
	// 	c.Err = model.MakeTaskNotFoundError(TaskId)
	// 	return
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


func getTasksWithTeam(c *Context, w http.ResponseWriter, r *http.Request) {
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

	teams,errteam := c.App.GetTeamMembersForUser(c.App.Session.UserId)
	if(errteam!=nil){
		return;
	}


	tasks, err := c.App.GetAllWithTeamId(teams[0].TeamId)
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


func updateTaskStatusQuickWithPost(c *Context, w http.ResponseWriter, r *http.Request) {
	var res string
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

	receiveDeptToId := make(map[string]string)
	receiveDeptToId["housekeeping"] = "r86qtsc49pbw7cdzigptsgm47a"
	fmt.Println("!!!POSTTIPM!!!new task.roomid", new_task.DueAt)
	var due_at_int64 time.Time = time.Unix(0, int64(new_task.DueAt) * int64(time.Millisecond))
	due_at_string := due_at_int64.Format("01月02日 15:04")

	//field struct

	field2 := field{Title: "客房号", Value: strconv.Itoa(new_task.RoomId), Short: true}
	field3 := field{Title: "要求完成时间", Value: due_at_string, Short: true}
	field4 := field{Title: "任务内容", Value: new_task.TaskType, Short: true}
	field5 := field{Title: "发单部门", Value: new_task.SendDept, Short: true}
	field6 := field{Title: "备注", Value: new_task.Note, Short: true}


	if r.URL.Query().Get("status") == "1" && old_task.Status == 0 {
		new_task.ConfirmAt = model.GetMillis()
		new_task.Status = 1
		integration_url_self := integration_url{Url: "http://localhost:8065/api/v4/tasks/"+ strconv.Itoa(new_task.TaskId) +"/updatepost?status=2"}
		action_obj := action{ Name: "确认完成", Integration: integration_url_self }
		action_self := []action{action_obj}
		field1 := field{Title: "任务状态", Value: "**执行中**", Short: true}
		field_self := []field{field1, field2, field3, field4, field5, field6}
		//attachment struct
		attachment_obj := attachment{Color:"#2CACE1", Fields:field_self, Actions:action_self}
		attachment_self := []attachment{attachment_obj}

		//props struct
		//props_self := props{Attachments: attachment_self }
		props_self := props{From_webhook:"true", Override_icon_url: "http://s575.com/Uploads/2018-10-31/20170pwu61540976213.png", Override_username: "灵奇任务助手", Webhook_display_name: "task_center", Attachments: attachment_self }


		//payload struct
		payload_self := payload2{Id: old_task.PostId, Is_pinned:false,Message: "任务"+strconv.Itoa(new_task.TaskId),File_ids:[]string{},Has_reactions:false, Props: props_self }


		//strinify playload
		jsonData,err2 := json.Marshal(payload_self)
		fmt.Println(string(jsonData)+"!!!!! JSON FILE OF PLAYLOAD AFTER MARSHAL !!!! ")
		if(err2!=nil){
			fmt.Println("!!!!!!! POSITION1 !!!!!!!", err2)
			return
		}



		var r2 http.Request
		r2.Body =  ioutil.NopCloser(bytes.NewReader(jsonData))
		c.Params.PostId = old_task.PostId
		res = updatePostWithReturn(c,w,&r2)


	} else if r.URL.Query().Get("status") == "2" && (old_task.Status == 1 || old_task.Status == 4){
		new_task.FinishAt = model.GetMillis()
		new_task.Status = 2
		field1 := field{Title: "任务状态", Value: "**已完成**", Short: true}
		field_self := []field{field1, field2, field3, field4, field5, field6}
		//attachment struct
		attachment_obj := attachment{Color:"#258A28", Fields:field_self}
		attachment_self := []attachment{attachment_obj}

		//props struct
		props_self := props{Attachments: attachment_self }

		//payload struct
		payload_self := payload2{Id: old_task.PostId, Is_pinned:false,Message: "任务"+strconv.Itoa(new_task.TaskId),File_ids:[]string{},Has_reactions:false, Props: props_self }


		//strinify playload
		jsonData,err2 := json.Marshal(payload_self)
		fmt.Println(string(jsonData)+"!!!!! JSON FILE OF PLAYLOAD AFTER MARSHAL !!!! ")
		if(err2!=nil){
			fmt.Println("!!!!!!! POSITION1 !!!!!!!", err2)
			return
		}



		var r2 http.Request
		r2.Body =  ioutil.NopCloser(bytes.NewReader(jsonData))
		c.Params.PostId = old_task.PostId
		res = updatePostWithReturn(c,w,&r2)

	} else if r.URL.Query().Get("status") == "3" && old_task.Status != 3 && old_task.Status != 2 {
		new_task.Status = 3
		field1 := field{Title: "任务状态", Value: "**已取消**", Short: true}
		field_self := []field{field1, field2, field3, field4, field5, field6}
		//attachment struct
		attachment_obj := attachment{Color:"#3D3C40", Fields:field_self}
		attachment_self := []attachment{attachment_obj}

		//props struct
		props_self := props{From_webhook:"true", Override_icon_url: "http://s575.com/Uploads/2018-10-31/20170pwu61540976213.png", Override_username: "灵奇任务助手", Webhook_display_name: "task_center", Attachments: attachment_self }


		//payload struct
		payload_self := payload2{Id: old_task.PostId, Is_pinned:false,Message: "任务"+strconv.Itoa(new_task.TaskId),File_ids:[]string{},Has_reactions:false, Props: props_self }


		//strinify playload
		jsonData,err2 := json.Marshal(payload_self)
		fmt.Println(string(jsonData)+"!!!!! JSON FILE OF PLAYLOAD AFTER MARSHAL !!!! ")
		if(err2!=nil){
			fmt.Println("!!!!!!! POSITION1 !!!!!!!", err2)
			return
		}

		var r2 http.Request
		r2.Body =  ioutil.NopCloser(bytes.NewReader(jsonData))
		c.Params.PostId = old_task.PostId
		res = updatePostWithReturn(c,w,&r2)


	} else if r.URL.Query().Get("status") == "4" && old_task.Status != 3 && old_task.Status != 2 && old_task.Status != 4 {
		//超时是4
		new_task.Status = 4
		integration_url_self := integration_url{Url: "http://localhost:8065/api/v4/tasks/"+ strconv.Itoa(new_task.TaskId) +"/updatepost?status=2"}
		action_obj := action{ Name: "确认完成", Integration: integration_url_self }
		action_self := []action{action_obj}
		field1 := field{Title: "任务状态", Value: "**已超时**", Short: true}
		field_self := []field{field1, field2, field3, field4, field5, field6}
		//attachment struct
		attachment_obj := attachment{Color:"#FC3D41", Fields:field_self, Actions:action_self}
		attachment_self := []attachment{attachment_obj}

		//props struct
		props_self := props{From_webhook:"true", Override_icon_url: "http://s575.com/Uploads/2018-10-31/20170pwu61540976213.png", Override_username: "灵奇任务助手", Webhook_display_name: "task_center", Attachments: attachment_self }

		//payload struct
		payload_self := payload2{Id: old_task.PostId, Is_pinned:false,Message: "任务"+strconv.Itoa(new_task.TaskId),File_ids:[]string{},Has_reactions:false, Props: props_self }


		//strinify playload
		jsonData,err2 := json.Marshal(payload_self)
		fmt.Println(string(jsonData)+"!!!!! JSON FILE OF PLAYLOAD AFTER MARSHAL !!!! ")
		if(err2!=nil){
			fmt.Println("!!!!!!! POSITION1 !!!!!!!", err2)
			return
		}

		var r2 http.Request
		r2.Body =  ioutil.NopCloser(bytes.NewReader(jsonData))
		c.Params.PostId = old_task.PostId
		res = updatePostWithReturn(c,w,&r2)


	}
	new_task.PostId = res
	task, err := c.App.UpdateTask(new_task)
	if err != nil {
		c.Err = err
		return
	}

	w.Write(task.ToJson())
}








func insertTaskWithPost(c *Context, w http.ResponseWriter, r *http.Request) {
	fmt.Println("!!!!!HERE!!!!")
	// new_task := model.TaskFromJson(r.Body)
	var new_task model.Task

	fmt.Println(c.App.Session.TeamMembers[0],"0 here!!!")
	fmt.Println(len(c.App.Session.TeamMembers),"length!!!")
	//fmt.Println(c.App.Session.TeamMembers[1],"1 here!!!")

	new_task = insertHelper(new_task,r)

	task, err := c.App.InsertTask(&new_task)
	if err != nil {
		c.Err = err
		return
	}


	// !!!! 做调用create post这个API的数据准备了 ！！！！ //

	//create a map { receive_dept: channel_id }
	//此处需要根据不同的酒店数据库情况修改
	//未来加入更多匹配值
	//receiveDeptToId := make(map[string]string)
	//receiveDeptToId["housekeeping"] = "r86qtsc49pbw7cdzigptsgm47a"

	teams,_ := c.App.GetTeamMembersForUser(c.App.Session.UserId)

	channel,errdisplay := c.App.GetByDisplayName(teams[0].TeamId,task.ReceiveDept)
	if errdisplay!=nil{
		c.Err = errdisplay
		return
	}

	// integration_url struct
	integration_url_self := integration_url{Url: "http://localhost:8065/api/v4/tasks/"+ strconv.Itoa(task.TaskId) +"/updatepost?status=1"}

	// action struct
	action_obj := action{ Name: "确认接收", Integration: integration_url_self }
	action_self := []action{action_obj}

	//convert due_at into normal Date format for field struct
	fmt.Println("!!!POSTTIPM!!!new task.roomid", new_task.DueAt)
	var due_at_int64 time.Time = time.Unix(0, int64(new_task.DueAt) * int64(time.Millisecond))
	due_at_string := due_at_int64.Format("01月02日 15:04")

	//field struct
	field1 := field{Title: "任务状态", Value: "**等待接收**", Short: true}
	field2 := field{Title: "客房号", Value: strconv.Itoa(new_task.RoomId), Short: true}
	field3 := field{Title: "要求完成时间", Value: due_at_string, Short: true}
	field4 := field{Title: "任务内容", Value: new_task.TaskType, Short: true}
	field5 := field{Title: "发单部门", Value: new_task.SendDept, Short: true}
	field6 := field{Title: "备注", Value: new_task.Note, Short: true}
	field_self := []field{field1, field2, field3, field4, field5, field6}

	//attachment struct
	attachment_obj := attachment{Color:"#9C58E3" ,Fields:field_self, Actions:action_self}
	attachment_self := []attachment{attachment_obj}

	//props struct
	props_self := props{From_webhook:"true", Override_icon_url: "http://s575.com/Uploads/2018-10-31/20170pwu61540976213.png", Override_username: "灵奇任务助手", Webhook_display_name: "task_center", Attachments: attachment_self }



	//payload struct
	payload_self := payload{Channel_id: channel.Id, Message: "任务"+strconv.Itoa(task.TaskId), Root_id: "", Props: props_self }


	//strinify playload
	jsonData,err2 := json.Marshal(payload_self)
	fmt.Println(string(jsonData)+"!!!!! JSON FILE OF PLAYLOAD AFTER MARSHAL !!!! ")
	if(err2!=nil){
		fmt.Println("!!!!!!! POSITION1 !!!!!!!", err2)
		return
	}

	post  :=  model.PostFromJson(bytes.NewBuffer(jsonData))
	if  post  ==  nil  {
		c.SetInvalidParam("post")
		return
	}
	fmt.Println("!!!!insert task with post")
	var r2 http.Request
	r2.Body =  ioutil.NopCloser(bytes.NewReader(jsonData))

	res := createPostWithReturn(c,w,&r2)
	task.PostId = res
	task.TeamId = teams[0].TeamId
	tasklateset, err := c.App.UpdateTask(task)



	w.Write(tasklateset.ToJson())
}
func insertTask(c *Context, w http.ResponseWriter, r *http.Request) {
	var new_task model.Task


	new_task = insertHelper(new_task,r)
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


func insertHelper(new_task model.Task,r *http.Request)(res model.Task){
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
	return new_task
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

// func assignTask(c *Context, w http.ResponseWriter, r *http.Request) {
// 	c.RequireUserId()
// 	c.RequireTaskId()
// 	if c.Err != nil {
// 	return
// 	}
// 	TaskId := c.Params.TaskId
// 	userId := c.Params.UserId

// 	if err := c.App.SessionHasPermissionToManageTask(c.App.Session, TaskId); err != nil {
// 	c.Err = err
// 	return
// 	}

// 	if user, err := c.App.GetUser(userId); err == nil {
// 	if user.IsTask {
// 	c.SetPermissionError(model.PERMISSION_ASSIGN_TASK)
// 	return
// 	}
// 	}

// 	task, err := c.App.UpdateTaskOwner(TaskId, userId)
// 	if err != nil {
// 	c.Err = err
// 	return
// 	}

// 	w.Write(task.ToJson())
// }
