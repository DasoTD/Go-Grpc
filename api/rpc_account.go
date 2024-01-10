package api
import (
	"context"
	"errors"

	db "github.com/dasotd/go_grpc/db/sqlc"
	"github.com/dasotd/go_grpc/pb"
	"github.com/dasotd/go_grpc/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func(server *Server) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.GetUsername(),
			FirstName:       req.GetFirstname(),
			LastName :		req.GetLastname(),
			Email:          req.GetEmail(),
		},
	}

	rsp := &pb.CreateAccountResponse{
		Account: arg,
	}
	return rsp, nil
	
}

func(server *Server) GetAccount(ctx context.Context, req *pb.GetAccountRequest)(*pb.GetAccountResponse, error){
	violations := validateGetUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	account, err := server.grpc.GetUser(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}

	rsp := &pb.GetAccountResponse{
		Account: account,
	}
	return rsp, nil

}

func validateCreateUserRequest(req *pb.CreateAccountRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("userName", err))
	}


	if err := val.ValidateFirstName(req.GetFirstname()); err != nil {
		violations = append(violations, fieldViolation("firstName", err))
	}
	if err := val.ValidateLastName(req.GetLastname()); err != nil {
		violations = append(violations, fieldViolation("lastName", err))
	}

	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}

func fieldViolation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
	}
}

func invalidArgumentError(violations []*errdetails.BadRequest_FieldViolation) error {
	badRequest := &errdetails.BadRequest{FieldViolations: violations}
	statusInvalid := status.New(codes.InvalidArgument, "invalid parameters")

	statusDetails, err := statusInvalid.WithDetails(badRequest)
	if err != nil {
		return statusInvalid.Err()
	}

	return statusDetails.Err()
}

func unauthenticatedError(err error) error {
	return status.Errorf(codes.Unauthenticated, "unauthorized: %s", err)
}

func validateGetUserRequest(req *pb.GetAccountRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}

	return violations
}