//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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

package v1

import (
	"context"
	"fmt"
	rv1 "github.com/lastbackend/lastbackend/internal/api/types/v1/request"
	"github.com/lastbackend/lastbackend/internal/api/types/v1/views"
	"github.com/lastbackend/lastbackend/pkg/client/types"
	"strconv"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	t "github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/internal/util/http/request"
)

type DeploymentClient struct {
	client *request.RESTClient

	namespace t.NamespaceSelfLink
	service   t.ServiceSelfLink
	selflink  t.DeploymentSelfLink
}

func (dc *DeploymentClient) Pod(args ...string) types.PodClientV1 {
	name := ""
	// Get any parameters passed to us out of the args variable into "real"
	// variables we created for them.
	for i := range args {
		switch i {
		case 0: // hostname
			name = args[0]
		default:
			panic("Wrong parameter count: (is allowed from 0 to 1)")
		}
	}
	return newPodClient(dc.client, dc.namespace.String(), t.KindDeployment, dc.selflink.String(), name)
}

func (dc *DeploymentClient) List(ctx context.Context) (*views.DeploymentList, error) {

	var s *views.DeploymentList
	var e *errors.Http

	err := dc.client.Get(fmt.Sprintf("/namespace/%s/service/%s/deployment", dc.namespace.String(), dc.service.Name())).
		AddHeader("Content-Type", "application/json").
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	if s == nil {
		list := make(views.DeploymentList, 0)
		s = &list
	}

	return s, nil
}

func (dc *DeploymentClient) Get(ctx context.Context) (*views.Deployment, error) {

	var s *views.Deployment
	var e *errors.Http

	err := dc.client.Get(fmt.Sprintf("/namespace/%s/service/%s/deployment/%s", dc.namespace.String(), dc.service.Name(), dc.selflink.Name())).
		AddHeader("Content-Type", "application/json").
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	return s, nil
}

func (dc *DeploymentClient) Create(ctx context.Context, opts *rv1.DeploymentManifest) (*views.Deployment, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	var s *views.Deployment
	var e *errors.Http

	err = dc.client.Post(fmt.Sprintf("/namespace/%s/service/%s/deployment", dc.namespace.String(), dc.service.Name())).
		AddHeader("Content-Type", "application/json").
		Body(body).
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	return s, nil
}

func (dc *DeploymentClient) Update(ctx context.Context, opts *rv1.DeploymentManifest) (*views.Deployment, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	var s *views.Deployment
	var e *errors.Http

	err = dc.client.Put(fmt.Sprintf("/namespace/%s/service/%s/deployment/%s", dc.namespace.String(), dc.service.Name(), dc.selflink.Name())).
		AddHeader("Content-Type", "application/json").
		Body(body).
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	return s, nil
}

func (dc *DeploymentClient) Remove(ctx context.Context, opts *rv1.DeploymentRemoveOptions) error {

	req := dc.client.Delete(fmt.Sprintf("/namespace/%s/service/%s/deployment/%s", dc.namespace.String(), dc.service.Name(), dc.selflink.Name())).
		AddHeader("Content-Type", "application/json")

	if opts != nil {
		if opts.Force {
			req.Param("force", strconv.FormatBool(opts.Force))
		}
	}

	var e *errors.Http

	if err := req.JSON(nil, &e); err != nil {
		return err
	}
	if e != nil {
		return errors.New(e.Message)
	}

	return nil
}

func newDeploymentClient(client *request.RESTClient, namespace, service, name string) *DeploymentClient {
	return &DeploymentClient{
		client:    client,
		namespace: *t.NewNamespaceSelfLink(namespace),
		service:   *t.NewServiceSelfLink(namespace, service),
		selflink:  *t.NewDeploymentSelfLink(namespace, service, name),
	}
}
