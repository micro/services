import * as m3o from "@m3o/m3o-node";

export class VehicleService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Lookup a UK vehicle by it's registration number
  lookup(request: LookupRequest): Promise<LookupResponse> {
    return this.client.call(
      "vehicle",
      "Lookup",
      request
    ) as Promise<LookupResponse>;
  }
}

export interface LookupRequest {
  // the vehicle registration number
  registration?: string;
}

export interface LookupResponse {
  // co2 emmissions
  co2Emissions?: number;
  // colour of vehicle
  colour?: string;
  // engine capacity
  engineCapacity?: number;
  // fuel type e.g petrol, diesel
  fuelType?: string;
  // date of last v5 issue
  lastV5Issued?: string;
  // make of vehicle
  make?: string;
  // month of first registration
  monthOfFirstRegistration?: string;
  // mot expiry
  motExpiry?: string;
  // mot status
  motStatus?: string;
  // registration number
  registration?: string;
  // tax due data
  taxDueDate?: string;
  // tax status
  taxStatus?: string;
  // type approvale
  typeApproval?: string;
  // wheel plan
  wheelplan?: string;
  // year of manufacture
  yearOfManufacture?: number;
}
