package auth

import (
	"sync"

	grpc "google.golang.org/grpc"

	authv1 "github.com/taehoio/idl/gen/go/taehoio/idl/services/auth/v1"
)

//go:generate mockgen -package auth -destination ./client_mock.go -mock_names AuthServuceClient=MockAuthServiceClient github.com/taehoio/idl/gen/go/taehoio/idl/services/auth/v1 AuthServiceClient
const serviceConfig = `{"loadBalancingPolicy":"round_robin"}`

var (
	once sync.Once
	cli  authv1.AuthServiceClient

	_ authv1.AuthServiceClient = (*MockAuthServiceClient)(nil)
)

func GetAuthServiceClient(serviceHost, caller string) authv1.AuthServiceClient {
	once.Do(func() {
		conn, _ := grpc.Dial(
			serviceHost,
			grpc.WithDefaultServiceConfig(serviceConfig),
		)

		cli = authv1.NewAuthServiceClient(conn)
	})

	return cli
}
