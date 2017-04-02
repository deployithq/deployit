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

package project_test

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/client/cmd/project"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"github.com/lastbackend/lastbackend/pkg/client/storage"
	h "github.com/lastbackend/lastbackend/pkg/util/http"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreate(t *testing.T) {

	const (
		name        string = "project"
		description string = "project describe"
		token       string = "mocktoken"
	)

	var (
		err error
		ctx = context.Mock()

		prct = new(types.Project)
	)

	ctx.Storage, err = storage.Init()
	if err != nil {
		t.Error(err)
		return
	}
	defer (func() {
		err = ctx.Storage.Clear()
		if err != nil {
			t.Error(err)
			return
		}
		err = ctx.Storage.Close()
		if err != nil {
			t.Error(err)
			return
		}
	})()

	ctx.Token = token

	//------------------------------------------------------------------------------------------
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tk := r.Header.Get("Authorization")
		assert.NotEmpty(t, tk, "token should be not empty")
		assert.Equal(t, tk, "Bearer "+token, "they should be equal")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
			return
		}

		var d = struct {
			Name        string `json:"name,omitempty"`
			Description string `json:"description,omitempty"`
		}{}

		err = json.Unmarshal(body, &d)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, d.Name, name, "they should be equal")
		assert.Equal(t, d.Description, description, "they should be equal")

		w.WriteHeader(200)
		_, err = w.Write([]byte(`{"id":"mock", "name":"` + name + `", "description":"` + description + `"}`))
		if err != nil {
			t.Error(err)
			return
		}
	}))
	defer server.Close()

	ctx.HTTP = h.New(server.URL)
	//------------------------------------------------------------------------------------------

	err = project.Create(name, description)
	if err != nil {
		t.Error(err)
		return
	}

	err = ctx.Storage.Get("project", prct)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, prct.Name, name)
}
