package middleware

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"yandex-gophkeeper-server/pkg/jwt"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"yandex-gophkeeper-server/internal/config"

	storageGRPCController "yandex-gophkeeper-server/internal/storage/delivery/grpc"
	"yandex-gophkeeper-server/internal/storage/entity"
	"yandex-gophkeeper-server/internal/storage/mocks"

	pb "github.com/havilcorp/yandex-gophkeeper-proto/save"
)

func TestJWTAuthMiddleware(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			t.Error(err)
		}
		w.WriteHeader(http.StatusOK)
	})

	jwtGood, err := jwt.GenerateJWT([]byte("jwtkeygood"), "1")
	if err != nil {
		t.Error(err)
	}

	jwtBad, err := jwt.GenerateJWT([]byte("jwtkeybad"), "1")
	if err != nil {
		t.Error(err)
	}

	type args struct {
		authorization string
	}
	tests := []struct {
		name       string
		args       args
		statusCode int
	}{
		{
			name: "good",
			args: args{
				authorization: fmt.Sprintf("Bearer %s", jwtGood),
			},
			statusCode: 200,
		},
		{
			name: "empty authorization",
			args: args{
				authorization: "",
			},
			statusCode: 401,
		},
		{
			name: "not have prefix bearer",
			args: args{
				authorization: "token",
			},
			statusCode: 401,
		},
		{
			name: "bad jwt",
			args: args{
				authorization: fmt.Sprintf("Bearer %s", jwtBad),
			},
			statusCode: 401,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/ping", nil)
			r.Header.Add("Authorization", tt.args.authorization)
			rw := httptest.NewRecorder()
			lm := JWTAuthMiddleware("jwtkeygood")(testHandler)
			lm.ServeHTTP(rw, r)
			res := rw.Result()
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			assert.Equal(t, res.StatusCode, tt.statusCode)
		})
	}
}

func TestAuthGRPCMiddleware(t *testing.T) {
	conf := config.Config{
		AddressGRPC: ":8081",
		JWTKey:      "jwt",
	}
	grpcListener, err := net.Listen("tcp", conf.AddressGRPC)
	if err != nil {
		t.Error(err)
	}
	serverGRPC := grpc.NewServer(grpc.ChainUnaryInterceptor(AuthGRPCMiddleware(conf.JWTKey)))
	uc := mocks.NewUseCase(t)
	uc.On("GetAll", 1).Return(&[]entity.Item{{
		ID:     1,
		UserId: 1,
		Data:   []byte(""),
		Meta:   "",
	}}, nil)
	pb.RegisterSaveServer(serverGRPC, storageGRPCController.NewHandler(&conf, uc))

	go func() {
		logrus.Printf("Сервер gRPC начал работу по адресу: %s\n", conf.AddressGRPC)
		if err := serverGRPC.Serve(grpcListener); err != nil {
			logrus.Error(err)
		}
		logrus.Printf("Сервер gRPC прекратил работу")
	}()

	connectGRPC, err := grpc.NewClient(conf.AddressGRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Error(err)
	}
	defer connectGRPC.Close()
	clientGRPC := pb.NewSaveClient(connectGRPC)

	jwtToken, err := jwt.GenerateJWT([]byte("jwt"), "1")
	if err != nil {
		t.Error(err)
	}
	header := metadata.New(map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", jwtToken),
	})
	ctx := metadata.NewOutgoingContext(context.Background(), header)
	clientGRPC.GetAll(ctx, &pb.GetAllRequest{})

	time.Sleep(time.Second * 2)

	if err := grpcListener.Close(); err != nil {
		logrus.Error(err)
	}
}
