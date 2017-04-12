package routes

import (
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/daemon/image"
	"github.com/lastbackend/lastbackend/pkg/daemon/namespace"
	"github.com/lastbackend/lastbackend/pkg/daemon/service"
	"github.com/lastbackend/lastbackend/pkg/daemon/service/routes/request"
	"github.com/lastbackend/lastbackend/pkg/daemon/service/views/v1"
	"github.com/lastbackend/lastbackend/pkg/errors"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"net/http"
)

func ServiceListH(w http.ResponseWriter, r *http.Request) {

	var (
		err error
		log = context.Get().GetLogger()
		id  = utils.Vars(r)["namespace"]
	)

	log.Debug("List service handler")
	ns := namespace.New()
	item, err := ns.Get(r.Context(), id)
	if err != nil {
		log.Error("Error: find namespace by id", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		errors.New("namespace").NotFound().Http(w)
		return
	}

	s := service.New()
	items, err := s.List(r.Context(), item.Meta.ID)
	if err != nil {
		log.Error("Error: find service list by user", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewServiceList(items).ToJson()
	if err != nil {
		log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}
}

func ServiceInfoH(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		log = context.Get().GetLogger()

		nid = utils.Vars(r)["namespace"]
		sid = utils.Vars(r)["service"]
	)

	log.Debug("Get service handler")

	ns := namespace.New()
	item, err := ns.Get(r.Context(), nid)
	if err != nil {
		log.Error("Error: find namespace by id", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		errors.New("namespace").NotFound().Http(w)
		return
	}

	s := service.New()
	svc, err := s.Get(r.Context(), item.Meta.Name, sid)
	if err != nil {
		log.Error("Error: find service by id", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if svc == nil {
		errors.New("service").NotFound().Http(w)
		return
	}

	response, err := v1.NewService(svc).ToJson()
	if err != nil {
		log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}
}

func ServiceCreateH(w http.ResponseWriter, r *http.Request) {

	var (
		err error
		log = context.Get().GetLogger()
		nid = utils.Vars(r)["namespace"]
		sid = utils.Vars(r)["service"]
	)

	log.Debug("Create service handler")

	// request body struct
	rq := new(request.RequestServiceCreateS)
	if err := rq.DecodeAndValidate(r.Body); err != nil {
		log.Error("Error: validation incomming data", err)
		errors.New("Invalid incomming data").Unknown().Http(w)
		return
	}

	ns := namespace.New()
	item, err := ns.Get(r.Context(), nid)
	if err != nil {
		log.Error("Error: find namespace by id", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		errors.New("namespace").NotFound().Http(w)
		return
	}

	s := service.New()
	svc, err := s.Get(r.Context(), item.Meta.ID, sid)
	if err != nil {
		log.Error("Error: find service by id", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if svc != nil {
		errors.New("service").NotUnique("name").Http(w)
		return
	}

	// Load template from registry
	if rq.Template != "" {
		// TODO: Send request for get template config from registry
		// TODO: Set service source with types.SourceTemplateType type field
		// TODO: Patch template config if need
		// TODO: Template provision
	}

	// If you are not using a template, then create a standard configuration template
	//if tpl == nil {
	// TODO: Generate default template for service
	//return
	//}

	// Patch config if exists custom configurations
	if rq.Config != nil {
		// TODO: If have custom config, then need patch this config
	} else {
		rq.Config = types.ServiceConfig{}.GetDefault()
	}

	if rq.Source != nil {
		img, err := image.Create(r.Context(), rq.Registry, rq.Source)
		if err != nil {
			log.Error("Error: insert service to db", err)
			errors.HTTP.InternalServerError(w)
			return
		}
		rq.Config.Image = img.Meta.Name
	} else {
		rq.Config.Image = rq.Image
	}

	svc, err = s.Create(r.Context(), item.Meta.ID, rq)
	if err != nil {
		log.Error("Error: insert service to db", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewService(svc).ToJson()
	if err != nil {
		log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}
}

func ServiceUpdateH(w http.ResponseWriter, r *http.Request) {

	var (
		err error
		log = context.Get().GetLogger()
		nid = utils.Vars(r)["namespace"]
		sid = utils.Vars(r)["service"]
	)

	log.Debug("Update service handler")

	// request body struct
	rq := new(request.RequestServiceUpdateS)
	if err := rq.DecodeAndValidate(r.Body); err != nil {
		log.Error("Error: validation incomming data", err)
		errors.New("Invalid incomming data").Unknown().Http(w)
		return
	}

	ns := namespace.New()
	item, err := ns.Get(r.Context(), nid)
	if err != nil {
		log.Error("Error: find namespace by id", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		errors.New("namespace").NotFound().Http(w)
		return
	}

	s := service.New()
	svc, err := s.Get(r.Context(), item.Meta.ID, sid)
	if err != nil {
		log.Error("Error: find service by id", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if svc == nil {
		errors.New("service").NotFound().Http(w)
		return
	}

	if rq.Name != "" {
		svc.Meta.Name = rq.Name
	}

	if rq.Description != "" {
		svc.Meta.Description = rq.Description
	}

	if rq.Config != nil {
		if err := svc.Config.Update(rq.Config); err != nil {
			log.Error("Error: update service config", err.Error())
			errors.New("service").BadParameter("config", err)
			return
		}
	}

	if rq.Domains != nil {
		svc.Domains = *rq.Domains
	}

	svc, err = s.Update(r.Context(), nid, svc)
	if err != nil {
		log.Error("Error: insert service to db", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	// TODO: spec generate
	response, err := v1.NewService(svc).ToJson()
	if err != nil {
		log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}
}

func ServiceRemoveH(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		log = context.Get().GetLogger()

		nid = utils.Vars(r)["namespace"]
		sid = utils.Vars(r)["service"]
	)

	log.Info("Remove service")

	ns := namespace.New()
	item, err := ns.Get(r.Context(), nid)
	if err != nil {
		log.Error("Error: find namespace by id", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		errors.New("namespace").NotFound().Http(w)
		return
	}

	s := service.New()
	svc, err := s.Get(r.Context(), item.Meta.Name, sid)
	if err != nil {
		log.Error("Error: find service by id", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if svc == nil {
		errors.New("service").NotFound().Http(w)
		return
	}

	// Todo: remove all activity by service name

	if err := s.Remove(r.Context(), item.Meta.ID, svc); err != nil {
		log.Error("Error: remove service from db", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}
}

func ServiceActivityListH(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`[]`)); err != nil {
		context.Get().GetLogger().Error("Error: write response", err.Error())
		return
	}
}

func ServiceLogsH(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`[]`)); err != nil {
		context.Get().GetLogger().Error("Error: write response", err.Error())
		return
	}
}