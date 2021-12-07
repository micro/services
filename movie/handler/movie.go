package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"

	pb "github.com/micro/services/movie/proto"
)

type Movie struct {
	Api string
	Key string
}

func New() *Movie {
	v, err := config.Get("movie.api")
	if err != nil {
		logger.Fatal("movie.api config not found: %v", err)
	}

	api := v.String("")
	if len(api) == 0 {
		logger.Fatal("movie.api config not found")
	}

	v, err = config.Get("movie.key")
	if err != nil {
		logger.Fatal("movie.key config not found: %v", err)
	}

	key := v.String("")
	if len(key) == 0 {
		logger.Fatal("movie.key config not found")
	}

	return &Movie{
		Api: api,
		Key: key,
	}
}

func (m *Movie) Search(_ context.Context, req *pb.SearchRequest, rsp *pb.SearchResponse) error {
	if req.Page == 0 {
		req.Page = 1
	}

	vals := url.Values{}
	vals.Set("api_key", m.Key)
	vals.Set("query", req.Query)
	vals.Set("language", req.Language)
	vals.Set("page", fmt.Sprintf("%d", req.Page))
	vals.Set("include_adult", "false")
	vals.Set("region", req.Region)
	if req.Year > 0 {
		vals.Set("year", fmt.Sprintf("%d", req.Year))
	}
	if req.PrimaryReleaseYear > 0 {
		vals.Set("primary_release_year", fmt.Sprintf("%d", req.PrimaryReleaseYear))
	}

	api := fmt.Sprintf("%s/search/movie?%s", m.Api, vals.Encode())

	resp, err := http.Get(api)
	if err != nil {
		logger.Errorf("Failed to get movie search results: %v\n", err)
		return errors.InternalServerError("movie.search", "failed to get movie search results")
	}

	if resp.StatusCode != http.StatusOK {
		logger.Errorf("Movie search api status code is not OK! status=%d\n", resp.StatusCode)
		return errors.InternalServerError("movie.search", fmt.Sprintf("movie search status is not 200, it's %d", resp.StatusCode))
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Errorf("Movie search close response body error: %v\n", err)
		}
	}()

	b, _ := ioutil.ReadAll(resp.Body)

	var respBody map[string]interface{}

	if err := json.Unmarshal(b, &respBody); err != nil {
		logger.Errorf("Failed to unmarshal movie search results: %v\n", err)
		return errors.InternalServerError("movie.search", "failed to unmarshal movie search results")
	}

	rsp.Page = int32(respBody["page"].(float64))
	rsp.TotalPages = int32(respBody["total_pages"].(float64))
	rsp.TotalResults = int32(respBody["total_results"].(float64))

	results := respBody["results"].([]interface{})

	for _, v := range results {
		info := v.(map[string]interface{})

		genreIds := make([]int32, 0)
		if ids, ok := info["genre_ids"].([]interface{}); ok {
			for _, id := range ids {
				genreIds = append(genreIds, int32(id.(float64)))
			}
		}

		postPath := ""
		if p, ok := info["poster_path"].(string); ok {
			postPath = p
		}

		backdropPath := ""
		if b, ok := info["backdrop_path"].(string); ok {
			backdropPath = b
		}

		rsp.Results = append(rsp.Results, &pb.MovieInfo{
			PosterPath:       postPath,
			Adult:            info["adult"].(bool),
			Overview:         info["overview"].(string),
			ReleaseDate:      info["release_date"].(string),
			GenreIds:         genreIds,
			Id:               int32(info["id"].(float64)),
			OriginalTitle:    info["original_title"].(string),
			OriginalLanguage: info["original_language"].(string),
			Title:            info["title"].(string),
			BackdropPath:     backdropPath,
			Popularity:       info["popularity"].(float64),
			VoteCount:        int32(info["vote_count"].(float64)),
			Video:            info["video"].(bool),
			VoteAverage:      info["vote_average"].(float64),
		})

	}

	return nil
}
