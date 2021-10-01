package handler

type Poi struct {
	ID              int32        `bson:"ID" json:"ID"`
	DataProviderID  int32        `bson:"DataProviderID" json:"DataProviderID"`
	DataProvider    DataProvider `bson:"DataProvider" json:"DataProvider"`
	OperatorID      int32        `bson:"OperatorID" json:"OperatorID"`
	OperatorInfo    Operator     `bson:"OperatorInfo" json:"OperatorInfo"`
	UsageTypeID     int32        `bson:"UsageTypeID" json:"UsageTypeID"`
	UsageType       UsageType    `bson:"UsageType" json:"UsageType"`
	Cost            string       `bson:"UsageCost" json:"UsageCost"`
	Address         Address      `bson:"AddressInfo" json:"AddressInfo"`
	Connections     []Connection `bson:"Connections" json:"Connections"`
	NumberOfPoints  int32        `bson:"NumberOfPoints" json:"NumberOfPoints"`
	GeneralComments string       `bson:"GeneralComments" json:"GeneralComments"`
	StatusTypeID    int32        `bson:"StatusTypeID" json:"StatusTypeID"`
	StatusType      StatusType   `bson:"StatusType" json:"StatusType"`
	SpatialPosition Position     `bson:"SpatialPosition" json:"SpatialPosition"`
}

type Position struct {
	Type        string    `bson:"type" json:"type"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"`
}

type Address struct {
	Title           string  `bson:"Title" json:"Title"`
	Latitude        float64 `bson:"Latitude" json:"Latitude"`
	Longitude       float64 `bson:"Longitude" json:"Longitude"`
	AddressLine1    string  `bson:"AddressLine1" json:"AddressLine1"`
	AddressLine2    string  `bson:"AddressLine2" json:"AddressLine2"`
	Town            string  `bson:"Town" json:"Town"`
	StateOrProvince string  `bson:"StateOrProvince" json:"StateOrProvince"`
	AccessComments  string  `bson:"AccessComments" json:"AccessComments"`
	Postcode        string  `bson:"Postcode" json:"Postcode"`
	CountryID       int32   `bson:"CountryID" json:"CountryID"`
	Country         Country `bson:"Country" json:"Country"`
}

type Country struct {
	ID            int32  `bson:"ID" json:"ID"`
	Title         string `bson:"Title" json:"Title"`
	ISOCode       string `bson:"ISOCode" json:"ISOCode"`
	ContinentCode string `bson:"ContinentCode" json:"ContinentCode"`
}

type Connection struct {
	TypeID        int32          `bson:"ConnectionTypeID" json:"ConnectionTypeID"`
	Type          ConnectionType `bson:"ConnectionType" json:"ConnectionType"`
	StatusTypeID  int32          `bson:"StatusTypeID" json:"StatusTypeID"`
	StatusType    StatusType     `bson:"StatusType" json:"StatusType"`
	LevelID       int32          `bson:"LevelID" json:"LevelID"`
	Level         ChargerType    `bson:"Level" json:"Level"`
	Amps          float64        `bson:"Amps" json:"Amps"`
	Voltage       float64        `bson:"Voltage" json:"Voltage"`
	Power         float64        `bson:"PowerKW" json:"PowerKW"`
	CurrentTypeID int32          `bson:"CurrentTypeID" json:"CurrentTypeID"`
	CurrentType   CurrentType    `bson:"CurrentType" json:"CurrentType"`
	Quantity      int32          `bson:"Quantity" json:"Quantity"`
	Reference     string         `bson:"Reference" json:"Reference"`
}

type ChargerType struct {
	ID                  int32  `bson:"ID" json:"ID"`
	Title               string `bson:"Title" json:"Title"`
	Comments            string `bson:"Comments" json:"Comments"`
	IsFastChargeCapable bool   `bson:"IsFastChargeCapable" json:"IsFastChargeCapable"`
}

type CurrentType struct {
	ID          int32  `bson:"ID" json:"ID"`
	Title       string `bson:"Title" json:"Title"`
	Description string `bson:"Description" json:"Description"`
}

type ConnectionType struct {
	ID             int32  `bson:"ID" json:"ID"`
	Title          string `bson:"Title" json:"Title"`
	FormalName     string `bson:"FormalName" json:"FormalName"`
	IsDiscontinued bool   `bson:"IsDiscontinued" json:"IsDiscontinued"`
	IsObsolete     bool   `bson:"IsObsolete" json:"IsObsolete"`
}

