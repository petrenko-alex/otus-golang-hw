package log

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

const unknown = "UNKNOWN"

type Interceptor struct {
	logger logger.Logger
}

func NewInterceptor(logger logger.Logger) *Interceptor {
	return &Interceptor{
		logger: logger,
	}
}

func (h *Interceptor) GetInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		end := time.Since(start)

		ip := unknown
		peerInfo, ok := peer.FromContext(ctx)
		if ok {
			ip = peerInfo.Addr.String()
		}

		userAgent := unknown
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			userAgent = md.Get("user-agent")[0]
		}

		statusCode := codes.Unknown
		if st, ok := status.FromError(err); ok {
			statusCode = st.Code()
		}

		logJSON, marshalErr := json.Marshal(
			struct {
				IP        string
				Datetime  string
				Method    string
				Status    string
				Time      string
				UserAgent string
			}{
				IP:        ip,
				Datetime:  time.Now().Format(time.RFC822),
				Method:    info.FullMethod,
				Status:    strconv.Itoa(int(statusCode)),
				Time:      end.String(),
				UserAgent: userAgent,
			},
		)
		if err != nil {
			h.logger.Error(marshalErr.Error())
		}

		h.logger.Info(string(logJSON))

		return resp, err
	}
}
