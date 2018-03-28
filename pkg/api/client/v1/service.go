//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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

	"encoding/json"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/api/client/http"
	"github.com/lastbackend/lastbackend/pkg/api/client/interfaces"
	rv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	vv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"io"
	"net/url"
	"strconv"
)

type ServiceClient struct {
	interfaces.Service
	client    http.Interface
	namespace string
	name      string
}

func (s *ServiceClient) Deployment(name string) *DeploymentClient {
	return newDeploymentClient(s.client, s.namespace, s.name, name)
}

func (s *ServiceClient) Trigger(name string) *TriggerClient {
	return newTriggerClient(s.client, s.namespace, s.name, name)
}

func (s *ServiceClient) Create(ctx context.Context, opts *rv1.ServiceCreateOptions) (*vv1.Service, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	req := s.client.Post(fmt.Sprintf("/namespace/%s/service", s.namespace)).
		AddHeader("Content-Type", "application/json").
		Body(body).
		Do()

	if err := req.Error(); err != nil {
		return nil, err
	}

	buf, err := req.Raw()
	if err != nil {
		return nil, err
	}

	if code := req.StatusCode(); 200 > code || code > 299 {
		var e *errors.Http
		if err := json.Unmarshal(buf, &e); err != nil {
			return nil, err
		}
		return nil, errors.New(e.Message)
	}

	var ss *vv1.Service

	if err := json.Unmarshal(buf, &ss); err != nil {
		return nil, err
	}

	return ss, nil
}

func (s *ServiceClient) List(ctx context.Context) (*vv1.ServiceList, error) {

	req := s.client.Get(fmt.Sprintf("/namespace/%s/service", s.namespace)).
		AddHeader("Content-Type", "application/json").
		Do()

	buf, err := req.Raw()
	if err != nil {
		return nil, err
	}

	if code := req.StatusCode(); 200 > code || code > 299 {
		var e *errors.Http
		if err := json.Unmarshal(buf, &e); err != nil {
			return nil, err
		}
		return nil, errors.New(e.Message)
	}

	var sl *vv1.ServiceList

	if len(buf) == 0 {
		list := make(vv1.ServiceList, 0)
		return &list, nil
	}

	if err := json.Unmarshal(buf, &sl); err != nil {
		return nil, err
	}

	return sl, nil
}

func (s *ServiceClient) Get(ctx context.Context) (*vv1.Service, error) {

	req := s.client.Get(fmt.Sprintf("/namespace/%s/service/%s", s.namespace, s.name)).
		AddHeader("Content-Type", "application/json").
		Do()

	buf, err := req.Raw()
	if err != nil {
		return nil, err
	}

	if code := req.StatusCode(); 200 > code || code > 299 {
		var e *errors.Http
		if err := json.Unmarshal(buf, &e); err != nil {
			return nil, err
		}
		return nil, errors.New(e.Message)
	}

	var ss *vv1.Service

	if err := json.Unmarshal(buf, &ss); err != nil {
		return nil, err
	}

	return ss, nil
}

func (s *ServiceClient) Update(ctx context.Context, opts *rv1.ServiceUpdateOptions) (*vv1.Service, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	req := s.client.Put(fmt.Sprintf("/namespace/%s/service/%s", s.namespace, s.name)).
		AddHeader("Content-Type", "application/json").
		Body(body).
		Do()

	buf, err := req.Raw()
	if err != nil {
		return nil, err
	}

	if code := req.StatusCode(); 200 > code || code > 299 {
		var e *errors.Http
		if err := json.Unmarshal(buf, &e); err != nil {
			return nil, err
		}
		return nil, errors.New(e.Message)
	}

	var ss *vv1.Service

	if err := json.Unmarshal(buf, &ss); err != nil {
		return nil, err
	}

	return ss, nil
}

func (s *ServiceClient) Remove(ctx context.Context, opts *rv1.ServiceRemoveOptions) error {

	v := url.Values{}

	if opts != nil {
		if opts.Force {
			v.Set("force", strconv.FormatBool(opts.Force))
		}
	}

	qs := v.Encode()

	if len(qs) != 0 {
		qs = "?" + qs
	}

	req := s.client.Delete(fmt.Sprintf("/namespace/%s/service/%s%s", s.namespace, s.name, qs)).
		AddHeader("Content-Type", "application/json").
		Do()

	buf, err := req.Raw()
	if err != nil {
		return err
	}

	if code := req.StatusCode(); 200 > code || code > 299 {
		var e *errors.Http
		if err := json.Unmarshal(buf, &e); err != nil {
			return err
		}
		return errors.New(e.Message)
	}

	return nil
}

func (s *ServiceClient) Logs(ctx context.Context, opts *rv1.ServiceLogsOptions) (io.ReadCloser, error) {

	v := url.Values{}

	if opts != nil {
		v.Set("deployment", opts.Deployment)
		v.Set("pod", opts.Pod)
		v.Set("container", opts.Container)

		if opts.Follow {
			v.Set("follow", strconv.FormatBool(opts.Follow))
		}
	}

	qs := v.Encode()

	if len(qs) != 0 {
		qs = "?" + qs
	}

	return s.client.Get(fmt.Sprintf("/namespace/%s/service/%s/logs%s", s.namespace, s.name, qs)).
		AddHeader("Content-Type", "application/json").
		Stream()
}

func newServiceClient(client http.Interface, namespace, name string) *ServiceClient {
	return &ServiceClient{client: client, namespace: namespace, name: name}
}
