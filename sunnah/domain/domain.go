package domain

type BookRequest struct {
	Data  []*Book `json:"data"`
	Total int32   `json:"total"`
	Limit int32   `json:"limit"`
}

type BookText struct {
	Lang string `json:"lang"`
	Name string `json:"name"`
}

type Book struct {
	BookNumber        string      `json:"bookNumber"`
	Book              []*BookText `json:"book"`
	HadithStartNumber int32       `json:"hadithStartNumber"`
	HadithEndNumber   int32       `json:"hadithEndNumber"`
	NumberOfHadith    int32       `json:"numberOfHadith"`
}

type ChaptersRequest struct {
	Data  []*Chapter `json:"data"`
	Total int32      `json:"total"`
	Limit int32      `json:"limit"`
}

type ChapterText struct {
	Lang          string `json:"lang"`
	ChapterNumber string `json:"chapterNumber"`
	ChapterTitle  string `json:"chapterTitle"`
}

type Chapter struct {
	BookNumber string         `json:"bookNumber"`
	ChapterId  string         `json:"chapterId"`
	Chapter    []*ChapterText `json:"chapter"`
}

type CollectionRequest struct {
	Data  []*Collection `json:"data"`
	Total int32         `json:"total"`
	Limit int32         `json:"limit"`
}

type CollectionText struct {
	Lang       string `json:"lang"`
	Title      string `json:"title"`
	ShortIntro string `json:"shortIntro"`
}

type Collection struct {
	Name                 string            `json:"name"`
	HasBooks             bool              `json:"hasBooks"`
	HasChapters          bool              `json:"hasChapters"`
	TotalHadith          int32             `json:"totalHadith"`
	TotalAvailableHadith int32             `json:"totalAvailableHadith"`
	Collection           []*CollectionText `json:"collection"`
}

type HadithsRequest struct {
	Data  []*Hadith `json:"data"`
	Total int32     `json:"total"`
	Limit int32     `json:"limit"`
}

type HadithText struct {
	Lang          string `json:"lang"`
	ChapterNumber string `json:"chapterNumber"`
	ChapterTitle  string `json:"chapterTitle"`
	Body          string `json:"body"`
}

type Hadith struct {
	Collection   string        `json:"collection"`
	BookNumber   string        `json:"bookNumber"`
	ChapterId    string        `json:"chapterId"`
	HadithNumber string        `json:"hadithNumber"`
	Hadith       []*HadithText `json:"hadith"`
}
