package gapi

import (
	"context"

	db "github.com/daniel/master-golang/db/sqlc"
	"github.com/daniel/master-golang/pb"
	"github.com/daniel/master-golang/utils"
	"github.com/daniel/master-golang/validation"

	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserReponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgError(violations)
	}

	hashPassword, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password %s", err)
	}

	arg := db.CreateUserParams{
		Name : req.GetName(),    
		Password: hashPassword,
		FullName : req.GetFullName(),
		Email : req.GetEmail(),
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation" :
					return nil, status.Errorf(codes.AlreadyExists, "user name exists : %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user %s", err)
	}


	response := &pb.UserReponse{
		User: converUser(user),
	}

	return response, nil
}


func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidateName(req.GetName()); err != nil {
		return append(violations, fieldViolation("name", err))
	}

	if err := validation.ValidateFullName(req.GetName()); err != nil {
		return append(violations, fieldViolation("full_name", err))
	}

	if err := validation.ValidatePassword(req.GetName()); err != nil {
		return append(violations, fieldViolation("password", err))
	}

	if err := validation.ValidateEmail(req.GetName()); err != nil {
		return append(violations, fieldViolation("email", err))
	}

	return
}