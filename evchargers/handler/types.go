package handler

type Poi struct {
	ID              int64        `bson:"ID"`
	DataProviderID  int64        `bson:"DataProviderID"`
	DataProvider    DataProvider `bson:"DataProvider"`
	OperatorID      int64        `bson:"OperatorID"`
	OperatorInfo    Operator     `bson:"OperatorInfo"`
	UsageTypeID     int64        `bson:"UsageTypeID"`
	UsageType       UsageType    `bson:"UsageType"`
	Cost            string       `bson:"UsageCost"`
	Address         Address      `bson:"AddressInfo"`
	Connections     []Connection `bson:"Connections"`
	NumberOfPoints  int64        `bson:"NumberOfPoints"`
	GeneralComments string       `bson:"GeneralComments"`
	StatusTypeID    int64        `bson:"StatusTypeID"`
	StatusType      StatusType   `bson:"StatusType"`
}

type Address struct {
	Title           string  `bson:"Title"`
	Latitude        float64 `bson:"Latitude"`
	Longitude       float64 `bson:"Longitude"`
	AddressLine1    string  `bson:"AddressLine1"`
	AddressLine2    string  `bson:"AddressLine2"`
	Town            string  `bson:"Town"`
	StateOrProvince string  `bson:"StateOrProvince"`
	AccessComments  string  `bson:"AccessComments"`
	Postcode        string  `bson:"Postcode"`
	CountryID       int64   `bson:"CountryID"`
	Country         Country `bson:"Country"`
}

type Country struct {
	ID            int64  `bson:"ID"`
	Title         string `bson:"Title"`
	ISOCode       string `bson:"ISOCode"`
	ContinentCode string `bson:"ContinentCode"`
}

type Connection struct {
	TypeID        int64          `bson:"ConnectionTypeID"`
	Type          ConnectionType `bson:"ConnectionType"`
	StatusTypeID  int64          `bson:"StatusTypeID"`
	StatusType    StatusType     `bson:"StatusType"`
	LevelID       int64          `bson:"LevelID"`
	Level         Level          `bson:"Level"`
	Amps          float64        `bson:"Amps"`
	Voltage       float64        `bson:"Voltage"`
	Power         float64        `bson:"PowerKW"`
	CurrentTypeID int64          `bson:"CurrentTypeID"`
	CurrentType   CurrentType    `bson:"CurrentType"`
	Quantity      int64          `bson:"Quantity"`
	Reference     string         `bson:"Reference"`
}

type Level struct {
	ID                  int64  `bson:"ID"`
	Title               string `bson:"Title"`
	Comments            string `bson:"Comments"`
	IsFastChargeCapable bool   `bson:"IsFastChargeCapable"`
}

type CurrentType struct {
	ID          int64  `bson:"ID"`
	Title       string `bson:"Title"`
	Description string `bson:"Description"`
}

type ConnectionType struct {
	ID             int64  `bson:"ID"`
	Title          string `bson:"Title"`
	FormalName     string `bson:"FormalName"`
	IsDiscontinued bool   `bson:"IsDiscontinued"`
	IsObsolete     bool   `bson:"IsObsolete"`
}

type DataProvider struct {
	ID                 int64              `bson:"ID"`
	Title              string             `bson:"Title"`
	WebsiteURL         string             `bson:"WebsiteURL"`
	Comments           string             `bson:"Comments"`
	DataProviderStatus DataProviderStatus `bson:"DataProviderStatusType"`
	IsOpenDataLicensed bool               `bson:"IsOpenDataLicensed"`
	License            string             `bson:"License"`
}

type DataProviderStatus struct {
	ID                int64  `bson:"ID"`
	Title             string `bson:"Title"`
	IsProviderEnabled bool   `bson:"IsProviderEnabled"`
}

type Operator struct {
	ID                  int64  `bson:"ID"`
	Title               string `bson:"Title"`
	WebsiteURL          string `bson:"WebsiteURL"`
	Comments            string `bson:"Comments"`
	PhonePrimary        string `bson:"PhonePrimaryContact"`
	PhoneSecondary      string `bson:"PhoneSecondaryContact"`
	IsPrivateIndividual bool   `bson:"IsPrivateIndividual"`
	ContactEmail        string `bson:"ContactEmail"`
	FaultReportEmail    string `bson:"FaultReportEmail"`
}

type UsageType struct {
	ID                   int64  `bson:"ID"`
	Title                string `bson:"Title"`
	IsPayAtLocation      bool   `bson:"IsPayAtLocation"`
	IsMembershipRequired bool   `bson:"IsMembershipRequired"`
	IsAccessKeyRequired  bool   `bson:"IsAccessKeyRequired"`
}

type StatusType struct {
	ID                int64  `bson:"ID"`
	Title             string `bson:"Title"`
	IsUsageSelectable bool   `bson:"IsUsageSelectable"`
	IsOperational     bool   `bson:"IsOperational"`
}
