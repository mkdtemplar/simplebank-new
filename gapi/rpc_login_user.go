package gapi

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	db "github.com/mkdtemplar/simplebank-new/db/sqlc"
	"github.com/mkdtemplar/simplebank-new/pb"
	"github.com/mkdtemplar/simplebank-new/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := s.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf(codes.NotFound.String(), "user not found %v", req.GetUsername())
		}
		return nil, fmt.Errorf(codes.Internal.String(), "server error %v", err)
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, fmt.Errorf(codes.Unauthenticated.String(), "user not authenticated %v", err)
	}
	accessToken, accessTokenPayload, err := s.tokenMaker.CreateToken(user.Username, s.config.AccessTokenDuration)
	if err != nil {
		return nil, fmt.Errorf(codes.Internal.String(), "failed to create access token %v", err)
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(user.Username, s.config.RefreshTokenDuration)
	if err != nil {
		return nil, fmt.Errorf(codes.Internal.String(), "failed to create refresh token %v", err)
	}

	mtdt := s.extractMetadata(ctx)

	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIp,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, fmt.Errorf(codes.Internal.String(), "failed to create session %v", err)
	}

	rsp := &pb.LoginUserResponse{
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessTokenPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}
	return rsp, nil
}
