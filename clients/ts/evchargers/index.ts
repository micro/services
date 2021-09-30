import * as m3o from "@m3o/m3o-node";

export class EvchargersService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Retrieve reference data as used by this API
  referenceData(request: ReferenceDataRequest): Promise<ReferenceDataResponse> {
    return this.client.call(
      "evchargers",
      "ReferenceData",
      request
    ) as Promise<ReferenceDataResponse>;
  }
  // Search by giving a coordinate and a max distance, or bounding box and optional filters
  search(request: SearchRequest): Promise<SearchResponse> {
    return this.client.call(
      "evchargers",
      "Search",
      request
    ) as Promise<SearchResponse>;
  }
}

export interface Address {
  // Any comments about how to access the charger
  accessComments?: string;
  addressLine1?: string;
  addressLine2?: string;
  country?: { [key: string]: any };
  countryId?: string;
  location?: Coordinates;
  postcode?: string;
  stateOrProvince?: string;
  title?: string;
  town?: string;
}

export interface BoundingBox {
  bottomLeft?: Coordinates;
  topRight?: Coordinates;
}

export interface ChargerType {
  comments?: string;
  id?: string;
  // Is this 40KW+
  isFastChargeCapable?: boolean;
  title?: string;
}

export interface CheckinStatusType {
  id?: string;
  isAutomated?: boolean;
  isPositive?: boolean;
  title?: string;
}

export interface Connection {
  // The amps offered
  amps?: number;
  connectionType?: ConnectionType;
  // The ID of the connection type
  connectionTypeId?: string;
  // The current
  current?: string;
  // The level of charging power available
  level?: string;
  // The power in KW
  power?: number;
  reference?: string;
  // The voltage offered
  voltage?: number;
}

export interface ConnectionType {
  formalName?: string;
  id?: string;
  isDiscontinued?: boolean;
  isObsolete?: boolean;
  title?: string;
}

export interface Coordinates {
  latitude?: number;
  longitude?: number;
}

export interface Country {
  continentCode?: string;
  id?: string;
  isoCode?: string;
  title?: string;
}

export interface CurrentType {
  description?: string;
  id?: string;
  title?: string;
}

export interface DataProvider {
  comments?: string;
  dataProviderStatusType?: DataProviderStatusType;
  id?: string;
  // How is this data licensed
  license?: string;
  title?: string;
  website?: string;
}

export interface DataProviderStatusType {
  id?: string;
  isProviderEnabled?: boolean;
  title?: string;
}

export interface Operator {
  comments?: string;
  contactEmail?: string;
  faultReportEmail?: string;
  id?: string;
  // Is this operator a private individual vs a company
  isPrivateIndividual?: boolean;
  phonePrimary?: string;
  phoneSecondary?: string;
  title?: string;
  website?: string;
}

export interface Poi {
  // The address
  address?: { [key: string]: any };
  // The connections available at this charge point
  connections?: Connection[];
  // The cost of charging
  cost?: string;
  // The ID of the data provider
  dataProviderId?: string;
  // The ID of the charger
  id?: string;
  // The number of charging points
  numPoints?: number;
  // The operator
  operator?: { [key: string]: any };
  // The ID of the operator of the charger
  operatorId?: string;
  // The type of usage
  usageType?: UsageType;
  // The type of usage for this charger point (is it public, membership required, etc)
  usageTypeId?: string;
}

export interface ReferenceDataRequest {}

export interface ReferenceDataResponse {
  // The types of charger
  chargerTypes?: ChargerType;
  // The types of checkin status
  checkinStatusTypes?: CheckinStatusType;
  // The types of connection
  connectionTypes?: ConnectionType;
  // The countries
  countries?: Country[];
  // The types of current
  currentTypes?: CurrentType;
  // The providers of the charger data
  dataProviders?: DataProvider;
  // The companies operating the chargers
  operators?: Operator[];
  // The status of the charger
  statusTypes?: StatusType;
  // The status of a submission
  submissionStatusTypes?: SubmissionStatusType;
  // The different types of usage
  usageTypes?: UsageType;
  // The types of user comment
  userCommentTypes?: UserCommentType;
}

export interface SearchRequest {
  // Bounding box to search within (top left and bottom right coordinates)
  box?: BoundingBox;
  // IDs of the connection type
  connectionTypes?: string;
  // Country ID
  countryId?: string;
  // Search distance from point in metres, defaults to 5000m
  distance?: number;
  // Supported charging levels
  levels?: string[];
  // Coordinates from which to begin search
  location?: Coordinates;
  // Maximum number of results to return, defaults to 100
  maxResults?: number;
  // Minimum power in KW. Note: data not available for many chargers
  minPower?: number;
  // IDs of the the EV charger operator
  operators?: string[];
  // Usage of the charge point (is it public, membership required, etc)
  usageTypes?: string;
}

export interface SearchResponse {
  pois?: Poi[];
}

export interface StatusType {
  id?: string;
  isOperational?: boolean;
  title?: string;
}

export interface SubmissionStatusType {
  id?: string;
  isLive?: boolean;
  title?: string;
}

export interface UsageType {
  id?: string;
  isAccessKeyRequired?: boolean;
  isMembershipRequired?: boolean;
  isPayAtLocation?: boolean;
  title?: string;
}

export interface UserCommentType {
  id?: string;
  title?: string;
}
