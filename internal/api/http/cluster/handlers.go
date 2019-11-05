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

package cluster

import (
	"github.com/lastbackend/lastbackend/internal/api/envs"
	"github.com/lastbackend/lastbackend/internal/api/types/v1"
	"github.com/lastbackend/lastbackend/internal/pkg/model"
	"net/http"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/tools/log"
)

const (
	logLevel  = 2
	logPrefix = "api:handler:cluster"
)

func ClusterInfoH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /cluster cluster clusterInfo
	//
	// Shows an info about current cluster
	//
	// ---
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: Cluster response
	//     schema:
	//       "$ref": "#/definitions/views_cluster"
	//   '500':
	//     description: Internal server error

	log.V(logLevel).Debugf("%s:info:> get cluster", logPrefix)

	var clm = model.NewClusterModel(r.Context(), envs.Get().GetStorage())

	cl, err := clm.Get()
	if err != nil {
		log.V(logLevel).Errorf("%s:info:> get cluster err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	if cl == nil {
		cl = new(types.Cluster)
		cl.Meta.SetDefault()
		cl.Meta.Name, cl.Meta.Description = envs.Get().GetClusterInfo()
		cl.Meta.SelfLink = *types.NewClusterSelfLink(cl.Meta.Name)
	}

	response, err := v1.View().Cluster().New(cl).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("%s:info:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("%s:info:> write response err: %s", logPrefix, err.Error())
		return
	}
}
