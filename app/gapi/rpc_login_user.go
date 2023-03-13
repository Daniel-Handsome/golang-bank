package gapi

import (
	"context"
	"database/sql"

	db "github.com/daniel/master-golang/db/sqlc"
	"github.com/daniel/master-golang/pb"
	"github.com/daniel/master-golang/utils"
	"github.com/daniel/master-golang/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	violations := validateLoginRequest(req)
	if violations != nil {
		return nil, invalidArgError(violations)
	}

	user, err := s.store.GetUser(ctx, req.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "fail to find user %v", req.Name)
		}
		return nil, status.Errorf(codes.Internal, "fail to get user %v", req.Name)
	}

	err = utils.CheckPassword(req.Password, user.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "password is unauth %v", req.Name)
	}

	// create token
	token, payload, err := s.tokenMaker.CreateToken(user.Name, s.config.Access_token_duration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create token error by internal %v", req.Name)
	}

	// refresh token
	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(user.Name, s.config.Access_token_duration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "refresh token create error %v", req.Name)
	}

	mtda := s.extractMetadata(ctx)
	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID : refreshPayload.ID,
		Username : user.Name,
		RefreshToken: refreshToken,
		UserAgent : mtda.UserAgent,
		ClientID : mtda.ClientIP,
		Isblacked: false,
		ExpiresAt: refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "session create error %v", req.Name)
	}

	response := &pb.LoginUserResponse{
		SessionId : session.ID.String(),
		AccessToken: token,
		AccessTokenAt: timestamppb.New(payload.ExpiredAt),
		RefreshToken: refreshToken,
		RefreshTokenAt: timestamppb.New(refreshPayload.ExpiredAt),
		User: converUser(user),
	}

	return response, nil 
}

func validateLoginRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidateName(req.GetName()); err != nil {
		return append(violations, fieldViolation("name", err))
	}

	if err := validation.ValidatePassword(req.GetName()); err != nil {
		return append(violations, fieldViolation("password", err))
	}

	return
}