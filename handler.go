package main

import (
	"log"

	pb "github.com/Sh4d1/wat-movie-api/proto/movie-api"
	"golang.org/x/net/context"
)

type service struct {
}

func (s *service) Get(ctx context.Context, req *pb.Movie, res *pb.Response) error {
	log.Println(req)
	//user, err := s.repo.Get(req.Id)
	//if err != nil {
	//	var err pb.Error
	//	log.Println("No user with id: ", req.Id)
	//	err.Code = 1
	//	err.Description = "User does not exist"
	//	res.Errors = append(res.Errors, &err)
	//	return nil
	//}
	//res.User = user
	return nil
}
