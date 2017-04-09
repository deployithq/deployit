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

package service

import (
	s "github.com/lastbackend/lastbackend/pkg/apis/views/v1/service"
	"github.com/lastbackend/lastbackend/pkg/client/cmd/project"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"github.com/lastbackend/lastbackend/pkg/errors"
)

func InspectCmd(name string) {

	ctx := context.Get()

	service, projectName, err := Inspect(name)
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	service.DrawTable(projectName)
}

func Inspect(name string) (*s.Service, string, error) {

	var (
		err     error
		ctx     = context.Get()
		er      = new(errors.Http)
		service = new(s.Service)
	)

	p, err := project.Current()
	if err != nil {
		return nil, "", errors.New(err.Error())
	}

	if p == nil {
		ctx.Log.Info("Project didn't select")
		return nil, "", nil
	}

	_, _, err = ctx.HTTP.
		GET("/project/"+p.Name+"/service/"+name).
		AddHeader("Authorization", "Bearer "+ctx.Token).
		Request(service, er)

	if err != nil {
		return nil, "", errors.New(err.Error())
	}

	if er.Code == 401 {
		return nil, "", errors.NotLoggedMessage
	}

	if er.Code != 0 {
		return nil, "", errors.New(er.Message)
	}

	return service, p.Name, nil
}
