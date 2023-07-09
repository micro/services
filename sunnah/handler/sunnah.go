package handler

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/micro/services/pkg/api"
	"github.com/micro/services/sunnah/domain"
	pb "github.com/micro/services/sunnah/proto"
	"micro.dev/v4/service/errors"
	"micro.dev/v4/service/logger"
)

var (
	apiUrl = "https://api.sunnah.com/v1/"
)

type Sunnah struct {
	apiKey string
}

func New(key string) *Sunnah {
	api.SetKey("X-API-Key", key)
	api.SetCache(true, 0)

	return &Sunnah{
		apiKey: key,
	}
}

func (s *Sunnah) Collections(ctx context.Context, req *pb.CollectionsRequest, rsp *pb.CollectionsResponse) error {
	var resp *domain.CollectionRequest

	if req.Limit <= 0 {
		req.Limit = 50
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	uri := fmt.Sprintf("%scollections?limit=%d&page=%d", apiUrl, req.Limit, req.Page)

	if err := api.Get(uri, &resp); err != nil {
		logger.Errorf("Failed to retrieve collections: %v", err)
		return errors.InternalServerError("sunnah.collections", "Failed to retrieve collections")
	}

	for _, c := range resp.Data {
		var arabicTitle string
		if len(c.Collection) > 1 && c.Collection[1].Lang == "ar" {
			arabicTitle = c.Collection[1].Title
		}

		rsp.Collections = append(rsp.Collections, &pb.Collection{
			Name:        c.Name,
			Title:       c.Collection[0].Title,
			ArabicTitle: arabicTitle,
			Hadiths:     c.TotalHadith,
			Summary:     c.Collection[0].ShortIntro,
		})
	}

	return nil
}

func (s *Sunnah) Books(ctx context.Context, req *pb.BooksRequest, rsp *pb.BooksResponse) error {
	var resp *domain.BookRequest

	if len(req.Collection) == 0 {
		return errors.BadRequest("sunnah.books", "missing collection name")
	}

	if req.Limit <= 0 {
		req.Limit = 50
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	uri := fmt.Sprintf("%scollections/%s/books?limit=%d&page=%d", apiUrl, req.Collection, req.Limit, req.Page)

	if err := api.Get(uri, &resp); err != nil {
		logger.Errorf("Failed to retrieve books: %v", err)
		return errors.InternalServerError("sunnah.books", "Failed to retrieve books")
	}

	rsp.Collection = req.Collection
	rsp.Total = resp.Total
	rsp.Limit = req.Limit
	rsp.Page = req.Page

	for _, b := range resp.Data {
		if len(b.Book) == 0 {
			continue
		}

		var arabicName string
		if len(b.Book) > 1 && b.Book[1].Lang == "ar" {
			arabicName = b.Book[1].Name
		}
		bkId, _ := strconv.Atoi(b.BookNumber)
		rsp.Books = append(rsp.Books, &pb.Book{
			Id:         int32(bkId),
			Name:       b.Book[0].Name,
			ArabicName: arabicName,
			Hadiths:    b.NumberOfHadith,
		})
	}

	return nil
}

func (s *Sunnah) Chapters(ctx context.Context, req *pb.ChaptersRequest, rsp *pb.ChaptersResponse) error {
	var resp *domain.ChaptersRequest

	if len(req.Collection) == 0 {
		return errors.BadRequest("sunnah.chapters", "missing collection name")
	}

	if req.Book == 0 {
		req.Book = 1
	}

	if req.Limit <= 0 {
		req.Limit = 50
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	uri := fmt.Sprintf("%scollections/%s/books/%d/chapters?limit=%d&page=%d", apiUrl, req.Collection, req.Book, req.Limit, req.Page)

	if err := api.Get(uri, &resp); err != nil {
		logger.Errorf("Failed to retrieve chapters: %v", err)
		return errors.InternalServerError("sunnah.chapters", "Failed to retrieve chapters")
	}

	rsp.Collection = req.Collection
	rsp.Book = req.Book
	rsp.Total = resp.Total
	rsp.Limit = req.Limit
	rsp.Page = req.Page

	for _, c := range resp.Data {
		if len(c.Chapter) == 0 {
			continue
		}

		var arabicTitle string
		if len(c.Chapter) > 1 && c.Chapter[1].Lang == "ar" {
			arabicTitle = c.Chapter[1].ChapterTitle
		}
		bkId, _ := strconv.Atoi(c.BookNumber)
		chNumber, _ := strconv.Atoi(strings.Split(c.ChapterId, ".")[0])
		rsp.Chapters = append(rsp.Chapters, &pb.Chapter{
			Id:          int32(chNumber),
			Key:         c.ChapterId,
			Book:        int32(bkId),
			Title:       c.Chapter[0].ChapterTitle,
			ArabicTitle: arabicTitle,
		})
	}

	return nil
}

func (s *Sunnah) Hadiths(ctx context.Context, req *pb.HadithsRequest, rsp *pb.HadithsResponse) error {
	var resp *domain.HadithsRequest

	if len(req.Collection) == 0 {
		return errors.BadRequest("sunnah.hadiths", "missing collection name")
	}

	if req.Book == 0 {
		req.Book = 1
	}

	if req.Limit <= 0 {
		req.Limit = 50
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	uri := fmt.Sprintf("%scollections/%s/books/%d/hadiths?limit=%d&page=%d", apiUrl, req.Collection, req.Book, req.Limit, req.Page)

	if err := api.Get(uri, &resp); err != nil {
		logger.Errorf("Failed to retrieve hadiths: %v", err)
		return errors.InternalServerError("sunnah.hadiths", "Failed to retrieve hadiths")
	}

	rsp.Collection = req.Collection
	rsp.Book = req.Book
	rsp.Total = resp.Total
	rsp.Limit = req.Limit
	rsp.Page = req.Page

	for _, h := range resp.Data {
		if len(h.Hadith) == 0 {
			continue
		}

		var arabicTitle string
		var arabicText string

		if len(h.Hadith) > 1 && h.Hadith[1].Lang == "ar" {
			arabicTitle = h.Hadith[1].ChapterTitle
			arabicText = h.Hadith[1].Body
		}

		chNumber, _ := strconv.Atoi(strings.Split(h.ChapterId, ".")[0])
		hId, _ := strconv.Atoi(h.HadithNumber)

		rsp.Hadiths = append(rsp.Hadiths, &pb.Hadith{
			Id:                 int32(hId),
			Chapter:            int32(chNumber),
			ChapterKey:         h.ChapterId,
			ChapterTitle:       h.Hadith[0].ChapterTitle,
			Text:               h.Hadith[0].Body,
			ArabicText:         arabicText,
			ArabicChapterTitle: arabicTitle,
		})
	}

	return nil
}
