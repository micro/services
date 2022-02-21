package handler

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	pauth "github.com/micro/services/pkg/auth"
	adminpb "github.com/micro/services/pkg/service/proto"
	"github.com/micro/services/pkg/tenant"
	pb "github.com/micro/services/search/proto"
	open "github.com/opensearch-project/opensearch-go"
	openapi "github.com/opensearch-project/opensearch-go/opensearchapi"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	indexNameRegex      = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9]$`)
	shortIndexNameRegex = regexp.MustCompile(`[a-zA-Z0-9]`)
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

type catIndicesEntry struct {
	Index string `json:"index"`
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

func isValidIndexName(s string) bool {
	if len(s) > 1 {
		return indexNameRegex.MatchString(s)
	}
	return shortIndexNameRegex.MatchString(s)
}

func (s *Search) CreateIndex(ctx context.Context, request *pb.CreateIndexRequest, response *pb.CreateIndexResponse) error {
	method := "search.CreateIndex"

	// TODO validate name https://opensearch.org/docs/latest/opensearch/rest-api/index-apis/create-index/#index-naming-restrictions
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.Unauthorized(method, "Unauthorized")
	}
	if !isValidIndexName(request.Index) {
		return errors.BadRequest(method, "Index name should contain only alphanumerics and hyphens")
	}
	req := openapi.IndicesCreateRequest{
		Index: indexName(tnt, request.Index),
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
	if request.Data == nil {
		return errors.BadRequest(method, "Missing data")
	}
	if len(request.Id) == 0 {
		request.Id = uuid.New().String()
	}
	if len(request.Index) == 0 {
		return errors.BadRequest(method, "Missing index")
	}
	if !isValidIndexName(request.Index) {
		return errors.BadRequest(method, "Index name should contain only alphanumerics and hyphens")
	}

	b, err := request.Data.MarshalJSON()
	if err != nil {
		return errors.BadRequest(method, "Error processing document")
	}
	req := openapi.CreateRequest{
		Index:      indexName(tnt, request.Index),
		DocumentID: request.Id,
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
	response.Record = &pb.Record{
		Id:   req.DocumentID,
		Data: request.Data,
	}

	return nil
}

func (s *Search) Delete(ctx context.Context, request *pb.DeleteRequest, response *pb.DeleteResponse) error {
	method := "search.Delete"
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.Unauthorized(method, "Unauthorized")
	}
	if len(request.Index) == 0 {
		return errors.BadRequest(method, "Missing index param")
	}
	req := openapi.DeleteRequest{
		Index:      indexName(tnt, request.Index),
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

func (s *Search) Search(ctx context.Context, request *pb.SearchRequest, response *pb.SearchResponse) error {
	method := "search.Search"
	if len(request.Index) == 0 {
		return errors.BadRequest(method, "Missing index param")
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
	if len(request.Query) == 0 {
		return errors.BadRequest(method, "Missing query param")
	}

	qs, err := parseQueryString(request.Query)
	if err != nil {
		log.Errorf("Error parsing string %s %s", request.Query, err)
		return errors.BadRequest(method, "%s", err)
	}
	b, _ := qs.MarshalJSON()
	req := openapi.SearchRequest{
		Index: []string{indexName(tnt, request.Index)},
		Body:  bytes.NewBuffer(b),
	}
	rsp, err := req.Do(ctx, s.client)
	if err != nil {
		log.Errorf("Error searching index %s", err)
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
		response.Records = append(response.Records, &pb.Record{
			Id:   v.ID,
			Data: vs,
		})
	}
	return nil
}

func (s *Search) DeleteIndex(ctx context.Context, request *pb.DeleteIndexRequest, response *pb.DeleteIndexResponse) error {
	method := "search.DeleteIndex"
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.Unauthorized(method, "Unauthorized")
	}
	if len(request.Index) == 0 {
		return errors.BadRequest(method, "Missing index param")
	}
	return s.deleteIndices(ctx, []string{indexName(tnt, request.Index)}, method)
}

func (s *Search) deleteIndices(ctx context.Context, indices []string, method string) error {
	req := openapi.IndicesDeleteRequest{
		Index: indices,
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
	log.Infof("Deleted indices: %v", indices)
	return nil

}

func (s *Search) DeleteData(ctx context.Context, request *adminpb.DeleteDataRequest, response *adminpb.DeleteDataResponse) error {
	method := "admin.DeleteData"
	_, err := pauth.VerifyMicroAdmin(ctx, method)
	if err != nil {
		return err
	}

	if len(request.TenantId) < 10 { // deliberate length check, don't want to unwittingly delete all the things
		return errors.BadRequest(method, "Missing tenant ID")
	}
	req := openapi.CatIndicesRequest{
		Format: "json",
	}
	rsp, err := req.Do(ctx, s.client)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	if rsp.IsError() {
		return err
	}
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	var entries []catIndicesEntry
	if err := json.Unmarshal(b, &entries); err != nil {
		return err
	}
	toDelete := []string{}
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Index, indexName(request.TenantId, "")) {
			continue
		}
		toDelete = append(toDelete, entry.Index)

	}
	if len(toDelete) > 0 {
		return s.deleteIndices(ctx, toDelete, method)
	}
	return nil
}
