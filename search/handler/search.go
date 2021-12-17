package handler

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/pkg/tenant"
	pb "github.com/micro/services/search/proto"
	open "github.com/opensearch-project/opensearch-go"
	openapi "github.com/opensearch-project/opensearch-go/opensearchapi"
	"google.golang.org/protobuf/types/known/structpb"
)

type Search struct {
	conf   conf
	client *open.Client
}

type conf struct {
	OpenAddr string `json:"open_addr"`
	User     string `json:"user"`
	Pass     string `json:"pass"`
	Insecure bool   `json:"insecure"`
}

type openSearchResponse struct {
	Took int64 `json:"took"`
	Hits hits  `json:"hits"`
}

type hits struct {
	Total map[string]interface{} `json:"total"`
	Hits  []hit                  `json:"hits"`
}
type hit struct {
	ID     string                 `json:"_id"`
	Score  float64                `json:"_score"`
	Source map[string]interface{} `json:"_source"`
}

func New(srv *service.Service) *Search {
	v, err := config.Get("micro.search")
	if err != nil {
		log.Fatalf("Failed to load config %s", err)
	}
	var c conf
	if err := v.Scan(&c); err != nil {
		log.Fatalf("Failed to load config %s", err)
	}
	if len(c.OpenAddr) == 0 || len(c.User) == 0 || len(c.Pass) == 0 {
		log.Fatalf("Missing configuration")
	}

	oc := open.Config{
		Addresses: []string{c.OpenAddr},
		Username:  c.User,
		Password:  c.Pass,
	}
	if c.Insecure {
		oc.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // For testing only. Use certificate for validation.
		}
	}

	client, err := open.NewClient(oc)
	if err != nil {
		log.Fatalf("Error configuring search client %s", err)
	}
	return &Search{
		conf:   c,
		client: client,
	}
}

func (s *Search) CreateIndex(ctx context.Context, request *pb.CreateIndexRequest, response *pb.CreateIndexResponse) error {
	method := "search.CreateIndex"

	// TODO validate name https://opensearch.org/docs/latest/opensearch/rest-api/index-apis/create-index/#index-naming-restrictions
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.Unauthorized(method, "Unauthorized")
	}
	req := openapi.CreateRequest{
		Index: indexName(tnt, request.IndexName),
		Body:  nil, // TODO populate with fields and their types
	}
	rsp, err := req.Do(ctx, s.client)
	if err != nil {
		log.Errorf("Error creating index %s", err)
		return errors.InternalServerError(method, "Error creating index")
	}
	defer rsp.Body.Close()
	if rsp.IsError() {
		log.Errorf("Error creating index %s", rsp.String())
		return errors.InternalServerError(method, "Error creating index")
	}
	return nil
}

func indexName(tnt, index string) string {
	return fmt.Sprintf("%s-%s", strings.ReplaceAll(tnt, "/", "-"), index)
}

func (s *Search) Index(ctx context.Context, request *pb.IndexRequest, response *pb.IndexResponse) error {
	method := "search.Index"
	// TODO validation
	// TODO validate name https://opensearch.org/docs/latest/opensearch/rest-api/index-apis/create-index/#index-naming-restrictions
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.Unauthorized(method, "Unauthorized")
	}
	if request.Document == nil {
		return errors.BadRequest(method, "Missing document param")
	}
	if len(request.Document.Id) == 0 {
		return errors.BadRequest(method, "Missing document.id param")
	}
	if len(request.Document.IndexName) == 0 {
		return errors.BadRequest(method, "Missing document.index_name param")
	}
	if request.Document.Contents == nil {
		return errors.BadRequest(method, "Missing document.contents param")
	}

	b, err := request.Document.Contents.MarshalJSON()
	if err != nil {
		return errors.BadRequest(method, "Error processing document")
	}
	req := openapi.IndexRequest{
		Index:      indexName(tnt, request.Document.IndexName),
		DocumentID: request.Document.Id,
		Body:       bytes.NewBuffer(b),
	}
	rsp, err := req.Do(ctx, s.client)
	if err != nil {
		log.Errorf("Error indexing doc %s", err)
		return errors.InternalServerError(method, "Error indexing document")
	}
	defer rsp.Body.Close()
	if rsp.IsError() {
		log.Errorf("Error indexing doc %s", rsp.String())
		return errors.InternalServerError(method, "Error indexing document")
	}
	return nil
}

func (s *Search) Delete(ctx context.Context, request *pb.DeleteRequest, response *pb.DeleteResponse) error {
	method := "search.Delete"
	// TODO validation
	// TODO validate name https://opensearch.org/docs/latest/opensearch/rest-api/index-apis/create-index/#index-naming-restrictions
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.Unauthorized(method, "Unauthorized")
	}
	req := openapi.DeleteRequest{
		Index:      indexName(tnt, request.IndexName),
		DocumentID: request.Id,
	}
	rsp, err := req.Do(ctx, s.client)
	if err != nil {
		log.Errorf("Error deleting doc %s", err)
		return errors.InternalServerError(method, "Error deleting document")
	}
	defer rsp.Body.Close()
	if rsp.IsError() {
		log.Errorf("Error deleting doc %s", rsp.String())
		return errors.InternalServerError(method, "Error deleting document")
	}
	return nil
}