type DataProvider struct {
	ID                 int32              `bson:"ID" json:"ID"`
	Title              string             `bson:"Title" json:"Title"`
	WebsiteURL         string             `bson:"WebsiteURL" json:"WebsiteURL"`
	Comments           string             `bson:"Comments" json:"Comments"`
	DataProviderStatus DataProviderStatus `bson:"DataProviderStatusType" json:"DataProviderStatusType"`
	IsOpenDataLicensed bool               `bson:"IsOpenDataLicensed" json:"IsOpenDataLicensed"`
	License            string             `bson:"License" json:"License"`
}

type DataProviderStatus struct {
	ID                int32  `bson:"ID" json:"ID"`
	Title             string `bson:"Title" json:"Title"`
	IsProviderEnabled bool   `bson:"IsProviderEnabled" json:"IsProviderEnabled"`
}

type Operator struct {
	ID                  int32  `bson:"ID" json:"ID"`
	Title               string `bson:"Title" json:"Title"`
	WebsiteURL          string `bson:"WebsiteURL" json:"WebsiteURL"`
	Comments            string `bson:"Comments" json:"Comments"`
	PhonePrimary        string `bson:"PhonePrimaryContact" json:"PhonePrimaryContact"`
	PhoneSecondary      string `bson:"PhoneSecondaryContact" json:"PhoneSecondaryContact"`
	IsPrivateIndividual bool   `bson:"IsPrivateIndividual" json:"IsPrivateIndividual"`
	ContactEmail        string `bson:"ContactEmail" json:"ContactEmail"`
	FaultReportEmail    string `bson:"FaultReportEmail" json:"FaultReportEmail"`
}

type UsageType struct {
	ID                   int32  `bson:"ID" json:"ID"`
	Title                string `bson:"Title" json:"Title"`
	IsPayAtLocation      bool   `bson:"IsPayAtLocation" json:"IsPayAtLocation"`
	IsMembershipRequired bool   `bson:"IsMembershipRequired" json:"IsMembershipRequired"`
	IsAccessKeyRequired  bool   `bson:"IsAccessKeyRequired" json:"IsAccessKeyRequired"`
}

type StatusType struct {
	ID                int32  `bson:"ID" json:"ID"`
	Title             string `bson:"Title" json:"Title"`
	IsUsageSelectable bool   `bson:"IsUsageSelectable" json:"IsUsageSelectable"`
	IsOperational     bool   `bson:"IsOperational" json:"IsOperational"`
}

type UserCommentType struct {
	ID    int32  `bson:"ID" json:"ID"`
	Title string `bson:"Title" json:"Title"`
}

type CheckinStatusType struct {
	ID                 int32  `bson:"ID" json:"ID"`
	Title              string `bson:"Title" json:"Title"`
	IsPositive         bool   `bson:"IsPositive" json:"IsPositive"`
	IsAutomatedCheckin bool   `bson:"IsAutomatedCheckin" json:"IsAutomatedCheckin"`
}

type ReferenceData struct {
	ChargerTypes          []ChargerType          `bson:"ChargerTypes" json:"ChargerTypes"`
	ConnectionTypes       []ConnectionType       `bson:"ConnectionTypes" json:"ConnectionTypes"`
	CurrentTypes          []CurrentType          `bson:"CurrentTypes" json:"CurrentTypes"`
	Countries             []Country              `bson:"Countries" json:"Countries"`
	DataProviders         []DataProvider         `bson:"DataProviders" json:"DataProviders"`
	Operators             []Operator             `bson:"Operators" json:"Operators"`
	StatusTypes           []StatusType           `bson:"StatusTypes" json:"StatusTypes"`
	UsageTypes            []UsageType            `bson:"UsageTypes" json:"UsageTypes"`
	UserCommentTypes      []UserCommentType      `bson:"UserCommentTypes" json:"UserCommentTypes"`
	CheckinStatusTypes    []CheckinStatusType    `bson:"CheckinStatusTypes" json:"CheckinStatusTypes"`
	SubmissionStatusTypes []SubmissionStatusType `bson:"SubmissionStatusTypes" json:"SubmissionStatusTypes"`
}

type SubmissionStatusType struct {
	ID     int32  `bson:"ID" json:"ID"`
	Title  string `bson:"Title" json:"Title"`
	IsLive bool   `bson:"IsLive" json:"IsLive"`
}
