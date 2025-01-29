package handler

import (
	"context"
	"fmt"
	"net/url"

	"github.com/micro/micro/v5/service/errors"
	"github.com/micro/micro/v5/service/logger"
	"github.com/micro/services/pkg/api"
	"github.com/micro/services/quran/domain"
	pb "github.com/micro/services/quran/proto"
)

const (
	// the default api url
	apiUrl = "https://api.quran.com/api/v4/"
	// TODO: allow multiple translations
	// the default translation id
	translationId = "131,20"
	// TODO: allow multiple interpretations
	// the default tafsir id
	tafsirId = "169"
	// TODO: make configurable
	arabicText = "text_imlaei"
)

type Quran struct{}

func New() *Quran {
	// enable the api cache
	api.SetCache(true, 0)

	return new(Quran)
}

// Chapters returns a list of the chapters of the Quran
func (q *Quran) Chapters(ctx context.Context, req *pb.ChaptersRequest, rsp *pb.ChaptersResponse) error {
	lang := "en"
	if len(req.Language) > 0 {
		lang = req.Language
	}

	var resp map[string][]*domain.Chapter
	if err := api.Get(apiUrl+"chapters?language="+lang, &resp); err != nil {
		logger.Errorf("Failed to retrieve chapters: %v", err)
		return errors.InternalServerError("quran.chapters", "Failed to retrieve chapters")
	}

	for _, c := range resp["chapters"] {
		rsp.Chapters = append(rsp.Chapters, &pb.Chapter{
			Id:              c.Id,
			RevelationPlace: c.RevelationPlace,
			RevelationOrder: c.RevelationOrder,
			PrefixBismillah: c.BismillahPrefix,
			Name:            c.NameSimple,
			ComplexName:     c.NameComplex,
			ArabicName:      c.NameArabic,
			TranslatedName:  c.TranslatedName.Name,
			Verses:          c.VersesCount,
			Pages:           c.Pages,
		})
	}

	return nil
}

// Retrieve the summary for a given chapter
func (q *Quran) Summary(ctx context.Context, req *pb.SummaryRequest, rsp *pb.SummaryResponse) error {
	lang := "en"
	if len(req.Language) > 0 {
		lang = req.Language
	}
	if req.Chapter <= 0 {
		return errors.BadRequest("quran.chapter-summary", "require chapter id")
	}

	var resp map[string]*domain.ChapterInfo
	uri := fmt.Sprintf(apiUrl+"chapters/%d/info?language=%s", req.Chapter, lang)

	if err := api.Get(uri, &resp); err != nil {
		logger.Errorf("Failed to retrieve chapter info: %v", err)
		return errors.InternalServerError("quran.chapter-summary", "Failed to retrieve chapter summary")
	}

	info := resp["chapter_info"]
	rsp.Chapter = info.ChapterId
	rsp.Summary = info.ShortText
	rsp.Source = info.Source
	rsp.Text = info.Text

	return nil
}

// Return the verses by chapter
func (q *Quran) Verses(ctx context.Context, req *pb.VersesRequest, rsp *pb.VersesResponse) error {
	lang := "en"
	if len(req.Language) > 0 {
		lang = req.Language
	}

	if req.Chapter <= 0 {
		return errors.BadRequest("quran.verses", "require chapter id")
	}

	// TODO: enable configuring translations
	// comma separated list of resource ids
	// https://quran.api-docs.io/v4/resources/translations
	translations := translationId
	// TODO: enable configuring tafirs
	// https://api.quran.com/api/v4/resources/tafsirs
	tafsirs := tafsirId

	uri := fmt.Sprintf(apiUrl+"verses/by_chapter/%d?language=%s", req.Chapter, lang)

	// additional fields we require
	// arabic text in imlaei script
	uri += "&fields=" + arabicText

	uri += "&words=true"
	uri += "&word_fields=code_v2,text_imlaei"

	if len(translations) > 0 && req.Translate {
		uri += "&translations=" + translations
		uri += "&translation_fields=resource_name"
	}

	if len(tafsirs) > 0 && req.Interpret {
		uri += "&tafsirs=" + tafsirs
	}

	if req.Page > 0 {
		uri += fmt.Sprintf("&page=%d", req.Page)
	}

	if req.Limit > 0 {
		uri += fmt.Sprintf("&per_page=%d", req.Limit)
	}

	var resp *domain.VersesByChapter

	if err := api.Get(uri, &resp); err != nil {
		logger.Errorf("Failed to retrieve verses: %v", err)
		return errors.InternalServerError("quran.verses", "Failed to retrieve verses")
	}

	rsp.Chapter = req.Chapter
	rsp.Page = resp.Pagination.CurrentPage
	rsp.TotalPages = resp.Pagination.TotalPages

	for _, verse := range resp.Verses {
		v := domain.VerseToProto(verse)
		// strip words if not asked for
		if req.Words != true {
			v.Words = nil
		}
		rsp.Verses = append(rsp.Verses, v)
	}

	return nil
}

// Return the search results for a given query
func (q *Quran) Search(ctx context.Context, req *pb.SearchRequest, rsp *pb.SearchResponse) error {
	if len(req.Query) == 0 {
		return errors.BadRequest("quran.search", "missing search query")
	}

	lang := "en"
	if len(req.Language) > 0 {
		lang = req.Language
	}

	if req.Limit <= 0 {
		req.Limit = 20
	}

	qq := url.Values{}
	qq.Set("q", req.Query)
	qq.Set("size", fmt.Sprintf("%d", req.Limit))
	qq.Set("page", fmt.Sprintf("%d", req.Page))
	qq.Set("language", lang)

	uri := fmt.Sprintf(apiUrl+"search?%s", qq.Encode())

	var resp map[string]*domain.SearchResults

	if err := api.Get(uri, &resp); err != nil {
		logger.Errorf("Failed to retrieve search results: %v", err)
		return errors.InternalServerError("quran.search", "Failed to retrieve search results")
	}

	rsp.Query = req.Query
	rsp.Page = resp["search"].CurrentPage
	rsp.TotalPages = resp["search"].TotalPages
	rsp.TotalResults = resp["search"].TotalResults

	for _, result := range resp["search"].Results {
		r := domain.ResultToProto(result)
		rsp.Results = append(rsp.Results, r)
	}

	return nil
}
