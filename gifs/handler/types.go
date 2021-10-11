package handler

type searchResponse struct {
	Data       []gif      `json:"data"`
	Pagination pagination `json:"pagination"`
}

type gif struct {
	ID       string  `json:"id"`
	Slug     string  `json:"slug"`
	URL      string  `json:"url"`
	ShortURL string  `json:"bitly_url"`
	EmbedURL string  `json:"embed_url"`
	Source   string  `json:"source"`
	Rating   string  `json:"rating"`
	Images   formats `json:"images"`
	Title    string  `json:"title"`
}

type formats struct {
	Original               format `json:"original"`
	Downsized              format `json:"downsized"`
	FixedHeight            format `json:"fixed_height"`
	FixedHeightStill       format `json:"fixed_height_still"`
	FixedHeightDownsampled format `json:"fixed_height_downsampled"`
	FixedWidth             format `json:"fixed_width"`
	FixedWidthStill        format `json:"fixed_width_still"`
	FixedWidthDownsampled  format `json:"fixed_width_downsampled"`
	FixedHeightSmall       format `json:"fixed_height_small"`
	FixedHeightSmallStill  format `json:"fixed_height_small_still"`
	FixedWidthSmall        format `json:"fixed_width_small"`
	FixedWidthSmallStill   format `json:"fixed_width_small_still"`
	DownsizedStill         format `json:"downsized_still"`
	DownsizedLarge         format `json:"downsized_large"`
	DownsizedMedium        format `json:"downsized_medium"`
	DownsizedSmall         format `json:"downsized_small"`
	OriginalStill          format `json:"original_still"`
	Looping                format `json:"looping"`
	Preview                format `json:"preview"`
	PreviewGif             format `json:"preview_gif"`
}

type format struct {
	Height   string `json:"height"`
	Width    string `json:"width"`
	Size     string `json:"size"`
	URL      string `json:"url"`
	MP4URL   string `json:"mp4_url"`
	MP4Size  string `json:"mp4_size"`
	WebpURL  string `json:"webp_url"`
	WebpSize string `json:"webp_size"`
}

type pagination struct {
	Offset     int32 `json:"offset"`
	TotalCount int32 `json:"total_count"`
	Count      int32 `json:"count"`
}
