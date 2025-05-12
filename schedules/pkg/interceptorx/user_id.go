package interceptorx

import (
	"context"
	grpcGen "github.com/tummerk/golang/schedules/internal/server/generated/grpc"
	"github.com/tummerk/golang/schedules/pkg/contextx"
	"github.com/tummerk/golang/schedules/pkg/utils"
	"google.golang.org/grpc"
	"strconv"
)

func UserIDInterceptor(Key []byte) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		var userID string

		switch r := req.(type) {
		case *grpcGen.UserID:
			if r.UserID != 0 {
				userID = strconv.Itoa(int(r.UserID))
			}
		case *grpcGen.GetScheduleRequest:

			if r.UserID != 0 {
				userID = strconv.Itoa(int(r.UserID))
			}
		case *grpcGen.CreateScheduleRequest:
			if r.UserId != 0 {
				userID = strconv.FormatInt(r.UserId, 10)
			}
		default:

			return handler(ctx, req)
		}

		if userID == "" {
			return handler(ctx, req)
		}
		maskedUserID, _ := utils.Encrypt(userID, Key)

		newCtx := contextx.WithUserID(ctx, contextx.UserID(userID))
		newCtx = contextx.WithMaskUserID(newCtx, contextx.MaskUserID(maskedUserID))

		return handler(newCtx, req)
	}
}
