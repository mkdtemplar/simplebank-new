package gapi

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	db "github.com/mkdtemplar/simplebank-new/db/sqlc"
	"github.com/mkdtemplar/simplebank-new/pb"
	"github.com/mkdtemplar/simplebank-new/util"
	"github.com/mkdtemplar/simplebank-new/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	var hashedPassword string
	var err error
	violations := validateUpdateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: sql.NullString{String: req.GetFullName(), Valid: req.FullName != nil},
		Email:    sql.NullString{String: req.GetEmail(), Valid: req.Email != nil},
	}

	if req.Password != nil {
		hashedPassword, err = util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, fmt.Errorf(codes.Internal.String(), "failed to hash password %v", err)
		}
	}

	arg.HashedPassword = sql.NullString{
		String: hashedPassword,
	}

	arg.PasswordChangetAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	user, err := s.store.UpdateUser(ctx, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, fmt.Errorf(codes.Internal.String(), "failed to update user %v", err)
	}

	rsp := &pb.UpdateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidateUserName(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if req.Password != nil {
		if err := validation.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}

	if req.FullName != nil {
		if err := validation.ValidateFullName(req.GetFullName()); err != nil {
			violations = append(violations, fieldViolation("full_name", err))
		}
	}

	if req.Email != nil {
		if err := validation.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}

	return violations
}
