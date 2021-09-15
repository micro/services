import * as m3o from "@m3o/m3o-node";

export class OtpService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Generate an OTP (one time pass) code
  generate(request: GenerateRequest): Promise<GenerateResponse> {
    return this.client.call(
      "otp",
      "Generate",
      request
    ) as Promise<GenerateResponse>;
  }
  // Validate the OTP code
  validate(request: ValidateRequest): Promise<ValidateResponse> {
    return this.client.call(
      "otp",
      "Validate",
      request
    ) as Promise<ValidateResponse>;
  }
}

export interface GenerateRequest {
  // expiration in seconds (default: 60)
  expiry?: number;
  // unique id, email or user to generate an OTP for
  id?: string;
  // number of characters (default: 6)
  size?: number;
}

export interface GenerateResponse {
  // one time pass code
  code?: string;
}

export interface ValidateRequest {
  // one time pass code to validate
  code?: string;
  // unique id, email or user for which the code was generated
  id?: string;
}

export interface ValidateResponse {
  // returns true if the code is valid for the ID
  success?: boolean;
}
