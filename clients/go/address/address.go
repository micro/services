package address

import (
	"github.com/m3o/m3o-go/client"
)

func NewAddressService(token string) *AddressService {
	return &AddressService{
		client: client.NewClient(&client.Options{
			Token: token,
		}),
	}
}

type AddressService struct {
	client *client.Client
}

// Lookup a list of UK addresses by postcode
func (t *AddressService) LookupPostcode(request *LookupPostcodeRequest) (*LookupPostcodeResponse, error) {
	rsp := &LookupPostcodeResponse{}
	return rsp, t.client.Call("address", "LookupPostcode", request, rsp)
}

type LookupPostcodeRequest struct {
	// UK postcode e.g SW1A 2AA
	Postcode string `json:"postcode"`
}

type LookupPostcodeResponse struct {
	Addresses []Record `json:"addresses"`
}

type Record struct {
	// building name
	BuildingName string `json:"buildingName"`
	// the county
	County string `json:"county"`
	// line one of address
	LineOne string `json:"lineOne"`
	// line two of address
	LineTwo string `json:"lineTwo"`
	// dependent locality
	Locality string `json:"locality"`
	// organisation if present
	Organisation string `json:"organisation"`
	// the postcode
	Postcode string `json:"postcode"`
	// the premise
	Premise string `json:"premise"`
	// street name
	Street string `json:"street"`
	// the complete address
	Summary string `json:"summary"`
	// post town
	Town string `json:"town"`
}
