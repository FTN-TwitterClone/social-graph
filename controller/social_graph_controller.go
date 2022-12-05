package controller

import (
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/codes"
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

	toUsername := mux.Vars(req)["username"]
	authUser := ctx.Value("authUser").(model.AuthUser)

	if authUser.Username == toUsername {
		http.Error(w, "Cant follow yourself", 400)
		return
	}
	err := sgc.socialGraphService.CreateFollow(ctx, authUser.Username, toUsername)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return
	}

}
func (sgc *SocialGraphController) RemoveFollow(w http.ResponseWriter, req *http.Request) {
	ctx, span := sgc.tracer.Start(req.Context(), "SocialGraphController.RemoveFollow")
	defer span.End()
	toUsername := mux.Vars(req)["username"]
	authUser := ctx.Value("authUser").(model.AuthUser)
	err := sgc.socialGraphService.RemoveFollow(ctx, authUser.Username, toUsername)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return
	}
}
func (sgc *SocialGraphController) GetFollowing(w http.ResponseWriter, req *http.Request) {
	ctx, span := sgc.tracer.Start(req.Context(), "SocialGraphController.GetFollowing")
	defer span.End()
	username := mux.Vars(req)["username"]
	users, err := sgc.socialGraphService.GetFollowing(ctx, username)
	if err != nil {
		return
	}
	err = json.EncodeJson(w, users)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return
	}
}
func (sgc *SocialGraphController) GetNumberOfFollowing(w http.ResponseWriter, req *http.Request) {
	ctx, span := sgc.tracer.Start(req.Context(), "SocialGraphController.GetNumberOfFollowing")
	defer span.End()
	username := mux.Vars(req)["username"]
	users, err := sgc.socialGraphService.GetFollowing(ctx, username)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return
	}
	err = json.EncodeJson(w, len(users))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return
	}
}

func (sgc *SocialGraphController) GetFollowers(w http.ResponseWriter, req *http.Request) {
	ctx, span := sgc.tracer.Start(req.Context(), "SocialGraphController.GetFollowers")
	defer span.End()
	username := mux.Vars(req)["username"]
	users, err := sgc.socialGraphService.GetFollowers(ctx, username)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return
	}
	err = json.EncodeJson(w, users)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return
	}
}
func (sgc *SocialGraphController) GetNumberOfFollowers(w http.ResponseWriter, req *http.Request) {
	ctx, span := sgc.tracer.Start(req.Context(), "SocialGraphController.GetNumberOfFollowers")
	defer span.End()
	username := mux.Vars(req)["username"]
	users, err := sgc.socialGraphService.GetFollowers(ctx, username)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return
	}
	err = json.EncodeJson(w, len(users))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return
	}
}

func (sgc *SocialGraphController) CheckIfFollowExists(w http.ResponseWriter, req *http.Request) {
	ctx, span := sgc.tracer.Start(req.Context(), "SocialGraphController.CheckIfFollowExists")
	defer span.End()
	authUser := ctx.Value("authUser").(model.AuthUser)
	to := mux.Vars(req)["username"]
	exists, err := sgc.socialGraphService.CheckIfFollowExists(ctx, authUser.Username, to)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return
	}
	err = json.EncodeJson(w, exists)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return
	}
}

func (sgc *SocialGraphController) AcceptRejectFollowRequest(w http.ResponseWriter, req *http.Request) {
	ctx, span := sgc.tracer.Start(req.Context(), "SocialGraphController.AcceptRejectFollowRequest")
	defer span.End()
	from := mux.Vars(req)["username"]
	authUser := ctx.Value("authUser").(model.AuthUser)
	approved, _ := json.DecodeJson[model.Approved](req.Body)
	err := sgc.socialGraphService.AcceptRejectFollowRequest(ctx, from, authUser.Username, approved.Approved)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return
	}
}

func (sgc *SocialGraphController) CheckIfFollowRequestExists(w http.ResponseWriter, req *http.Request) {
	ctx, span := sgc.tracer.Start(req.Context(), "SocialGraphController.CheckIfFollowRequestExists")
	defer span.End()
	authUser := ctx.Value("authUser").(model.AuthUser)
	from := mux.Vars(req)["username"]
	exists, err := sgc.socialGraphService.CheckIfFollowRequestExists(ctx, from, authUser.Username)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return
	}
	err = json.EncodeJson(w, exists)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return
	}
}

func (sgc *SocialGraphController) GetAllFollowRequests(w http.ResponseWriter, req *http.Request) {
	ctx, span := sgc.tracer.Start(req.Context(), "SocialGraphController.GetAllFollowRequests")
	defer span.End()
	authUser := ctx.Value("authUser").(model.AuthUser)
	users, err := sgc.socialGraphService.GetAllFollowRequests(ctx, authUser.Username)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return
	}
	err = json.EncodeJson(w, users)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return
	}
}
func (sgc *SocialGraphController) GetRecommendationsProfile(w http.ResponseWriter, req *http.Request) {
	ctx, span := sgc.tracer.Start(req.Context(), "SocialGraphController.GetRecommendationsProfile")
	defer span.End()
	authUser := ctx.Value("authUser").(model.AuthUser)
	users, err := sgc.socialGraphService.GetRecommendationsProfile(ctx, authUser.Username)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return
	}
	err = json.EncodeJson(w, users)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return
	}
}
