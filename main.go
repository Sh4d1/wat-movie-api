package main

import (
	"errors"
	"log"
	"os"

	"golang.org/x/net/context"

	pb "github.com/Sh4d1/wat-movie-api/proto/movie-api"
	userService "github.com/Sh4d1/wat-user-service/proto/user"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
	k8s "github.com/micro/kubernetes/go/micro"
)

var (
	srv micro.Service
)

func main() {

	if os.Getenv("DEV") == "true" {
		srv = micro.NewService(
			micro.Name("wat.movie.api"),
			micro.WrapHandler(AuthWrapper),
		)
	} else {
		srv = k8s.NewService(
			micro.Name("wat.movie.api"),
			micro.WrapHandler(AuthWrapper),
		)
	}
	srv.Init()

	pb.RegisterMovieAPIHandler(srv.Server(), &service{})

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}

func AuthWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, resp interface{}) error {
		if os.Getenv("DISABLE_AUTH") == "true" {
			return fn(ctx, req, resp)
		}
		meta, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.New("no auth meta-data found in request")
		}
		token := meta["Token"]
		log.Println("Authenticating with token: ", token)
		authClient := userService.NewUserServiceClient("wat.user", srv.Client())
		_, err := authClient.ValidateToken(ctx, &userService.Token{
			Token: token,
		})
		if err != nil {
			return err
		}
		err = fn(ctx, req, resp)
		return err
	}
}