func parseQueryDef(qd *pb.QueryDef) *simplejson.Json {
	js := recurseParseQueryDef(qd)
	ret := simplejson.New()
	ret.Set("query", js)
	return ret
}

func recurseParseQueryDef(qd *pb.QueryDef) *simplejson.Json {
	qs := simplejson.New()
	boolean := "must"
	if strings.ToLower(qd.Operator) == "or" {
		boolean = "should"
	}

	terms := []*simplejson.Json{}
	for _, v := range qd.Fields {
		matchType := "match"
		if qd.Prefix {
			matchType = "match_bool_prefix"
		}
		js := simplejson.New()
		js.SetPath([]string{matchType, v.FieldName}, v.Value)
		terms = append(terms, js)
	}
	// TODO reinstate once we fix protoc openapi3 recursive generation
	//for _, v := range qd.Queries {
	//	terms = append(terms, recurseParseQueryDef(v))
	//}

	qs.SetPath([]string{"bool", boolean}, terms)
	return qs
}

func (s *Search) Search(ctx context.Context, request *pb.SearchRequest, response *pb.SearchResponse) error {
	method := "search.Search"
	// TODO validation
	// TODO validate name https://opensearch.org/docs/latest/opensearch/rest-api/index-apis/create-index/#index-naming-restrictions
	if len(request.IndexName) == 0 {
		return errors.BadRequest(method, "Missing index_name param")
	}

	// Search models to support https://opensearch.org/docs/latest/opensearch/ux/
	// - Simple query
	// - Autocomplete (prefix) queries
	// - pagination
	// - Sorting
	//
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.Unauthorized(method, "Unauthorized")
	}

	// TODO fuzzy
	if request.Query == nil {
		return errors.BadRequest(method, "Missing query param")
	}

	qs := parseQueryDef(request.Query)
	// https://opensearch.org/docs/latest/opensearch/query-dsl/index/
	//q := fmt.Sprintf(`{"query": {"simple_query_string": {"query": "%s"} }}`, request.Query)
	b, _ := qs.MarshalJSON()
	log.Infof("Querying %v", string(b))
	req := openapi.SearchRequest{
		Index: []string{indexName(tnt, request.IndexName)},
		Body:  bytes.NewBuffer(b), // TODO - do we create our own DSL or just pass through the user string? support both?? simple and complex query

	}
	rsp, err := req.Do(ctx, s.client)
	if err != nil {
		log.Errorf("Error indexing doc %s", err)
		return errors.InternalServerError(method, "Error searching documents")
	}
	defer rsp.Body.Close()
	if rsp.IsError() {
		if rsp.StatusCode == 404 { // index not found
			return errors.NotFound(method, "Index not found")
		}
		log.Errorf("Error searching index %s", rsp.String())
		return errors.InternalServerError(method, "Error searching documents")
	}
	b, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Errorf("Error searching index %s", rsp.String())
		return errors.InternalServerError(method, "Error searching documents")
	}
	var os openSearchResponse
	if err := json.Unmarshal(b, &os); err != nil {
		log.Errorf("Error unmarshalling doc %s", err)
		return errors.InternalServerError(method, "Error searching documents")
	}
	log.Infof("%s", string(b))
	for _, v := range os.Hits.Hits {
		vs, err := structpb.NewStruct(v.Source)
		if err != nil {
			log.Errorf("Error unmarshalling doc %s", err)
			return errors.InternalServerError(method, "Error searching documents")
		}
		response.Documents = append(response.Documents, &pb.Document{
			Id:        v.ID,
			IndexName: request.IndexName,
			Contents:  vs,
		})
	}
	return nil
}

func (s *Search) DeleteIndex(ctx context.Context, request *pb.DeleteIndexRequest, response *pb.DeleteIndexResponse) error {
	method := "search.DeleteIndex"
	// TODO validation
	// TODO validate name https://opensearch.org/docs/latest/opensearch/rest-api/index-apis/create-index/#index-naming-restrictions
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.Unauthorized(method, "Unauthorized")
	}
	req := openapi.DeleteRequest{
		Index: indexName(tnt, request.IndexName),
	}
	rsp, err := req.Do(ctx, s.client)
	if err != nil {
		log.Errorf("Error deleting index %s", err)
		return errors.InternalServerError(method, "Error deleting index")
	}
	defer rsp.Body.Close()
	if rsp.IsError() {
		log.Errorf("Error deleting index %s", rsp.String())
		return errors.InternalServerError(method, "Error deleting index")
	}
	return nil
}
