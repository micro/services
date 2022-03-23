package handler

import (
	"context"
	"fmt"

	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/nft/domain"
	pb "github.com/micro/services/nft/proto"
	"github.com/micro/services/pkg/api"
	"google.golang.org/protobuf/types/known/structpb"
)

// OpenSea handler
type OpenSea struct {
	apiKey string
	// embed nft api
	*Nft
}

var (
	openseaURL = "https://api.opensea.io/api/v1"
)

func New() *OpenSea {
	v, err := config.Get("nft.key")
	if err != nil {
		logger.Fatal("nft.key config not found: %v", err)
	}

	key := v.String("")
	if len(key) == 0 {
		logger.Fatal("nft.key config not found")
	}

	// set the api key
	api.SetKey("X-API-KEY", key)
	// set the cache
	api.SetCache(true, time.Minute * 5)

	return &OpenSea{
		apiKey: key,
		Nft:    new(Nft),
	}
}

func (o *OpenSea) Assets(ctx context.Context, req *pb.AssetsRequest, rsp *pb.AssetsResponse) error {
	uri := openseaURL + "/assets"
	params := "?"

	limit := int32(20)
	order := "desc"
	orderBy := ""

	if req.Limit > 0 {
		limit = req.Limit
	}

	if req.Order == "asc" {
		order = "asc"
	}

	switch req.OrderBy {
	case "sale_date", "sale_count", "sale_price", "total_price":
		orderBy = req.OrderBy
	}

	params += fmt.Sprintf("limit=%d&order_direction=%s", limit, order)

	if len(req.Cursor) > 0 {
		params += fmt.Sprintf("&cursor=%s", req.Cursor)
	}

	if len(orderBy) > 0 {
		params += "&order_by=" + orderBy
	}

	if len(req.Collection) > 0 {
		params += "&collection=" + req.Collection
	}

	var resp domain.AssetsResponse

	if err := api.Get(uri+params, &resp); err != nil {
		return errors.InternalServerError("nft.assets", "failed to get assets: %v", err)
	}

	for _, asset := range resp.Assets {
		rsp.Assets = append(rsp.Assets, assetToPb(asset))
	}
	rsp.Next = resp.Next
	rsp.Previous = resp.Previous

	return nil
}

func paymentTokenToPb(token *domain.Token) *pb.Token {
	if token == nil {
		return &pb.Token{}
	}
	return &pb.Token{
		Id:       token.Id,
		Name:     token.Name,
		Symbol:   token.Symbol,
		Address:  token.Address,
		ImageUrl: token.ImageUrl,
		Decimals: token.Decimals,
		// converting to string for backwards compat
		EthPrice: fmt.Sprintf("%v", token.EthPrice),
		UsdPrice: fmt.Sprintf("%v", token.UsdPrice),
	}
}

func contractToPb(contract *domain.Contract) *pb.Contract {
	if contract == nil {
		return &pb.Contract{}
	}
	return &pb.Contract{
		Name:          contract.Name,
		Description:   contract.Description,
		Address:       contract.Address,
		Type:          contract.Type,
		CreatedAt:     contract.CreatedAt,
		Owner:         contract.Owner,
		Schema:        contract.Schema,
		Symbol:        contract.Symbol,
		PayoutAddress: contract.PayoutAddress,
		SellerFees:    contract.SellerFees,
	}
}

func collectionToPb(collection *domain.Collection) *pb.Collection {
	if collection == nil {
		return &pb.Collection{}
	}
	ret := &pb.Collection{
		Name:                  collection.Name,
		Description:           collection.Description,
		Slug:                  collection.Slug,
		ImageUrl:              collection.ImageUrl,
		CreatedAt:             collection.CreatedAt,
		PayoutAddress:         collection.PayoutAddress,
		ExternalLink:          collection.ExternalLink,
		BannerImageUrl:        collection.BannerImageUrl,
		SellerFees:            collection.DevSellerFeeBasisPoints,
		SafelistRequestStatus: collection.SafelistRequestStatus,
		PrimaryAssetContracts: func() []*pb.Contract {
			cons := make([]*pb.Contract, len(collection.PrimaryAssetContracts))
			for i, c := range collection.PrimaryAssetContracts {
				cons[i] = contractToPb(&c)
			}
			return cons
		}(),
		PaymentTokens: func() []*pb.Token {
			toks := make([]*pb.Token, len(collection.PaymentTokens))
			for i, t := range collection.PaymentTokens {
				toks[i] = paymentTokenToPb(&t)
			}
			return toks
		}(),
		Editors: collection.Editors,
	}
	ret.Traits, _ = structpb.NewStruct(collection.Traits)
	ret.Stats, _ = structpb.NewStruct(collection.Stats)
	return ret
}

