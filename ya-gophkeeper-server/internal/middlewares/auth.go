package middleware

import (
	"context"
	"net/http"
	"strings"

	"yandex-gophkeeper-server/pkg/jwt"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthMiddleware3(jwtKey string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			if authorization == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			token, ok := strings.CutPrefix(authorization, "Bearer ")
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			userID, err := jwt.VerifyJWT([]byte(jwtKey), token)
			if err != nil {
				logrus.Error(err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			r.Header.Add("X-User-ID", userID)
			h.ServeHTTP(w, r)
		})
	}
}

func AuthGRPCMiddleware(jwtKey string) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logrus.Printf("[gRPC INFO] %s", info.FullMethod)
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			logrus.Info("context not have authhorization")
		}
		if len(md.Get("Authorization")) != 1 {
			err := status.Error(codes.Unauthenticated, "len(auth) not 1")
			logrus.Error(err)
			return nil, err
		}
		authorization := md.Get("Authorization")[0]
		logrus.Info(authorization)
		if authorization == "" {
			err := status.Error(codes.Unauthenticated, "authorization is empty")
			logrus.Error(err)
			return nil, err
		}
		token, ok := strings.CutPrefix(authorization, "Bearer ")
		if !ok {
			err := status.Error(codes.Unauthenticated, "Bearer is empty")
			logrus.Error(err)
			return nil, err
		}
		userID, err := jwt.VerifyJWT([]byte(jwtKey), token)
		if err != nil {
			err := status.Error(codes.Unauthenticated, "jwt not correct")
			logrus.Error(err)
			return nil, err
		}
		ctx = context.WithValue(ctx, "X-User-ID", userID)
		return handler(ctx, req)
	}
}
