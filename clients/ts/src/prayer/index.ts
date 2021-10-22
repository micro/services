import * as m3o from "@m3o/m3o-node";

export class PrayerService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Get the prayer (salah) times for a location on a given date
  times(request: TimesRequest): Promise<TimesResponse> {
    return this.client.call(
      "prayer",
      "Times",
      request
    ) as Promise<TimesResponse>;
  }
}

export interface PrayerTime {
  // asr time
  asr?: string;
  // date for prayer times in YYYY-MM-DD format
  date?: string;
  // fajr time
  fajr?: string;
  // isha time
  isha?: string;
  // maghrib time
  maghrib?: string;
  // time of sunrise
  sunrise?: string;
  // zuhr time
  zuhr?: string;
}

export interface TimesRequest {
  // optional date in YYYY-MM-DD format, otherwise uses today
  date?: string;
  // number of days to request times for
  days?: number;
  // optional latitude used in place of location
  latitude?: number;
  // location to retrieve prayer times for.
  // this can be a specific address, city, etc
  location?: string;
  // optional longitude used in place of location
  longitude?: number;
}

export interface TimesResponse {
  // date of request
  date?: string;
  // number of days
  days?: number;
  // latitude of location
  latitude?: number;
  // location for the request
  location?: string;
  // longitude of location
  longitude?: number;
  // prayer times for the given location
  times?: PrayerTime[];
}
