package api

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/docker/swarmkit/api"
	"github.com/gorilla/mux"
	ct "golang.org/x/net/context"
)

// GET /services
func listService(c *context, w http.ResponseWriter, r *http.Request) {
	sresp, err := c.swarmkitAPI.ListServices(ct.TODO(), &api.ListServicesRequest{})
	if err != nil {
		errResponse(w, r, err, c)
		return
	}

	c.render.JSON(w, http.StatusOK, sresp.Services)
}

// POST /services/create
// {
//    name:"redis",
//    image:"redis:3.0.5",
//    labels:{"com.docker.test":"test"},        // service label (key=value)
//    mode:"",                                  // one of replicated, global
//    replicas: 1,                              // number of replicas for the service (only works in replicated service mode)
//    image: "redis:3.0.5",                     // container image
//    args: [],                                 // container args
//    env: [],                                  // container env
//    ports: [],                                // ports
//    network:"",                               // network name
//    memory-reservation: "",                   // amount of reserved memory (e.g. 512m)
//    memory-limit: "",                         // memory limit (e.g. 512m)
//    cpu-reservation:"",                       // number of CPU cores reserved (e.g. 0.5)
//    cpu-limit:"",                             // CPU cores limit (e.g. 0.5)
//    update-parallelism:0,                     // task update parallelism (0 = all at once)
//    update-delay:"0s",                        // delay between task updates (0s = none)
//    restart-condition:"any",                  // condition to restart the task (any, failure, none)
//    restart-delay:"5s",                       // delay between task restarts
//    restart-max-attempts:0,                   // maximum number of restart attempts (0 = unlimited)
//    restart-window:"0s",                      // time window to evaluate restart attempts (0 = unbound)
//    constraint:[],                            // Placement constraint (node.labels.key==value)
//    bind:[],                                  // define a bind mount
//    volume:[],                                // define a volume mount
// }
func createService(c *context, w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		cspec = &createSpec{}
	)
	if err = DecoderRequest(r, cspec); err != nil {
		err = fmt.Errorf("Parse params for create service error:%v", err)
		errResponse(w, r, err, c)
		return
	}

	if len(strings.TrimSpace(cspec.Name)) == 0 || len(strings.TrimSpace(cspec.Image)) == 0 {
		err = fmt.Errorf("name and image are mandatory")
		errResponse(w, r, err, c)
		return
	}

	spec := &api.ServiceSpec{
		Mode: &api.ServiceSpec_Replicated{
			Replicated: &api.ReplicatedService{
				Replicas: 1,
			},
		},
		Task: api.TaskSpec{
			Runtime: &api.TaskSpec_Container{
				Container: &api.ContainerSpec{},
			},
		},
	}

	if err = merge(cspec, spec, c.swarmkitAPI); err != nil {
		errResponse(w, r, err, c)
		return
	}

	var csResp *api.CreateServiceResponse
	if csResp, err = c.swarmkitAPI.CreateService(ct.TODO(), &api.CreateServiceRequest{Spec: spec}); err != nil {
		errResponse(w, r, err, c)
		return
	}

	c.render.JSON(w, http.StatusOK, csResp.Service)
}

// GET /services/{serviceid:.*}?all=1
//    all:0 only display running
//		  1 display all
//	  default 0
func inspectService(c *context, w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		gsResp    *api.GetServiceResponse
		lsTask    *api.ListTasksResponse
		serviceid = mux.Vars(r)["serviceid"]
		tasks     = make([]*api.Task, 0)
	)

	if gsResp, err = c.swarmkitAPI.GetService(ct.TODO(), &api.GetServiceRequest{ServiceID: serviceid}); err != nil {
		errResponse(w, r, err, c)
		return
	}

	if lsTask, err = c.swarmkitAPI.ListTasks(ct.TODO(), &api.ListTasksRequest{
		Filters: &api.ListTasksRequest_Filters{
			ServiceIDs: []string{gsResp.Service.ID},
		},
	}); err != nil {
		errResponse(w, r, err, c)
		return
	}

	for _, t := range lsTask.Tasks {
		if t.Status.State == api.TaskStateRunning {
			tasks = append(tasks, t)
		}
	}

	c.render.JSON(w, http.StatusOK, map[string]interface{}{
		"service": gsResp.Service,
		"tasks":   tasks,
	})
}

// POST /services/{serviceid:.*}/update
func updateService(c *context, w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		gsResp    *api.GetServiceResponse
		cspec     = &createSpec{}
		serviceid = mux.Vars(r)["serviceid"]
	)

	if len(strings.TrimSpace(serviceid)) <= 1 {
		err = errors.New("service ID missing")
		errResponse(w, r, err, c)
		return
	}

	if err = DecoderRequest(r, cspec); err != nil {
		err = fmt.Errorf("Parse params for create service error:%v", err)
		errResponse(w, r, err, c)
		return
	}

	if gsResp, err = c.swarmkitAPI.GetService(ct.TODO(), &api.GetServiceRequest{ServiceID: serviceid}); err != nil {
		errResponse(w, r, err, c)
		return
	}

	service := gsResp.Service
	spec := service.Spec.Copy()
	if err = merge(cspec, spec, c.swarmkitAPI); err != nil {
		errResponse(w, r, err, c)
		return
	}

	if reflect.DeepEqual(spec, &service.Spec) {
		err = errors.New("no changes detected")
		errResponse(w, r, err, c)
		return
	}

	var usResp *api.UpdateServiceResponse
	if usResp, err = c.swarmkitAPI.UpdateService(ct.TODO(), &api.UpdateServiceRequest{
		ServiceID:      service.ID,
		ServiceVersion: &service.Meta.Version,
		Spec:           spec,
	}); err != nil {
		errResponse(w, r, err, c)
		return
	}

	c.render.JSON(w, http.StatusOK, map[string]interface{}{"id": usResp.Service.ID})
}

// DELETE /services/{serviceid:.*}
func deleteService(c *context, w http.ResponseWriter, r *http.Request) {}
