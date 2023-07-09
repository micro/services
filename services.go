package services

import (
	"github.com/micro/services/address/proto"
	"github.com/micro/services/ai/proto"
	"github.com/micro/services/app/proto"
	"github.com/micro/services/avatar/proto"
	"github.com/micro/services/bitcoin/proto"
	"github.com/micro/services/cache/proto"
	"github.com/micro/services/carbon/proto"
	"github.com/micro/services/chat/proto"
	"github.com/micro/services/comments/proto"
	"github.com/micro/services/contact/proto"
	"github.com/micro/services/cron/proto"
	"github.com/micro/services/crypto/proto"
	"github.com/micro/services/currency/proto"
	"github.com/micro/services/db/proto"
	"github.com/micro/services/dns/proto"
	"github.com/micro/services/email/proto"
	"github.com/micro/services/ethereum/proto"
	"github.com/micro/services/evchargers/proto"
	"github.com/micro/services/event/proto"
	"github.com/micro/services/file/proto"
	"github.com/micro/services/forex/proto"
	"github.com/micro/services/function/proto"
	"github.com/micro/services/geocoding/proto"
	"github.com/micro/services/gifs/proto"
	"github.com/micro/services/github/proto"
	"github.com/micro/services/google/proto"
	"github.com/micro/services/helloworld/proto"
	"github.com/micro/services/holidays/proto"
	"github.com/micro/services/id/proto"
	"github.com/micro/services/image/proto"
	"github.com/micro/services/ip/proto"
	"github.com/micro/services/lists/proto"
	"github.com/micro/services/location/proto"
	"github.com/micro/services/memegen/proto"
	"github.com/micro/services/minecraft/proto"
	"github.com/micro/services/movie/proto"
	"github.com/micro/services/mq/proto"
	"github.com/micro/services/news/proto"
	"github.com/micro/services/nft/proto"
	"github.com/micro/services/notes/proto"
	"github.com/micro/services/otp/proto"
	"github.com/micro/services/password/proto"
	"github.com/micro/services/ping/proto"
	"github.com/micro/services/place/proto"
	"github.com/micro/services/postcode/proto"
	"github.com/micro/services/prayer/proto"
	"github.com/micro/services/price/proto"
	"github.com/micro/services/qr/proto"
	"github.com/micro/services/quran/proto"
	"github.com/micro/services/routing/proto"
	"github.com/micro/services/rss/proto"
	"github.com/micro/services/search/proto"
	"github.com/micro/services/secret/proto"
	"github.com/micro/services/sentiment/proto"
	"github.com/micro/services/sms/proto"
	"github.com/micro/services/space/proto"
	"github.com/micro/services/spam/proto"
	"github.com/micro/services/stock/proto"
	"github.com/micro/services/stream/proto"
	"github.com/micro/services/sunnah/proto"
	"github.com/micro/services/thumbnail/proto"
	"github.com/micro/services/time/proto"
	"github.com/micro/services/translate/proto"
	"github.com/micro/services/tunnel/proto"
	"github.com/micro/services/twitter/proto"
	"github.com/micro/services/url/proto"
	"github.com/micro/services/user/proto"
	"github.com/micro/services/vehicle/proto"
	"github.com/micro/services/wallet/proto"
	"github.com/micro/services/weather/proto"
	"github.com/micro/services/youtube/proto"
	"micro.dev/v4/service/client"
)

type Client struct {
	Address    address.AddressService
	Ai         ai.AiService
	App        app.AppService
	Avatar     avatar.AvatarService
	Bitcoin    bitcoin.BitcoinService
	Cache      cache.CacheService
	Carbon     carbon.CarbonService
	Chat       chat.ChatService
	Comments   comments.CommentsService
	Contact    contact.ContactService
	Cron       cron.CronService
	Crypto     crypto.CryptoService
	Currency   currency.CurrencyService
	Db         db.DbService
	Dns        dns.DnsService
	Email      email.EmailService
	Ethereum   ethereum.EthereumService
	Evchargers evchargers.EvchargersService
	Event      event.EventService
	File       file.FileService
	Forex      forex.ForexService
	Function   function.FunctionService
	Geocoding  geocoding.GeocodingService
	Gifs       gifs.GifsService
	Github     github.GithubService
	Google     google.GoogleService
	Helloworld helloworld.HelloworldService
	Holidays   holidays.HolidaysService
	Id         id.IdService
	Image      image.ImageService
	Ip         ip.IpService
	Lists      lists.ListsService
	Location   location.LocationService
	Memegen    memegen.MemegenService
	Minecraft  minecraft.MinecraftService
	Movie      movie.MovieService
	Mq         mq.MqService
	News       news.NewsService
	Nft        nft.NftService
	Notes      notes.NotesService
	Otp        otp.OtpService
	Password   password.PasswordService
	Ping       ping.PingService
	Place      place.PlaceService
	Postcode   postcode.PostcodeService
	Prayer     prayer.PrayerService
	Price      price.PriceService
	Qr         qr.QrService
	Quran      quran.QuranService
	Routing    routing.RoutingService
	Rss        rss.RssService
	Search     search.SearchService
	Secret     secret.SecretService
	Sentiment  sentiment.SentimentService
	Sms        sms.SmsService
	Space      space.SpaceService
	Spam       spam.SpamService
	Stock      stock.StockService
	Stream     stream.StreamService
	Sunnah     sunnah.SunnahService
	Thumbnail  thumbnail.ThumbnailService
	Time       time.TimeService
	Translate  translate.TranslateService
	Tunnel     tunnel.TunnelService
	Twitter    twitter.TwitterService
	Url        url.UrlService
	User       user.UserService
	Vehicle    vehicle.VehicleService
	Wallet     wallet.WalletService
	Weather    weather.WeatherService
	Youtube    youtube.YoutubeService
}

