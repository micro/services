package m3o

import (
	"github.com/micro/services/clients/go/address"
	"github.com/micro/services/clients/go/answer"
	"github.com/micro/services/clients/go/cache"
	"github.com/micro/services/clients/go/crypto"
	"github.com/micro/services/clients/go/currency"
	"github.com/micro/services/clients/go/db"
	"github.com/micro/services/clients/go/email"
	"github.com/micro/services/clients/go/emoji"
	"github.com/micro/services/clients/go/file"
	"github.com/micro/services/clients/go/forex"
	"github.com/micro/services/clients/go/geocoding"
	"github.com/micro/services/clients/go/helloworld"
	"github.com/micro/services/clients/go/id"
	"github.com/micro/services/clients/go/image"
	"github.com/micro/services/clients/go/ip"
	"github.com/micro/services/clients/go/location"
	"github.com/micro/services/clients/go/otp"
	"github.com/micro/services/clients/go/postcode"
	"github.com/micro/services/clients/go/quran"
	"github.com/micro/services/clients/go/routing"
	"github.com/micro/services/clients/go/rss"
	"github.com/micro/services/clients/go/sentiment"
	"github.com/micro/services/clients/go/sms"
	"github.com/micro/services/clients/go/stock"
	"github.com/micro/services/clients/go/stream"
	"github.com/micro/services/clients/go/thumbnail"
	"github.com/micro/services/clients/go/time"
	"github.com/micro/services/clients/go/url"
	"github.com/micro/services/clients/go/user"
	"github.com/micro/services/clients/go/weather"
)

func NewClient(token string) *Client {
	return &Client{
		token: token,

		AddressService:    address.NewAddressService(token),
		AnswerService:     answer.NewAnswerService(token),
		CacheService:      cache.NewCacheService(token),
		CryptoService:     crypto.NewCryptoService(token),
		CurrencyService:   currency.NewCurrencyService(token),
		DbService:         db.NewDbService(token),
		EmailService:      email.NewEmailService(token),
		EmojiService:      emoji.NewEmojiService(token),
		FileService:       file.NewFileService(token),
		ForexService:      forex.NewForexService(token),
		GeocodingService:  geocoding.NewGeocodingService(token),
		HelloworldService: helloworld.NewHelloworldService(token),
		IdService:         id.NewIdService(token),
		ImageService:      image.NewImageService(token),
		IpService:         ip.NewIpService(token),
		LocationService:   location.NewLocationService(token),
		OtpService:        otp.NewOtpService(token),
		PostcodeService:   postcode.NewPostcodeService(token),
		QuranService:      quran.NewQuranService(token),
		RoutingService:    routing.NewRoutingService(token),
		RssService:        rss.NewRssService(token),
		SentimentService:  sentiment.NewSentimentService(token),
		SmsService:        sms.NewSmsService(token),
		StockService:      stock.NewStockService(token),
		StreamService:     stream.NewStreamService(token),
		ThumbnailService:  thumbnail.NewThumbnailService(token),
		TimeService:       time.NewTimeService(token),
		UrlService:        url.NewUrlService(token),
		UserService:       user.NewUserService(token),
		WeatherService:    weather.NewWeatherService(token),
	}
}

type Client struct {
	token string

	AddressService    *address.AddressService
	AnswerService     *answer.AnswerService
	CacheService      *cache.CacheService
	CryptoService     *crypto.CryptoService
	CurrencyService   *currency.CurrencyService
	DbService         *db.DbService
	EmailService      *email.EmailService
	EmojiService      *emoji.EmojiService
	FileService       *file.FileService
	ForexService      *forex.ForexService
	GeocodingService  *geocoding.GeocodingService
	HelloworldService *helloworld.HelloworldService
	IdService         *id.IdService
	ImageService      *image.ImageService
	IpService         *ip.IpService
	LocationService   *location.LocationService
	OtpService        *otp.OtpService
	PostcodeService   *postcode.PostcodeService
	QuranService      *quran.QuranService
	RoutingService    *routing.RoutingService
	RssService        *rss.RssService
	SentimentService  *sentiment.SentimentService
	SmsService        *sms.SmsService
	StockService      *stock.StockService
	StreamService     *stream.StreamService
	ThumbnailService  *thumbnail.ThumbnailService
	TimeService       *time.TimeService
	UrlService        *url.UrlService
	UserService       *user.UserService
	WeatherService    *weather.WeatherService
}
