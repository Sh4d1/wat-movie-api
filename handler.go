package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	pb "github.com/Sh4d1/wat-movie-api/proto/movie-api"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

type service struct {
}

type OmdbMovie struct {
	Actors     string `json:"Actors"`
	Awards     string `json:"Awards"`
	BoxOffice  string `json:"BoxOffice"`
	Country    string `json:"Country"`
	DVD        string `json:"DVD"`
	Director   string `json:"Director"`
	Genre      string `json:"Genre"`
	Language   string `json:"Language"`
	Metascore  string `json:"Metascore"`
	Plot       string `json:"Plot"`
	Poster     string `json:"Poster"`
	Production string `json:"Production"`
	Rated      string `json:"Rated"`
	Ratings    []struct {
		Source string `json:"Source"`
		Value  string `json:"Value"`
	} `json:"Ratings"`
	Released   string `json:"Released"`
	Response   string `json:"Response"`
	Runtime    string `json:"Runtime"`
	Title      string `json:"Title"`
	Type       string `json:"Type"`
	Website    string `json:"Website"`
	Writer     string `json:"Writer"`
	Year       string `json:"Year"`
	ImdbID     string `json:"imdbID"`
	ImdbRating string `json:"imdbRating"`
	ImdbVotes  string `json:"imdbVotes"`
}

type OmdbMovieList struct {
	Response string `json:"Response"`
	Search   []struct {
		Poster string `json:"Poster"`
		Title  string `json:"Title"`
		Type   string `json:"Type"`
		Year   string `json:"Year"`
		ImdbID string `json:"imdbID"`
	} `json:"Search"`
	TotalResults string `json:"totalResults"`
	Error        string `json:"Error"`
}

func (s *service) List(ctx context.Context, req *pb.ListRequest, res *pb.Response) error {
	apiKey := os.Getenv("OMDB_API_KEY")
	log.Println(apiKey)
	name := strings.Replace(url.QueryEscape(req.Name), " ", "+", -1)
	if req.Page == 0 {
		req.Page = 1
	}

	url := fmt.Sprintf("http://www.omdbapi.com/?s=%s&type=movie&page=%d&apikey=%s", name, req.Page, apiKey)

	get, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("GetReq: ", err)
		return err
	}
	client := ProxyAwareHttpClient()
	resp, err := client.Do(get)
	if err != nil {
		log.Printf("Do: ", err)
		return err
	}

	defer resp.Body.Close()

	var movies OmdbMovieList

	if err := json.NewDecoder(resp.Body).Decode(&movies); err != nil {
		log.Println("JSON: ", err)
		return err
	}

	if movies.Response == "False" {
		return errors.NotFound("wat.movie.api", "No movies found")
	}

	for _, m := range movies.Search {
		movie := &pb.Movie{}
		movie.Imdb = m.ImdbID
		movie.Name = m.Title
		movie.Year = m.Year
		movie.Poster = m.Poster
		res.Movies = append(res.Movies, movie)

	}

	return nil
}

func (s *service) Get(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {

	apiKey := os.Getenv("OMDB_API_KEY")
	log.Println(apiKey)
	name := strings.Replace(url.QueryEscape(req.Name), " ", "+", -1)

	url := fmt.Sprintf("http://www.omdbapi.com/?t=%s&type=movie&apikey=%s", name, apiKey)
	get, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("GetReq: ", err)
		return err
	}
	client := ProxyAwareHttpClient()

	resp, err := client.Do(get)
	if err != nil {
		log.Printf("Do: ", err)
		return err
	}

	defer resp.Body.Close()

	var movie OmdbMovie

	if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
		log.Println("JSON: ", err)
		return err
	}

	if movie.Response == "False" {
		return errors.NotFound("wat.movie.api", "Movie not found")
	}

	ret := &pb.Movie{}
	ret.Imdb = movie.ImdbID
	ret.Name = movie.Title
	ret.Year = movie.Year
	ret.Poster = movie.Poster
	ret.Genre = movie.Genre
	ret.Plot = movie.Plot
	ret.Rating = movie.ImdbRating

	res.Movie = ret
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

func ProxyAwareHttpClient() *http.Client {
	proxyServer, isSet := os.LookupEnv("GO_HTTP_PROXY")
	log.Println("ENV:")
	log.Println(os.Getenv("GO_HTTP_PROXY"))
	var proxyUrl *url.URL
	if isSet {
		log.Println(proxyServer)
		proxyUrl, err := url.Parse(proxyServer)
		if err != nil {
			log.Println("Invalid proxy url ", proxyUrl)
		}
	}

	httpTransport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}
	httpClient := &http.Client{Transport: httpTransport}
	return httpClient
}
