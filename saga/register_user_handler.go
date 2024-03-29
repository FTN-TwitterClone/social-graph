package saga

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/net/context"
	"os"
	"social-graph/repository"
	"social-graph/tracing"
)

type RegisterUserHandler struct {
	tracer trace.Tracer
	conn   *nats.Conn
	repo   repository.SocialGraphRepository
}

func NewRegisterUserHandler(tracer trace.Tracer, repo repository.SocialGraphRepository) (*RegisterUserHandler, error) {
	natsHost := os.Getenv("NATS_HOST")
	natsPort := os.Getenv("NATS_PORT")

	url := fmt.Sprintf("nats://%s:%s", natsHost, natsPort)

	connection, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	h := &RegisterUserHandler{
		tracer: tracer,
		conn:   connection,
		repo:   repo,
	}

	_, err = connection.Subscribe(REGISTER_COMMAND, h.handleCommand)
	if err != nil {
		return nil, err
	}

	return h, nil
}

func (h RegisterUserHandler) handleCommand(msg *nats.Msg) {
	remoteCtx, err := tracing.GetNATSParentContext(msg)
	if err != nil {

	}

	ctx, span := otel.Tracer("social-graph").Start(trace.ContextWithRemoteSpanContext(context.Background(), remoteCtx), "RegisterUserHandler.handleCommand")
	defer span.End()

	var c RegisterUserCommand

	err = json.Unmarshal(msg.Data, &c)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return
	}

	switch c.Command {
	case SaveSocialGraph:
		h.handleSaveSocialGraph(ctx, c.User)
	}
}

func (h RegisterUserHandler) handleSaveSocialGraph(ctx context.Context, user NewUser) {
	handlerCtx, span := h.tracer.Start(ctx, "RegisterUserHandler.handleSaveSocialGraph")
	defer span.End()

	//TODO: save user with all details
	var private bool

	if user.Role == "ROLE_USER" {
		private = true
	} else {
		private = false
	}

	err := h.repo.CreateNewUser(handlerCtx, user.Username, private)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())

		h.sendReply(handlerCtx, RegisterUserReply{
			Reply: SocialGraphFail,
			User:  user,
		})

		return
	}

	h.sendReply(handlerCtx, RegisterUserReply{
		Reply: SocialGraphSuccess,
		User:  user,
	})

}

func (h RegisterUserHandler) sendReply(ctx context.Context, r RegisterUserReply) {
	_, span := h.tracer.Start(ctx, "RegisterUserHandler.sendReply")
	defer span.End()

	headers := nats.Header{}
	headers.Set(tracing.TRACE_ID, span.SpanContext().TraceID().String())
	headers.Set(tracing.SPAN_ID, span.SpanContext().SpanID().String())

	data, err := json.Marshal(r)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return
	}

	msg := nats.Msg{
		Subject: REGISTER_REPLY,
		Header:  headers,
		Data:    data,
	}

	err = h.conn.PublishMsg(&msg)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
}
