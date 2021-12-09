package handler

import (
	"context"
	"fmt"

	"github.com/micro/services/pkg/api"
	"github.com/micro/services/nft/domain"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/errors"
	pb "github.com/micro/services/nft/proto"
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

	return &OpenSea {
		apiKey: key,
		Nft: new(Nft),
	}
}

func (o *OpenSea) Assets(ctx context.Context, req *pb.AssetsRequest, rsp *pb.AssetsResponse) error {
	uri := openseaURL + "/assets"
	params := "?"

	limit := int32(20)
	offset := int32(0)
	order := "desc"
	orderBy := ""

	if req.Limit > 0 {
		limit = req.Limit
	}

	if req.Offset > 0 {
		offset = req.Offset
	}

	if req.Order == "asc" {
		order = "asc"
	}

	switch req.OrderBy {
	case "sale_date", "sale_count", "sale_price", "total_price":
		orderBy = req.OrderBy
	}

	params += fmt.Sprintf("limit=%d&offset=%d&order_direction=%s",
		limit, offset, order)

	if len(orderBy) > 0 {
		params += "&order_by=" + orderBy
	}

	if len(req.Collection) > 0 {
		params += "&collection=" + req.Collection
	}

	var resp domain.AssetsResponse

	if err := api.Get(uri + params, &resp); err != nil {
		return errors.InternalServerError("nft.assets", "failed to get assets: %v", err)
	}

	for _, asset := range resp.Assets {
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
					ToAccount: &domain.User{User: new(domain.Username)},
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
				AssetTokenId: asset.LastSale.Asset.TokenId,
				AssetDecimals: asset.LastSale.Asset.Decimals,
				EventType: asset.LastSale.EventType,
				EventTimestamp: asset.LastSale.EventTimestamp,
				TotalPrice: asset.LastSale.TotalPrice,
				Quantity: asset.LastSale.Quantity,
				CreatedAt: asset.LastSale.CreatedAt,
				Transaction: &pb.Transaction{
					Id: asset.LastSale.Transaction.Id,
					Timestamp: asset.LastSale.Transaction.Timestamp,
					BlockHash: asset.LastSale.Transaction.BlockHash,
					BlockNumber: asset.LastSale.Transaction.BlockNumber,
					FromAccount: &pb.User{
						Username: asset.LastSale.Transaction.FromAccount.User.Username,
						ProfileUrl: asset.LastSale.Transaction.FromAccount.ProfileUrl,
						Address: asset.LastSale.Transaction.FromAccount.Address,
					},
					ToAccount: &pb.User{
						Username: asset.LastSale.Transaction.ToAccount.User.Username,
						ProfileUrl: asset.LastSale.Transaction.ToAccount.ProfileUrl,
						Address: asset.LastSale.Transaction.ToAccount.Address,
					},
					TransactionHash: asset.LastSale.Transaction.TransactionHash,
					TransactionIndex: asset.LastSale.Transaction.TransactionIndex,
				},
				PaymentToken: &pb.Token{
					Id: asset.LastSale.PaymentToken.Id,
					Name: asset.LastSale.PaymentToken.Name,
					Symbol: asset.LastSale.PaymentToken.Symbol,
					Address: asset.LastSale.PaymentToken.Address,
					ImageUrl: asset.LastSale.PaymentToken.ImageUrl,
					Decimals: asset.LastSale.PaymentToken.Decimals,
					EthPrice: asset.LastSale.PaymentToken.EthPrice,
					UsdPrice: asset.LastSale.PaymentToken.UsdPrice,
				},
			}
		}

		rsp.Assets = append(rsp.Assets, &pb.Asset{
			Name: asset.Name,
			Description: asset.Description,
			Id: asset.Id,
			TokenId: asset.TokenId,
			ImageUrl: asset.ImageUrl,
			Sales: asset.Sales,
			Permalink: asset.Permalink,
			Contract: &pb.Contract{
				Name: asset.Contract.Name,
				Description: asset.Contract.Description,
				Address: asset.Contract.Address,
				Type: asset.Contract.Type,
				CreatedAt: asset.Contract.CreatedAt,
				Owner: asset.Contract.Owner,
				Schema: asset.Contract.Schema,
				Symbol: asset.Contract.Symbol,
				PayoutAddress: asset.Contract.PayoutAddress,
				SellerFees: asset.Contract.SellerFees,
			},
			Collection: &pb.Collection{
				Name: asset.Collection.Name,
				Description: asset.Collection.Description,
				Slug: asset.Collection.Slug,
				ImageUrl: asset.Collection.ImageUrl,
				CreatedAt: asset.Collection.CreatedAt,
				PayoutAddress: asset.Collection.PayoutAddress,
			},
			Owner: &pb.User{
				Username: asset.Owner.User.Username,
				ProfileUrl: asset.Owner.ProfileUrl,
				Address: asset.Owner.Address,
			},
			Creator: &pb.User{
				Username: asset.Creator.User.Username,
				ProfileUrl: asset.Creator.ProfileUrl,
				Address: asset.Creator.Address,
			},
			LastSale: lastSale,
			Presale: asset.Presale,
			ListingDate: asset.ListingDate,
		})
	}

	return nil
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

	params += fmt.Sprintf("limit=%d&offset=%d",limit, offset)

	var resp domain.CollectionsResponse

	if err := api.Get(uri + params, &resp); err != nil {
		return errors.InternalServerError("nft.collections", "failed to get collections: %v", err)
	}

	for _, collection := range resp.Collections {
		rsp.Collections = append(rsp.Collections, &pb.Collection{
			Name: collection.Name,
			Description: collection.Description,
			Slug: collection.Slug,
			ImageUrl: collection.ImageUrl,
			CreatedAt: collection.CreatedAt,
			PayoutAddress: collection.PayoutAddress,
		})
	}

	return nil
}

