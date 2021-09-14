package domain

import (
	"strings"

	pb "github.com/micro/services/quran/proto"
)

var (
	// the default tafsir author
	tafsir = map[int32]string{
		169: "Tafsir Ibn Kathir",
	}
)

type Chapter struct {
	Id              int32           `json:"id"`
	RevelationPlace string          `json:"revelation_place"`
	RevelationOrder int32           `json:"revelation_order"`
	BismillahPrefix bool            `json:"bismillah_pre"`
	NameSimple      string          `json:"name_simple"`
	NameComplex     string          `json:"name_complex"`
	NameArabic      string          `json:"name_arabic"`
	VersesCount     int32           `json:"verses_count"`
	Pages           []int32         `json:"pages"`
	TranslatedName  *TranslatedName `json:"translated_name"`
}

type TranslatedName struct {
	LanguageName string `json:"language_name"`
	Name         string `json:"name"`
}

type TranslationText struct {
	Id           int32  `json:"id"`
	ResourceId   int32  `json:"resource_id"`
	Text         string `json:"text"`
	ResourceName string `json:"resource_name"`
}

type Translation struct {
	Text         string `json:"text"`
	LanguageName string `json:"language_name"`
}

type Transliteration struct {
	Text         string `json:"text"`
	LanguageName string `json:"language_name"`
}

type Tafsir struct {
	Id         int32  `json:"id"`
	ResourceId int32  `json:"resource_id"`
	Text       string `json:"text"`
}

type ChapterInfo struct {
	Id           int32  `json:"id"`
	ChapterId    int32  `json:"chapter_id"`
	LanguageName string `json:"language_name"`
	ShortText    string `json:"short_text"`
	Source       string `json:"source"`
	Text         string `json:"text"`
}

type Pagination struct {
	PerPage      int32 `json:"per_page"`
	CurrentPage  int32 `json:"current_page"`
	NextPage     int32 `json:"next_page"`
	TotalPages   int32 `json:"total_pages"`
	TotalRecords int32 `json:"total_records"`
}

type Verse struct {
	Id           int32              `json:"id"`
	VerseNumber  int32              `json:"verse_number"`
	VerseKey     string             `json:"verse_key"`
	JuzNumber    int32              `json:"juz_number"`
	HizbNumber   int32              `json:"hizb_number"`
	RubNumber    int32              `json:"rub_number"`
	PageNumber   int32              `json:"page_number"`
	Translations []*TranslationText `json:"translations"`
	Tafsirs      []*Tafsir          `json:"tafsirs"`
	Words        []*Word            `json:"words"`
	TextImlaei   string             `json:"text_imlaei"`
}

type Word struct {
	Id              int32        `json:"id"`
	Position        int32        `json:"position"`
	AudioUrl        string       `json:"audio_url"`
	CharTypeName    string       `json:"char_type_name"`
	CodeV1          string       `json:"code_v1"`
	PageNumber      int32        `json:"page_number"`
	LineNumber      int32        `json:"line_number"`
	Text            string       `json:"text_imlaei"`
	Code            string       `json:"code_v2"`
	Translation     *Translation `json:"translation"`
	Transliteration *Translation `json:"transliteration"`
}

type Result struct {
	VerseId      int32                `json:"verse_id"`
	VerseKey     string               `json:"verse_key"`
	Text         string               `json:"text"`
	Translations []*SearchTranslation `json:"translations"`
}

type SearchResults struct {
	Query        string    `json:"query"`
	TotalResults int32     `json:"total_results"`
	CurrentPage  int32     `json:"current_page"`
	TotalPages   int32     `json:"total_pages"`
	Results      []*Result `json:"results"`
}

type SearchTranslation struct {
	ResourceId int32  `json:"resource_id"`
	Text       string `json:"text"`
	Name       string `json:"name"`
}

type VersesByChapter struct {
	Pagination *Pagination `json:"pagination"`
	Verses     []*Verse    `json:"verses"`
}

func VerseToProto(verse *Verse) *pb.Verse {
	var transliteration []string
	var translation []string
	var words []*pb.Word

	for _, word := range verse.Words {
		words = append(words, &pb.Word{
			Id:              word.Id,
			Position:        word.Position,
			CharType:        word.CharTypeName,
			Page:            word.PageNumber,
			Line:            word.LineNumber,
			Text:            word.Text,
			Code:            word.Code,
			Translation:     word.Translation.Text,
			Transliteration: word.Transliteration.Text,
		})

		// skip the end
		if word.CharTypeName == "end" {
			continue
		}

		translation = append(translation, word.Translation.Text)
		transliteration = append(transliteration, word.Transliteration.Text)
	}

	var translations []*pb.Translation

	for _, tr := range verse.Translations {
		translations = append(translations, &pb.Translation{
			Id:     tr.Id,
			Source: tr.ResourceName,
			Text:   tr.Text,
		})
	}

	var interpretations []*pb.Interpretation

	for _, tf := range verse.Tafsirs {
		interpretations = append(interpretations, &pb.Interpretation{
			Id:     tf.Id,
			Source: tafsir[tf.ResourceId],
			Text:   tf.Text,
		})
	}

	return &pb.Verse{
		Id:              verse.Id,
		Key:             verse.VerseKey,
		Number:          verse.VerseNumber,
		Page:            verse.PageNumber,
		Text:            verse.TextImlaei,
		TranslatedText:  strings.Join(translation, " "),
		Transliteration: strings.Join(transliteration, " "),
		Words:           words,
		Translations:    translations,
		Interpretations: interpretations,
	}
}

func ResultToProto(r *Result) *pb.Result {
	var translations []*pb.Translation

	for _, tr := range r.Translations {
		translations = append(translations, &pb.Translation{
			Id:     tr.ResourceId,
			Source: tr.Name,
			Text:   tr.Text,
		})
	}

	return &pb.Result{
		VerseKey:     r.VerseKey,
		VerseId:      r.VerseId,
		Text:         r.Text,
		Translations: translations,
	}
}