func assetToPb(asset *domain.Asset) *pb.Asset {
	if asset.Creator == nil {
		asset.Creator = &domain.User{
			User: &domain.Username{},
		}
	}
	if asset.Creator.User == nil {
		asset.Creator.User = &domain.Username{}
	}
	if asset.Owner == nil {
		asset.Owner = &domain.User{
			User: &domain.Username{},
		}
	}
	if asset.Owner.User == nil {
		asset.Owner.User = &domain.Username{}
	}
	if asset.Collection == nil {
		asset.Collection = new(domain.Collection)
	}
	if asset.Contract == nil {
		asset.Contract = new(domain.Contract)
	}

	lastSale := new(pb.Sale)

	if asset.LastSale != nil {
		if asset.LastSale.Transaction == nil {
			asset.LastSale.Transaction = &domain.Transaction{
				FromAccount: &domain.User{User: new(domain.Username)},
				ToAccount:   &domain.User{User: new(domain.Username)},
			}
		}
		if asset.LastSale.Transaction.FromAccount == nil {
			asset.LastSale.Transaction.FromAccount = &domain.User{User: new(domain.Username)}
		}
		if asset.LastSale.Transaction.FromAccount.User == nil {
			asset.LastSale.Transaction.FromAccount.User = new(domain.Username)
		}
		if asset.LastSale.Transaction.ToAccount == nil {
			asset.LastSale.Transaction.ToAccount = &domain.User{User: new(domain.Username)}
		}
		if asset.LastSale.Transaction.ToAccount.User == nil {
			asset.LastSale.Transaction.ToAccount.User = new(domain.Username)
		}
		if asset.LastSale.PaymentToken == nil {
			asset.LastSale.PaymentToken = new(domain.Token)
		}

		lastSale = &pb.Sale{
			AssetTokenId:   asset.LastSale.Asset.TokenId,
			AssetDecimals:  asset.LastSale.Asset.Decimals,
			EventType:      asset.LastSale.EventType,
			EventTimestamp: asset.LastSale.EventTimestamp,
			TotalPrice:     asset.LastSale.TotalPrice,
			Quantity:       asset.LastSale.Quantity,
			CreatedAt:      asset.LastSale.CreatedAt,
			Transaction: &pb.Transaction{
				Id:          asset.LastSale.Transaction.Id,
				Timestamp:   asset.LastSale.Transaction.Timestamp,
				BlockHash:   asset.LastSale.Transaction.BlockHash,
				BlockNumber: asset.LastSale.Transaction.BlockNumber,
				FromAccount: &pb.User{
					Username:   asset.LastSale.Transaction.FromAccount.User.Username,
					ProfileUrl: asset.LastSale.Transaction.FromAccount.ProfileUrl,
					Address:    asset.LastSale.Transaction.FromAccount.Address,
				},
				ToAccount: &pb.User{
					Username:   asset.LastSale.Transaction.ToAccount.User.Username,
					ProfileUrl: asset.LastSale.Transaction.ToAccount.ProfileUrl,
					Address:    asset.LastSale.Transaction.ToAccount.Address,
				},
				TransactionHash:  asset.LastSale.Transaction.TransactionHash,
				TransactionIndex: asset.LastSale.Transaction.TransactionIndex,
			},
			PaymentToken: paymentTokenToPb(asset.LastSale.PaymentToken),
		}
	}
	traits := make([]*structpb.Struct, len(asset.Traits))
	for i, t := range asset.Traits {
		traits[i], _ = structpb.NewStruct(t)
	}

	return &pb.Asset{
		Name:        asset.Name,
		Description: asset.Description,
		Id:          asset.Id,
		TokenId:     asset.TokenId,
		ImageUrl:    asset.ImageUrl,
		Sales:       asset.Sales,
		Permalink:   asset.Permalink,
		Contract:    contractToPb(asset.Contract),
		Collection:  collectionToPb(asset.Collection),
		Owner: &pb.User{
			Username:   asset.Owner.User.Username,
			ProfileUrl: asset.Owner.ProfileUrl,
			Address:    asset.Owner.Address,
		},
		Creator: &pb.User{
			Username:   asset.Creator.User.Username,
			ProfileUrl: asset.Creator.ProfileUrl,
			Address:    asset.Creator.Address,
		},
		LastSale:    lastSale,
		Presale:     asset.Presale,
		ListingDate: asset.ListingDate,
		Traits:      traits,
	}
}

func (o *OpenSea) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	return errors.BadRequest("nft.create", "coming soon")
}

func (o *OpenSea) Collections(ctx context.Context, req *pb.CollectionsRequest, rsp *pb.CollectionsResponse) error {
	uri := openseaURL + "/collections"
	params := "?"

	limit := int32(20)
	offset := int32(0)

	if req.Limit > 0 {
		limit = req.Limit
	}

	if req.Offset > 0 {
		offset = req.Offset
	}

	params += fmt.Sprintf("limit=%d&offset=%d", limit, offset)

	var resp domain.CollectionsResponse

	if err := api.Get(uri+params, &resp); err != nil {
		return errors.InternalServerError("nft.collections", "failed to get collections: %v", err)
	}

	for _, collection := range resp.Collections {
		rsp.Collections = append(rsp.Collections, collectionToPb(collection))
	}

	return nil
}

func (o *OpenSea) Asset(ctx context.Context, req *pb.AssetRequest, rsp *pb.AssetResponse) error {
	if len(req.ContractAddress) == 0 {
		return errors.BadRequest("nft.asset", "Missing contract address param")
	}
	if len(req.TokenId) == 0 {
		return errors.BadRequest("nft.asset", "Missing token id param")
	}

	uri := fmt.Sprintf("%s/asset/%s/%s", openseaURL, req.ContractAddress, req.TokenId)

	var resp domain.Asset

	if err := api.Get(uri, &resp); err != nil {
		return errors.InternalServerError("nft.collection", "failed to get collection: %v", err)
	}

	rsp.Asset = assetToPb(&resp)

	return nil
}

func (o *OpenSea) Collection(ctx context.Context, req *pb.CollectionRequest, rsp *pb.CollectionResponse) error {
	if len(req.Slug) == 0 {
		return errors.BadRequest("nft.collection", "Missing slug param")
	}

	uri := fmt.Sprintf("%s/collection/%s", openseaURL, req.Slug)

	var resp domain.CollectionResponse

	if err := api.Get(uri, &resp); err != nil {
		return errors.InternalServerError("nft.collection", "failed to get collection: %v", err)
	}

	rsp.Collection = collectionToPb(resp.Collection)

	return nil
}
