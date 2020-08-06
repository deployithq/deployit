//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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
	rv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"strconv"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/util/http/request"
)

type RouteClient struct {
	client *request.RESTClient

	namespace string
	name      string
}

func (rc *RouteClient) Create(ctx context.Context, opts *rv1.RouteManifest) (*views.Route, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	var s *views.Route
	var e *errors.Http

	err = rc.client.Post(fmt.Sprintf("/namespace/%s/route", rc.namespace)).
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

func (rc *RouteClient) List(ctx context.Context) (*views.RouteList, error) {

	var s *views.RouteList
	var e *errors.Http

	err := rc.client.Get(fmt.Sprintf("/namespace/%s/route", rc.namespace)).
		AddHeader("Content-Type", "application/json").
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	if s == nil {
		list := make(views.RouteList, 0)
		s = &list
	}

	return s, nil
}

func (rc *RouteClient) Get(ctx context.Context) (*views.Route, error) {

	var s *views.Route
	var e *errors.Http

	err := rc.client.Get(fmt.Sprintf("/namespace/%s/route/%s", rc.namespace, rc.name)).
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

func (rc *RouteClient) Update(ctx context.Context, opts *rv1.RouteManifest) (*views.Route, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	var s *views.Route
	var e *errors.Http

	err = rc.client.Put(fmt.Sprintf("/namespace/%s/route/%s", rc.namespace, rc.name)).
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

func (rc *RouteClient) Remove(ctx context.Context, opts *rv1.RouteRemoveOptions) error {

	var e *errors.Http

	req := rc.client.Delete(fmt.Sprintf("/namespace/%s/route/%s", rc.namespace, rc.name)).
		AddHeader("Content-Type", "application/json")

	if opts != nil {
		if opts.Force {
			req.Param("force", strconv.FormatBool(opts.Force))
		}
	}

	if err := req.JSON(nil, &e); err != nil {
		return err
	}
	if e != nil {
		return errors.New(e.Message)
	}

	return nil
}

func newRouteClient(client *request.RESTClient, namespace, name string) *RouteClient {
	return &RouteClient{client: client, namespace: namespace, name: name}
}
