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

package routes

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/errors"
	"net/http"
	"time"
)

func BuildListH(w http.ResponseWriter, _ *http.Request) {

	var (
		ctx = context.Get()
	)

	ctx.Log.Debug("Get boold list handler")

	builds, err := ctx.Storage.Build().ListByImage("", "")
	if err != nil {
		ctx.Log.Error(err)
		errors.HTTP.InternalServerError(w)
		return
	}

	buf, er := json.Marshal(builds)
	if er != nil {
		ctx.Log.Error(er.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, er = w.Write(buf)
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}

func BuildCreateH(w http.ResponseWriter, _ *http.Request) {

	var (
		ctx = context.Get()
	)

	ctx.Log.Debug("Get create build handler")

	b := new(types.Build)
	b.Created = time.Now()
	b.Updated = time.Now()

	build, err := ctx.Storage.Build().Insert(b)
	if err != nil {
		ctx.Log.Error(err)
		errors.HTTP.InternalServerError(w)
		return
	}

	buf, er := json.Marshal(build)
	if er != nil {
		ctx.Log.Error(er.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, er = w.Write(buf)
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}
