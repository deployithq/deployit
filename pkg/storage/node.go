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

package storage

import (
	"context"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"regexp"
	"time"
)

const nodeStorage = "node"

// Namespace Service type for interface in interfaces folder
type NodeStorage struct {
	INode
	util   IUtil
	Client func() (store.IStore, store.DestroyFunc, error)
}

func (s *NodeStorage) List(ctx context.Context) ([]*types.Node, error) {
	const filter = `\b(.+)` + nodeStorage + `\/(.+)\/(meta|state)\b`
	nodes := []*types.Node{}
	client, destroy, err := s.Client()
	if err != nil {
		return nodes, err
	}
	defer destroy()

	key := s.util.Key(ctx, nodeStorage)
	if err := client.List(ctx, key, filter, &nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}

func (s *NodeStorage) Get(ctx context.Context, hostname string) (*types.Node, error) {

	const filter = `\b.+` + nodeStorage + `\/.+\/(meta|state)\b`

	var (
		node = new(types.Node)
	)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	key := s.util.Key(ctx, nodeStorage, hostname)
	if err := client.Map(ctx, key, filter, node); err != nil {

		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}

		return nil, err
	}

	node.Spec.Pods = make(map[string]types.PodNodeSpec)
	keySpec := s.util.Key(ctx, nodeStorage, hostname, "spec", "pods")
	if err := client.Map(ctx, keySpec, "", node.Spec.Pods); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return node, nil
		}

		return nil, err
	}

	return node, nil
}

func (s *NodeStorage) Insert(ctx context.Context, node *types.Node) error {

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyMeta := s.util.Key(ctx, nodeStorage, node.Meta.Hostname, "meta")
	if err := tx.Create(keyMeta, &node.Meta, 0); err != nil {
		fmt.Println("meta", err.Error())
		return err
	}

	keyState := s.util.Key(ctx, nodeStorage, node.Meta.Hostname, "state")
	if err := tx.Create(keyState, &node.State, 0); err != nil {
		fmt.Println("meta", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("commit", err.Error())
		return err
	}

	return nil
}

func (s *NodeStorage) UpdateMeta(ctx context.Context, meta *types.NodeMeta) error {
	meta.Updated = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)
	keyMeta := s.util.Key(ctx, nodeStorage, meta.Hostname, "meta")
	if err := tx.Update(keyMeta, meta, 0); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}

func (s *NodeStorage) UpdateState(ctx context.Context, node *types.Node) error {
	node.Meta.Updated = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)
	keyMeta := s.util.Key(ctx, nodeStorage, node.Meta.Hostname, "meta")
	if err := tx.Update(keyMeta, &node.Meta, 0); err != nil {
		return err
	}

	keyState := s.util.Key(ctx, nodeStorage, node.Meta.Hostname, "state")
	if err := tx.Update(keyState, &node.Meta, 0); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *NodeStorage) InsertPod(ctx context.Context, meta *types.NodeMeta, pod *types.PodNodeSpec) error {
	meta.Updated = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)
	keyMeta := s.util.Key(ctx, nodeStorage, meta.Hostname, "meta")
	if err := tx.Update(keyMeta, meta, 0); err != nil {
		return err
	}

	keyPod := s.util.Key(ctx, nodeStorage, meta.Hostname, "pod", pod.Meta.ID)
	if err := tx.Create(keyPod, pod, 0); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *NodeStorage) UpdatePod(ctx context.Context, meta *types.NodeMeta, pod *types.PodNodeSpec) error {
	meta.Updated = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)
	keyMeta := s.util.Key(ctx, nodeStorage, meta.Hostname, "meta")
	if err := tx.Update(keyMeta, meta, 0); err != nil {
		return err
	}

	keyPod := s.util.Key(ctx, nodeStorage, meta.Hostname, "pod", pod.Meta.ID)
	if err := tx.Update(keyPod, pod, 0); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *NodeStorage) RemovePod(ctx context.Context, meta *types.NodeMeta, pod *types.PodNodeSpec) error {
	meta.Updated = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)
	keyMeta := s.util.Key(ctx, nodeStorage, meta.Hostname, "meta")
	if err := tx.Update(keyMeta, meta, 0); err != nil {
		return err
	}

	keyPod := s.util.Key(ctx, nodeStorage, meta.Hostname, "pod", pod.Meta.ID)
	tx.Delete(keyPod)

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *NodeStorage) Remove(ctx context.Context, node *types.Node) error {
	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)
	key := s.util.Key(ctx, nodeStorage, node.Meta.Hostname)
	tx.DeleteDir(key)

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *NodeStorage) Watch(ctx context.Context, service chan *types.Node) error {
	const filter = `\b.+` + nodeStorage + `\/(.+)\/alive\b`

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := s.util.Key(ctx, nodeStorage)
	cb := func(action, key string) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 2 {
			return
		}

	}

	client.Watch(ctx, key, filter, cb)
	return nil
}

func newNodeStorage(config store.Config, util IUtil) *NodeStorage {
	s := new(NodeStorage)
	s.util = util
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}