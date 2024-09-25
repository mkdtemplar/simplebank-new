package gapi

import (
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	db "github.com/mkdtemplar/simplebank-new/db/sqlc"
	"github.com/mkdtemplar/simplebank-new/pb"
	"github.com/mkdtemplar/simplebank-new/util"
	"github.com/mkdtemplar/simplebank-new/validation"
	"github.com/mkdtemplar/simplebank-new/worker"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
)

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidateUserName(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := validation.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := validation.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}

	if err := validation.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, fmt.Errorf(codes.Internal.String(), "failed to hash password %v", err)
	}
	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.GetUsername(),
			HashedPassword: hashedPassword,
			FullName:       req.GetFullName(),
			Email:          req.GetEmail(),
		},
		AfterCreate: func(user db.User) error {
			taskPayload := &worker.PayloadSendVerifyEmail{Username: user.Username}
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}

			return s.taskDistributor.DistributeTaskVerifyEmail(ctx, taskPayload, opts...)
		},
	}

	txResult, err := s.store.CreateUserTx(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, fmt.Errorf(codes.AlreadyExists.String(), "username allready exist %v", err)
		}
		return nil, fmt.Errorf(codes.Internal.String(), "failed to create user %v", err)
	}
	

	rsp := &pb.CreateUserResponse{
		User: convertUser(txResult.User),
	}
	return rsp, nil
}
