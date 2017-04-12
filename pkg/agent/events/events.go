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

package events

import (
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/daemon/node/views/v1"
	"github.com/lastbackend/lastbackend/pkg/errors"
)

type Event struct {
}

func New() *Event {
	return new(Event)
}

func NewEvent() *types.Event {
	var event = new(types.Event)
	event.Meta.Hostname = "local"
	event.State.Allocated.Pods = 1
	return event
}

func (e *Event) Send(event *types.Event) error {

	var (
		er       = new(errors.Http)
		http     = context.Get().GetHttpClient()
		log      = context.Get().GetLogger()
		endpoint = "/node/event"
		spec     = new(v1.Spec)
	)

	log.Debugf("Send event request to: %s", endpoint)
	_, _, err := http.
		PUT(endpoint).
		AddHeader("Content-Type", "application/json").
		BodyJSON(event).
		Request(spec, er)
	if err != nil {
		log.Errorf("Send request error: %s", err.Error())
		return err
	}

	if er.Code == 401 {
		log.Error("401")
		return nil
	}

	if er.Code != 0 {
		log.Error(er.Code)
		return errors.New(er.Message)
	}

	s, _ := spec.ToJson()
	log.Debug(string(s))

	return nil
}