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

package etcd3

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/lastbackend/lastbackend/pkg/util/serializer"
	s "github.com/lastbackend/lastbackend/pkg/storage/store"
	"path"
)

func New(client *clientv3.Client, codec serializer.Codec, prefix string) s.IStore {
	return newStore(client, true, codec, prefix)
}

func NewWithNoQuorumRead(client *clientv3.Client, codec serializer.Codec, prefix string) s.IStore {
	return newStore(client, false, codec, prefix)
}

func newStore(client *clientv3.Client, quorumRead bool, codec serializer.Codec, prefix string) *store {
	var result = &store{
		client:     client,
		codec:      codec,
		pathPrefix: path.Join("/", prefix),
	}
	if !quorumRead {
		result.opts = append(result.opts, clientv3.WithSerializable())
	}
	return result
}