func NewClient(c client.Client) *Client {
	return &Client{
		Address:    address.NewAddressService("address", c),
		Ai:         ai.NewAiService("ai", c),
		App:        app.NewAppService("app", c),
		Avatar:     avatar.NewAvatarService("avatar", c),
		Bitcoin:    bitcoin.NewBitcoinService("bitcoin", c),
		Cache:      cache.NewCacheService("cache", c),
		Carbon:     carbon.NewCarbonService("carbon", c),
		Chat:       chat.NewChatService("chat", c),
		Comments:   comments.NewCommentsService("comments", c),
		Contact:    contact.NewContactService("contact", c),
		Cron:       cron.NewCronService("cron", c),
		Crypto:     crypto.NewCryptoService("crypto", c),
		Currency:   currency.NewCurrencyService("currency", c),
		Db:         db.NewDbService("db", c),
		Dns:        dns.NewDnsService("dns", c),
		Email:      email.NewEmailService("email", c),
		Ethereum:   ethereum.NewEthereumService("ethereum", c),
		Evchargers: evchargers.NewEvchargersService("evchargers", c),
		Event:      event.NewEventService("event", c),
		File:       file.NewFileService("file", c),
		Forex:      forex.NewForexService("forex", c),
		Function:   function.NewFunctionService("function", c),
		Geocoding:  geocoding.NewGeocodingService("geocoding", c),
		Gifs:       gifs.NewGifsService("gifs", c),
		Github:     github.NewGithubService("github", c),
		Google:     google.NewGoogleService("google", c),
		Helloworld: helloworld.NewHelloworldService("helloworld", c),
		Holidays:   holidays.NewHolidaysService("holidays", c),
		Id:         id.NewIdService("id", c),
		Image:      image.NewImageService("image", c),
		Ip:         ip.NewIpService("ip", c),
		Lists:      lists.NewListsService("lists", c),
		Location:   location.NewLocationService("location", c),
		Memegen:    memegen.NewMemegenService("memegen", c),
		Minecraft:  minecraft.NewMinecraftService("minecraft", c),
		Movie:      movie.NewMovieService("movie", c),
		Mq:         mq.NewMqService("mq", c),
		News:       news.NewNewsService("news", c),
		Nft:        nft.NewNftService("nft", c),
		Notes:      notes.NewNotesService("notes", c),
		Otp:        otp.NewOtpService("otp", c),
		Password:   password.NewPasswordService("password", c),
		Ping:       ping.NewPingService("ping", c),
		Place:      place.NewPlaceService("place", c),
		Postcode:   postcode.NewPostcodeService("postcode", c),
		Prayer:     prayer.NewPrayerService("prayer", c),
		Price:      price.NewPriceService("price", c),
		Qr:         qr.NewQrService("qr", c),
		Quran:      quran.NewQuranService("quran", c),
		Routing:    routing.NewRoutingService("routing", c),
		Rss:        rss.NewRssService("rss", c),
		Search:     search.NewSearchService("search", c),
		Secret:     secret.NewSecretService("secret", c),
		Sentiment:  sentiment.NewSentimentService("sentiment", c),
		Sms:        sms.NewSmsService("sms", c),
		Space:      space.NewSpaceService("space", c),
		Spam:       spam.NewSpamService("spam", c),
		Stock:      stock.NewStockService("stock", c),
		Stream:     stream.NewStreamService("stream", c),
		Sunnah:     sunnah.NewSunnahService("sunnah", c),
		Thumbnail:  thumbnail.NewThumbnailService("thumbnail", c),
		Time:       time.NewTimeService("time", c),
		Translate:  translate.NewTranslateService("translate", c),
		Tunnel:     tunnel.NewTunnelService("tunnel", c),
		Twitter:    twitter.NewTwitterService("twitter", c),
		Url:        url.NewUrlService("url", c),
		User:       user.NewUserService("user", c),
		Vehicle:    vehicle.NewVehicleService("vehicle", c),
		Wallet:     wallet.NewWalletService("wallet", c),
		Weather:    weather.NewWeatherService("weather", c),
		Youtube:    youtube.NewYoutubeService("youtube", c),
	}
}
