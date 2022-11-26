package controller

import (
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"social-graph/controller/json"
	"social-graph/model"
	"social-graph/service"
)

type SocialGraphController struct {
	socialGraphService service.SocialGraphService
	tracer             trace.Tracer
}

func NewSocialGraphController(socialGraphService *service.SocialGraphService, tracer trace.Tracer) *SocialGraphController {
	return &SocialGraphController{
		*socialGraphService,
		tracer,
	}
}

func (sgc *SocialGraphController) CreateFollow(w http.ResponseWriter, req *http.Request) {
	ctx, span := sgc.tracer.Start(req.Context(), "SocialGraphController.CreateFollow")
	defer span.End()

	f, _ := json.DecodeJson[model.Follows](req.Body)
	err := sgc.socialGraphService.CreateFollow(ctx, &f)
	if err != nil {
		return
	}
}
func (sgc *SocialGraphController) RemoveFollow(w http.ResponseWriter, req *http.Request) {
	f, _ := json.DecodeJson[model.Follows](req.Body)
	err := sgc.socialGraphService.RemoveFollow(&f)
	if err != nil {
		return
	}
}
func (sgc *SocialGraphController) GetFollowing(w http.ResponseWriter, req *http.Request) {
	username := mux.Vars(req)["username"]
	users, err := sgc.socialGraphService.GetFollowing(username)
	if err != nil {
		return
	}
	err = json.EncodeJson(w, users)
	if err != nil {
		return
	}
}
func (sgc *SocialGraphController) GetFollowers(w http.ResponseWriter, req *http.Request) {
	username := mux.Vars(req)["username"]
	users, err := sgc.socialGraphService.GetFollowers(username)
	if err != nil {
		return
	}
	err = json.EncodeJson(w, users)
	if err != nil {
		return
	}
}
