package gapi

import (
	db "github.com/daniel/master-golang/db/sqlc"
	"github.com/daniel/master-golang/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func converUser(user db.User) *pb.User {
	return &pb.User{
		Uame: user.Name,
		FullName: user.FullName,
		Email: user.Email,
		PasswordChangeAt: timestamppb.New(user.PasswordChangeAt),
		CreateAt: timestamppb.New(user.CreateAt),
	}
}