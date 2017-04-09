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
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/satori/go.uuid"
	"golang.org/x/oauth2"
)

type IUtil interface {
	Key(ctx context.Context, pattern ...string) string
}

type IStorage interface {
	Activity() IActivity
	Build() IBuild
	Hook() IHook
	Image() IImage
	Project() IProject
	Service() IService
	Vendor() IVendor
	Volume() IVolume
}

type IActivity interface {
	Insert(ctx context.Context, activity *types.Activity) (*types.Activity, error)
	ListProjectActivity(ctx context.Context, id uuid.UUID) (*types.ActivityList, error)
	ListServiceActivity(ctx context.Context, id uuid.UUID) (*types.ActivityList, error)
	RemoveByProject(ctx context.Context, id uuid.UUID) error
	RemoveByService(ctx context.Context, id uuid.UUID) error
}

type IBuild interface {
	GetByID(ctx context.Context, imageID, id uuid.UUID) (*types.Build, error)
	ListByImage(ctx context.Context, id uuid.UUID) (*types.BuildList, error)
	Insert(ctx context.Context, imageID uuid.UUID, source *types.BuildSource) (*types.Build, error)
}

type IHook interface {
	GetByToken(ctx context.Context, token string) (*types.Hook, error)
	ListByImage(ctx context.Context, id uuid.UUID) (*types.HookList, error)
	ListByService(ctx context.Context, id uuid.UUID) (*types.HookList, error)
	Insert(ctx context.Context, hook *types.Hook) (*types.Hook, error)
	Remove(ctx context.Context, id uuid.UUID) error
	RemoveByService(ctx context.Context, id uuid.UUID) error
}

type IProject interface {
	GetByID(ctx context.Context, id uuid.UUID) (*types.Project, error)
	GetByName(ctx context.Context, name string) (*types.Project, error)
	List(ctx context.Context) (*types.ProjectList, error)
	Insert(ctx context.Context, name, description string) (*types.Project, error)
	Update(ctx context.Context, project *types.Project) (*types.Project, error)
	Remove(ctx context.Context, id uuid.UUID) error
}

type IService interface {
	GetByID(ctx context.Context, project, id uuid.UUID) (*types.Service, error)
	GetByName(ctx context.Context, project uuid.UUID, name string) (*types.Service, error)
	ListByProject(ctx context.Context, project uuid.UUID) (*types.ServiceList, error)
	Insert(ctx context.Context, project uuid.UUID, name, description string, config *types.ServiceConfig) (*types.Service, error)
	Update(ctx context.Context, project uuid.UUID, service *types.Service) (*types.Service, error)
	Remove(ctx context.Context, project, id uuid.UUID) error
	RemoveByProject(ctx context.Context, project uuid.UUID) error
}

type IImage interface {
	Get(ctx context.Context, name string) (*types.Image, error)
	Insert(ctx context.Context, name string, source *types.ImageSource) (*types.Image, error)
	Update(ctx context.Context, image *types.Image) (*types.Image, error)
}

type IVendor interface {
	Insert(ctx context.Context, owner, name, host, serviceID string, token *oauth2.Token) error
	Get(ctx context.Context, name string) (*types.Vendor, error)
	List(ctx context.Context) (map[string]*types.Vendor, error)
	Update(ctx context.Context, vendor *types.Vendor) error
	Remove(ctx context.Context, vendorName string) error
}

type IVolume interface {
	GetByToken(ctx context.Context, token string) (*types.Volume, error)
	ListByProject(ctx context.Context, project uuid.UUID) (*types.VolumeList, error)
	Insert(ctx context.Context, volume *types.Volume) (*types.Volume, error)
	Remove(ctx context.Context, id uuid.UUID) error
}
