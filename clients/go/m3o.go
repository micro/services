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
	"github.com/micro/services/clients/go/evchargers"
	"github.com/micro/services/clients/go/file"
	"github.com/micro/services/clients/go/forex"
	"github.com/micro/services/clients/go/function"
	"github.com/micro/services/clients/go/geocoding"
	"github.com/micro/services/clients/go/gifs"
	"github.com/micro/services/clients/go/google"
	"github.com/micro/services/clients/go/helloworld"
	"github.com/micro/services/clients/go/holidays"
	"github.com/micro/services/clients/go/id"
	"github.com/micro/services/clients/go/image"
	"github.com/micro/services/clients/go/ip"
	"github.com/micro/services/clients/go/location"
	"github.com/micro/services/clients/go/notes"
	"github.com/micro/services/clients/go/otp"
	"github.com/micro/services/clients/go/postcode"
	"github.com/micro/services/clients/go/prayer"
	"github.com/micro/services/clients/go/qr"
	"github.com/micro/services/clients/go/quran"
	"github.com/micro/services/clients/go/routing"
	"github.com/micro/services/clients/go/rss"
	"github.com/micro/services/clients/go/sentiment"
	"github.com/micro/services/clients/go/sms"
	"github.com/micro/services/clients/go/stock"
	"github.com/micro/services/clients/go/stream"
	"github.com/micro/services/clients/go/sunnah"
	"github.com/micro/services/clients/go/thumbnail"
	"github.com/micro/services/clients/go/time"
	"github.com/micro/services/clients/go/twitter"
	"github.com/micro/services/clients/go/url"
	"github.com/micro/services/clients/go/user"
	"github.com/micro/services/clients/go/vehicle"
	"github.com/micro/services/clients/go/weather"
	"github.com/micro/services/clients/go/youtube"
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
		EvchargersService: evchargers.NewEvchargersService(token),
		FileService:       file.NewFileService(token),
		ForexService:      forex.NewForexService(token),
		FunctionService:   function.NewFunctionService(token),
		GeocodingService:  geocoding.NewGeocodingService(token),
		GifsService:       gifs.NewGifsService(token),
		GoogleService:     google.NewGoogleService(token),
		HelloworldService: helloworld.NewHelloworldService(token),
		HolidaysService:   holidays.NewHolidaysService(token),
		IdService:         id.NewIdService(token),
		ImageService:      image.NewImageService(token),
		IpService:         ip.NewIpService(token),
		LocationService:   location.NewLocationService(token),
		NotesService:      notes.NewNotesService(token),
		OtpService:        otp.NewOtpService(token),
		PostcodeService:   postcode.NewPostcodeService(token),
		PrayerService:     prayer.NewPrayerService(token),
		QrService:         qr.NewQrService(token),
		QuranService:      quran.NewQuranService(token),
		RoutingService:    routing.NewRoutingService(token),
		RssService:        rss.NewRssService(token),
		SentimentService:  sentiment.NewSentimentService(token),
		SmsService:        sms.NewSmsService(token),
		StockService:      stock.NewStockService(token),
		StreamService:     stream.NewStreamService(token),
		SunnahService:     sunnah.NewSunnahService(token),
		ThumbnailService:  thumbnail.NewThumbnailService(token),
		TimeService:       time.NewTimeService(token),
		TwitterService:    twitter.NewTwitterService(token),
		UrlService:        url.NewUrlService(token),
		UserService:       user.NewUserService(token),
		VehicleService:    vehicle.NewVehicleService(token),
		WeatherService:    weather.NewWeatherService(token),
		YoutubeService:    youtube.NewYoutubeService(token),
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
	EvchargersService *evchargers.EvchargersService
	FileService       *file.FileService
	ForexService      *forex.ForexService
	FunctionService   *function.FunctionService
	GeocodingService  *geocoding.GeocodingService
	GifsService       *gifs.GifsService
	GoogleService     *google.GoogleService
	HelloworldService *helloworld.HelloworldService
	HolidaysService   *holidays.HolidaysService
	IdService         *id.IdService
	ImageService      *image.ImageService
	IpService         *ip.IpService
	LocationService   *location.LocationService
	NotesService      *notes.NotesService
	OtpService        *otp.OtpService
	PostcodeService   *postcode.PostcodeService
	PrayerService     *prayer.PrayerService
	QrService         *qr.QrService
	QuranService      *quran.QuranService
	RoutingService    *routing.RoutingService
	RssService        *rss.RssService
	SentimentService  *sentiment.SentimentService
	SmsService        *sms.SmsService
	StockService      *stock.StockService
	StreamService     *stream.StreamService
	SunnahService     *sunnah.SunnahService
	ThumbnailService  *thumbnail.ThumbnailService
	TimeService       *time.TimeService
	TwitterService    *twitter.TwitterService
	UrlService        *url.UrlService
	UserService       *user.UserService
	VehicleService    *vehicle.VehicleService
	WeatherService    *weather.WeatherService
	YoutubeService    *youtube.YoutubeService
}
