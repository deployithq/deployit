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

package runtime

import (
	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/pkg/storage"
)

type Runtime struct {
	routes  *mux.Router
	storage storage.IStorage
}

func New(storage storage.IStorage) *Runtime {

	runtime := new(Runtime)
	runtime.storage = storage

	return runtime
}
