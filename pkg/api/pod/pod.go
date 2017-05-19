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

package pod

import (
	"context"
	"fmt"
	ctx "github.com/lastbackend/lastbackend/pkg/api/context"
	"github.com/lastbackend/lastbackend/pkg/api/namespace"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	h "github.com/lastbackend/lastbackend/pkg/util/http"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

type pod struct {
	Context   context.Context
	Namespace types.Meta
}

func (p *pod) Set(pod types.Pod) error {
	var (
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)

	log.Debugf("update pod %s state: %s", pod.Meta.Name, pod.State.State)

	svc, err := storage.Service().GetByPodName(p.Context, pod.Meta.Name)
	if err != nil {
		log.Errorf("Error: get service by pod name %s from db: %s", pod.Meta.Name, err.Error())
		if err.Error() == store.ErrKeyNotFound {
			log.Debugf("Pod %s not found", pod.Meta.Name)
			return nil
		}
		return err
	}

	ns := namespace.New(p.Context)
	item, err := ns.Get(svc.Meta.Namespace)
	if err != nil {
		log.Error("Error: find namespace by name", err.Error())
		return err
	}
	if item == nil {
		return err
	}

	pd, e := storage.Pod().GetByName(p.Context, item.Meta.Name, pod.Meta.Name)
	if e != nil {
		log.Errorf("Error: get pod from db: %s", e.Error())
		return err
	}

	pd.Containers = pod.Containers
	pd.Meta = pod.Meta
	pd.State = pod.State

	if pd.State.State == types.StateDestroyed {
		log.Debugf("Service: Set pods: remove deleted pod: %s", pd.Meta.Name)
		if err := storage.Pod().Remove(p.Context, item.Meta.Name, pd); err != nil {
			log.Errorf("Error: set pod to db: %s", err)
			return err
		}
		delete(svc.Pods, pd.Meta.Name)

		if len(svc.Pods) == 0 && svc.State.State == types.StateDestroy {
			if err = storage.Service().Remove(p.Context, svc); err != nil {
				log.Errorf("Error: remove service from db: %s", err)
				return err
			}
		}

		return nil
	}

	if err := storage.Pod().Update(p.Context, item.Meta.Name, pd); err != nil {
		log.Errorf("Error: set pod to db: %s", err)
		return err
	}

	return nil
}

func Logs(c context.Context, ns, pod, container string, stream io.Writer, done chan bool) error {

	const buffer_size = 1024

	var (
		log      = ctx.Get().GetLogger()
		storage  = ctx.Get().GetStorage()
		buffer   = make([]byte, buffer_size)
		doneChan = make(chan bool, 1)
	)

	log.Debug("Service: get service logs")

	svc, err := storage.Service().GetByPodName(c, pod)
	if err != nil {
		log.Errorf("Error: get service by pod name %s from db: %s", pod, err.Error())
		if err.Error() == store.ErrKeyNotFound {
			log.Debugf("Pod %s not found", pod)
			return nil
		}
		return err
	}

	_ns := namespace.New(c)
	item, err := _ns.Get(svc.Meta.Namespace)
	if err != nil {
		log.Error("Error: find namespace by name", err.Error())
		return err
	}
	if item == nil {
		return err
	}

	if ns != item.Meta.Name {
		log.Error("Error: access denied")
		return errors.New("access denied")
	}

	p, e := storage.Pod().GetByName(c, item.Meta.Name, pod)
	if e != nil {
		log.Errorf("Error: get pod from db: %s", e.Error())
		return err
	}

	var cnt string
	for c := range p.Containers {
		if c == container {
			cnt = container
			break
		}
	}
	if cnt == "" {
		log.Error("Error: access denied")
		return errors.New("access denied")
	}

	n, err := storage.Node().Get(c, p.Meta.Hostname)
	if err != nil {
		log.Error("Error: find namespace by name", err.Error())
		return err
	}

	uri := fmt.Sprintf("%s:%d", n.Meta.IP, n.Meta.Port)
	client, err := h.New(uri, &h.ReqOpts{TLS: false})
	if err != nil {
		return err
	}

	_, res, err := client.
		GET(fmt.Sprintf("/container/%s/logs", cnt)).Do()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-doneChan:
				res.Body.Close()
				return
			default:
				n, err := res.Body.Read(buffer)
				if err != nil {
					log.Errorf("Error read bytes from stream %s", err)
					res.Body.Close()
					return
				}

				_, err = func(p []byte) (n int, err error) {
					n, err = stream.Write(p)
					if err != nil {
						log.Errorf("Error write bytes from stream %s", err)
						return n, err
					}
					if f, ok := stream.(http.Flusher); ok {
						f.Flush()
					}
					return n, nil
				}(buffer[0:n])
				if err != nil {
					log.Errorf("Error written to stream %s", err)
					return
				}

				for i := 0; i < n; i++ {
					buffer[i] = 0
				}
			}
		}
	}()

	<-done

	close(doneChan)

	return nil
}

func New(ctx context.Context) *pod {
	return &pod{
		Context: ctx,
	}
}
