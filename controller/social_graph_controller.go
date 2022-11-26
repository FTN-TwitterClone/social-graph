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

	if f.From.Username == f.To.Username {
		http.Error(w, "Cant follow yourself", 400)
		return
	}
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
func (sgc *SocialGraphController) GetNumberOfFollowing(w http.ResponseWriter, req *http.Request) {
	username := mux.Vars(req)["username"]
	users, err := sgc.socialGraphService.GetFollowing(username)
	if err != nil {
		return
	}
	err = json.EncodeJson(w, len(users))
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
func (sgc *SocialGraphController) GetNumberOfFollowers(w http.ResponseWriter, req *http.Request) {
	username := mux.Vars(req)["username"]
	users, err := sgc.socialGraphService.GetFollowers(username)
	if err != nil {
		return
	}
	err = json.EncodeJson(w, len(users))
	if err != nil {
		return
	}
}

func (sgc *SocialGraphController) CheckIfFollowExists(w http.ResponseWriter, req *http.Request) {
	from := mux.Vars(req)["from"]
	to := mux.Vars(req)["to"]
	exists, err := sgc.socialGraphService.CheckIfFollowExists(from, to)
	if err != nil {
		return
	}
	err = json.EncodeJson(w, exists)
	if err != nil {
		return
	}
}

func (sgc *SocialGraphController) AcceptRejectFollowRequest(w http.ResponseWriter, req *http.Request) {
	from := mux.Vars(req)["from"]
	to := mux.Vars(req)["to"]
	approved, _ := json.DecodeJson[model.Approved](req.Body)
	err := sgc.socialGraphService.AcceptRejectFollowRequest(from, to, approved.Approved)
	if err != nil {
		return
	}
}
