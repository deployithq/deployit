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

package views

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/internal/pkg/types"
)

type DeploymentView struct{}

func (dv *DeploymentView) New(obj *types.Deployment) *Deployment {
	d := Deployment{}
	d.SetMeta(obj.Meta)
	d.SetStatus(obj.Status)
	d.SetSpec(obj.Spec)

	return &d
}

func (d *Deployment) SetMeta(obj types.DeploymentMeta) {
	meta := DeploymentMeta{}

	meta.Name = obj.Name
	meta.Description = obj.Description
	meta.Version = obj.Version
	meta.SelfLink = obj.SelfLink.String()
	meta.Namespace = obj.Namespace
	meta.Service = obj.Service
	meta.Endpoint = obj.Endpoint
	meta.Updated = obj.Updated
	meta.Created = obj.Created

	meta.Labels = make(map[string]string, 0)
	for k, v := range obj.Labels {
		meta.Labels[k] = v
	}

	d.Meta = meta
}

func (d *Deployment) SetStatus(obj types.DeploymentStatus) {
	d.Status = DeploymentStatusInfo{
		State:   obj.State,
		Message: obj.Message,
	}
}

func (d *Deployment) SetSpec(obj types.DeploymentSpec) {
	mv := new(ManifestView)
	var spec = DeploymentSpec{
		Replicas: obj.Replicas,
		Template: mv.NewManifestSpecTemplate(obj.Template),
		Selector: mv.NewManifestSpecSelector(obj.Selector),
	}

	d.Spec = spec
}

func (d *Deployment) ToJson() ([]byte, error) {
	return json.Marshal(d)
}

func (dv *DeploymentView) NewList(obj *types.DeploymentList) *DeploymentList {
	dl := make(DeploymentList, 0)
	for _, d := range obj.Items {
		dv := new(DeploymentView)
		dp := dv.New(d)
		dl = append(dl, dp)
	}
	return &dl
}

func (di *DeploymentList) ToJson() ([]byte, error) {
	return json.Marshal(di)
}
