package api

import (
	"net/http"
	"strings"

	"github.com/docker/swarmkit/api"
	"github.com/gorilla/mux"
	ct "golang.org/x/net/context"
)

// GET /tasks?all=1&quiet=1
//    all:0 only display running
//		  1 display all
//	  default 0
func listTasks(c *context, w http.ResponseWriter, r *http.Request) {
	var (
		allStr       = r.URL.Query().Get("all")
		all          = false
		err          error
		listTaskResp *api.ListTasksResponse
		tasks        []*api.Task
	)

	if len(strings.TrimSpace(allStr)) != 0 && allStr == "1" {
		all = true
	}

	if listTaskResp, err = c.swarmkitAPI.ListTasks(ct.TODO(), &api.ListTasksRequest{}); err != nil {
		errResponse(w, r, err, c)
		return
	}
	for _, t := range listTaskResp.Tasks {
		if all || t.DesiredState <= api.TaskStateRunning {
			tasks = append(tasks, t)
		}
	}

	c.render.JSON(w, http.StatusOK, tasks)
}

// GET /tasks/{taskid:.*}
func inspectTasks(c *context, w http.ResponseWriter, r *http.Request) {
	var (
		taskid      = mux.Vars(r)["taskid"]
		err         error
		getTaskResp *api.GetTaskResponse
	)

	if getTaskResp, err = c.swarmkitAPI.GetTask(ct.TODO(), &api.GetTaskRequest{TaskID: taskid}); err != nil {
		errResponse(w, r, err, c)
		return
	}

	c.render.JSON(w, http.StatusOK, getTaskResp.Task)
}

//DELETE /tasts/{taskid:.*}
func removeTasks(c *context, w http.ResponseWriter, r *http.Request) {
	var (
		taskid = mux.Vars(r)["taskid"]
		err    error
	)

	if _, err = c.swarmkitAPI.RemoveTask(ct.TODO(), &api.RemoveTaskRequest{TaskID: taskid}); err != nil {
		errResponse(w, r, err, c)
		return
	}
	c.render.JSON(w, http.StatusOK, taskid)
}
