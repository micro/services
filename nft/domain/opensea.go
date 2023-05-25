package domain

import "encoding/json"

type AssetsResponse struct {
	Assets   []*Asset `json:"assets"`
	Next     string   `json:"next"`
	Previous string   `json:"previous"`
}

type CollectionsResponse struct {
	Collections []*Collection `json:"collections"`
}

type CollectionResponse struct {
	Collection *Collection `json:"collection"`
}

type Asset struct {
	Id          int32       `json:"id"`
	TokenId     string      `json:"token_id"`
	Sales       int32       `json:"num_sales"`
	ImageUrl    string      `json:"image_url"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Permalink   string      `json:"permalink"`
	Contract    *Contract   `json:"asset_contract"`
	Collection  *Collection `json:"collection"`
	Creator     *User       `json:"creator"`
	Owner       *User       `json:"owner"`
	LastSale    *Sale       `json:"last_sale,omitempty"`
	Presale     bool        `json:"is_presale"`
	ListingDate string      `json:"listing_date,omitempty"`

	Traits []map[string]interface{} `json:"traits,omitempty"`
}

type Contract struct {
	// name of contract
	Name string `json:"name,omitempty"`
	// ethereum address
	Address string `json:"address,omitempty"`
	// type of contract e.g "semi-fungible"
	Type string `json:"asset_contract_type,omitempty"`
	// timestamp of creation
	CreatedAt string `json:"created_date,omitempty"`
	// owner id
	Owner int32 `json:"owner,omitempty"`
	// aka "ERC1155"
	Schema string `json:"schema_name,omitempty"`
	// related symbol
	Symbol string `json:"symbol,omitempty"`
	// description of contract
	Description string `json:"description,omitempty"`
	// payout address
	PayoutAddress string `json:"payout_address,omitempty"`
	// seller fees
	SellerFees interface{} `json:"seller_fees_basis_points,omitempty"`
}

type Collection struct {
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
	Slug          string `json:"slug,omitempty"`
	ImageUrl      string `json:"image_url,omitempty"`
	CreatedAt     string `json:"created_date,omitempty"`
	PayoutAddress string `json:"payout_address,omitempty"`

	ExternalLink            string                 `json:"external_link,omitempty"`
	BannerImageUrl          string                 `json:"banner_image_url,omitempty"`
	DevSellerFeeBasisPoints interface{}            `json:"dev_seller_fee_basis_points,omitempty"`
	SafelistRequestStatus   string                 `json:"safelist_request_status,omitempty"`
	PrimaryAssetContracts   []Contract             `json:"primary_asset_contracts,omitempty"`
	Traits                  map[string]interface{} `json:"traits,omitempty"`
	PaymentTokens           []Token                `json:"payment_tokens,omitempty"`
	Editors                 []string               `json:"editors,omitempty"`
	Stats                   map[string]interface{} `json:"stats,omitempty"`
}

type User struct {
	User       *Username `json:"user"`
	ProfileUrl string    `json:"profile_img_url,omitempty"`
	Address    string    `json:"address,omitempty"`
}

type Username struct {
	Username string `json:"username",omitempty"`
}

type SaleAsset struct {
	TokenId  string `json:"token_id"`
	Decimals int32  `json:"decimals"`
}

type Sale struct {
	Asset          *SaleAsset   `json:"asset"`
	EventType      string       `json:"event_type,omitempty"`
	EventTimestamp string       `json:"event_timestamp,omitempty"`
	TotalPrice     string       `json:"total_price,omitempty"`
	Quantity       string       `json:"quantity,omitempty"`
	CreatedAt      string       `json:"created_date,omitempty"`
	Transaction    *Transaction `json:"transaction,omitempty"`
	PaymentToken   *Token       `json:"payment_token,omitempty"`
}

type Transaction struct {
	Id               int32       `json:"id,omitempty"`
	Timestamp        string      `json:"timestamp,omitempty"`
	BlockHash        string      `json:"block_hash,omitempty"`
	BlockNumber      interface{} `json:"block_number,omitempty"`
	FromAccount      *User       `json:"from_account,omitempty"`
	ToAccount        *User       `json:"to_account,omitempty"`
	TransactionHash  string      `json:"transaction_hash,omitempty"`
	TransactionIndex interface{} `json:"transaction_index,omitempty"`
}

type Token struct {
	Id       int32       `json:"id,omitempty"`
	Name     string      `json:"name,omitempty"`
	Symbol   string      `json:"symbol,omitempty"`
	Address  string      `json:"address,omitempty"`
	ImageUrl string      `json:"image_url,omitempty"`
	Decimals int32       `json:"decimals,omitempty"`
	EthPrice json.Number `json:"eth_price,omitempty"`
	UsdPrice json.Number `json:"usd_price,omitempty"`
}
