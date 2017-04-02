//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package routes

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/daemon/api/views/v1"
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/errors"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func ProjectListH(w http.ResponseWriter, r *http.Request) {

	var (
		err     error
		session *types.Session
		ctx     = c.Get()
	)

	ctx.Log.Debug("List project handler")

	session = utils.Session(r)
	if session == nil {
		ctx.Log.Error(http.StatusText(http.StatusUnauthorized))
		errors.HTTP.Unauthorized(w)
		return
	}

	projectList, err := ctx.Storage.Project().ListByUser(session.Username)
	if err != nil {
		ctx.Log.Error("Error: find projects by user", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewProjectList(projectList).ToJson()
	if err != nil {
		ctx.Log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		ctx.Log.Error("Error: write response", err.Error())
		return
	}
}

func ProjectInfoH(w http.ResponseWriter, r *http.Request) {

	var (
		err          error
		session      *types.Session
		project      *types.Project
		ctx          = c.Get()
		params       = utils.Vars(r)
		projectParam = params["project"]
	)

	ctx.Log.Info("Get project handler")

	session = utils.Session(r)
	if session == nil {
		ctx.Log.Error(http.StatusText(http.StatusUnauthorized))
		errors.HTTP.Unauthorized(w)
		return
	}

	project, err = ctx.Storage.Project().GetByName(session.Username, projectParam)
	if err != nil {
		ctx.Log.Error("Error: find project by id", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if project == nil {
		errors.New("project").NotFound().Http(w)
		return
	}

	response, err := v1.NewProject(project).ToJson()
	if err != nil {
		ctx.Log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		ctx.Log.Error("Error: write response", err.Error())
		return
	}

}

//func ProjectActivityListH(w http.ResponseWriter, r *http.Request) {
//
//	var (
//		err          error
//		session      *types.Session
//		projectModel *types.Project
//		ctx          = c.Get()
//		params       = utils.Vars(r)
//		projectParam = params["project"]
//	)
//
//	ctx.Log.Debug("List project activity handler")
//
//	s, ok := context.GetOk(r, `session`)
//	if !ok {
//		ctx.Log.Error("Error: get session context")
//		e.New("user").Unauthorized().Http(w)
//		return
//	}
//
//	session = s.(*types.Session)
//
//	projectModel, err = ctx.Storage.Project().GetByNameOrID(session.Username, projectParam)
//	if err != nil {
//		ctx.Log.Error("Error: find project by id", err.Error())
//		e.HTTP.InternalServerError(w)
//		return
//	}
//	if projectModel == nil {
//		e.New("project").NotFound().Http(w)
//		return
//	}
//
//	activityListModel, err := ctx.Storage.Activity().ListProjectActivity(session.Username, projectModel.ID)
//	if err != nil {
//		ctx.Log.Error("Error: find projects by user", err)
//		e.HTTP.InternalServerError(w)
//		return
//	}
//
//	response, err := activityListModel.ToJson()
//	if err != nil {
//		ctx.Log.Error("Error: convert struct to json", err.Error())
//		e.HTTP.InternalServerError(w)
//		return
//	}
//
//	ctx.Log.Info(string(response))
//
//	w.WriteHeader(http.StatusOK)
//	_, err = w.Write(response)
//	if err != nil {
//		ctx.Log.Error("Error: write response", err.Error())
//		return
//	}
//}

type projectCreateS struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (s *projectCreateS) decodeAndValidate(reader io.Reader) *errors.Err {

	var (
		ctx = c.Get()
	)

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		ctx.Log.Error(err)
		return errors.New("user").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("project").IncorrectJSON(err)
	}

	if s.Name == nil {
		return errors.New("project").BadParameter("name")
	}

	*s.Name = strings.ToLower(*s.Name)

	if len(*s.Name) < 4 && len(*s.Name) > 64 && !validator.IsProjectName(*s.Name) {
		return errors.New("project").BadParameter("name")
	}

	if s.Description == nil {
		s.Description = new(string)
	}

	return nil
}

func ProjectCreateH(w http.ResponseWriter, r *http.Request) {

	var (
		session *types.Session
		ctx     = c.Get()
	)

	ctx.Log.Debug("Create project handler")

	session = utils.Session(r)
	if session == nil {
		ctx.Log.Error(http.StatusText(http.StatusUnauthorized))
		errors.HTTP.Unauthorized(w)
		return
	}

	// request body struct
	rq := new(projectCreateS)
	if err := rq.decodeAndValidate(r.Body); err != nil {
		ctx.Log.Error("Error: validation incomming data", err)
		errors.New("Invalid incomming data").Unknown().Http(w)
		return
	}

	project, err := ctx.Storage.Project().GetByName(session.Username, *rq.Name)
	if err != nil {
		ctx.Log.Error("Error: check exists by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	if project != nil {
		errors.New("project").NotUnique("name").Http(w)
		return
	}

	project, err = ctx.Storage.Project().Insert(session.Username, *rq.Name, *rq.Description)
	if err != nil {
		ctx.Log.Error("Error: insert project to db", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewProject(project).ToJson()
	if err != nil {
		ctx.Log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		ctx.Log.Error("Error: write response", err.Error())
		return
	}
}

//type projectUpdateS struct {
//	Name        *string `json:"name,omitempty"`
//	Description *string `json:"description,omitempty"`
//}
//
//func (s *projectUpdateS) decodeAndValidate(reader io.Reader) *e.Err {
//
//	var (
//		err error
//		ctx = c.Get()
//	)
//
//	body, err := ioutil.ReadAll(reader)
//	if err != nil {
//		ctx.Log.Error(err)
//		return e.New("user").Unknown(err)
//	}
//
//	err = json.Unmarshal(body, s)
//	if err != nil {
//		return e.New("project").IncorrectJSON(err)
//	}
//
//	if s.Name == nil {
//		return e.New("project").BadParameter("name")
//	}
//
//	*s.Name = strings.ToLower(*s.Name)
//
//	if len(*s.Name) < 4 && len(*s.Name) > 64 && !validator.IsProjectName(*s.Name) {
//		return e.New("project").BadParameter("name")
//	}
//
//	if s.Description == nil {
//		s.Description = new(string)
//	}
//
//	return nil
//}
//
//func ProjectUpdateH(w http.ResponseWriter, r *http.Request) {
//
//	var (
//		err          error
//		session      *types.Session
//		projectModel *types.Project
//		ctx          = c.Get()
//		params       = utils.Vars(r)
//		projectParam = params["project"]
//	)
//
//	ctx.Log.Debug("Update project handler")
//
//	s, ok := context.GetOk(r, `session`)
//	if !ok {
//		ctx.Log.Error("Error: get session context")
//		e.New("user").Unauthorized().Http(w)
//		return
//	}
//
//	session = s.(*types.Session)
//
//	// request body struct
//	rq := new(projectUpdateS)
//	if err := rq.decodeAndValidate(r.Body); err != nil {
//		ctx.Log.Error("Error: validation incomming data", err)
//		err.Http(w)
//		return
//	}
//
//	projectModel, err = ctx.Storage.Project().GetByNameOrID(session.Uid, projectParam)
//	if err != nil {
//		ctx.Log.Error("Error: find project by id", err.Error())
//		e.HTTP.InternalServerError(w)
//		return
//	}
//	if projectModel == nil {
//		e.New("project").NotFound().Http(w)
//		return
//	}
//
//	if rq.Name == nil || *rq.Name == "" {
//		rq.Name = &projectModel.Name
//	}
//
//	if !validator.IsUUID(projectParam) && projectModel.Name != *rq.Name {
//		exists, err := ctx.Storage.Project().ExistByName(projectModel.User, *rq.Name)
//		if err != nil {
//			e.HTTP.InternalServerError(w)
//		}
//		if exists {
//			e.New("project").NotUnique("name").Http(w)
//			return
//		}
//	}
//
//	projectModel.Name = *rq.Name
//	projectModel.Description = *rq.Description
//
//	projectModel, err = ctx.Storage.Project().Update(projectModel)
//	if err != nil {
//		ctx.Log.Error("Error: insert project to db", err.Error())
//		e.HTTP.InternalServerError(w)
//		return
//	}
//
//	response, err := projectModel.ToJson()
//	if err != nil {
//		ctx.Log.Error("Error: convert struct to json", err.Error())
//		e.HTTP.InternalServerError(w)
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//	_, err = w.Write(response)
//	if err != nil {
//		ctx.Log.Error("Error: write response", err.Error())
//		return
//	}
//}

func ProjectRemoveH(w http.ResponseWriter, r *http.Request) {

	var (
		ctx          = c.Get()
		session      *types.Session
		params       = utils.Vars(r)
		projectParam = params["project"]
	)

	ctx.Log.Info("Remove project")

	session = utils.Session(r)
	if session == nil {
		ctx.Log.Error(http.StatusText(http.StatusUnauthorized))
		errors.HTTP.Unauthorized(w)
		return
	}

	projectModel, err := ctx.Storage.Project().GetByName(session.Username, projectParam)
	if err != nil {
		ctx.Log.Error("Error: find project by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if projectModel == nil {
		errors.New("project").NotFound().Http(w)
		return
	}

	// Todo: remove all services by project id
	// Todo: remove all activity by project id

	//err = ctx.Storage.Service().RemoveByProject(session.Username, projectParam)
	//if err != nil {
	//	ctx.Log.Error("Error: remove services from db", err)
	//	e.HTTP.InternalServerError(w)
	//	return
	//}

	//err = ctx.Storage.Activity().RemoveByProject(session.Username, projectParam)
	//if err != nil {
	//	ctx.Log.Error("Error: remove activity from db", err)
	//	e.HTTP.InternalServerError(w)
	//	return
	//}

	err = ctx.Storage.Project().Remove(session.Username, projectParam)
	if err != nil {
		ctx.Log.Error("Error: remove project from db", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte{})
	if err != nil {
		ctx.Log.Error("Error: write response", err.Error())
		return
	}
}
