import * as m3o from "@m3o/m3o-node";

export class WeatherService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Get the weather forecast for the next 1-10 days
  forecast(request: ForecastRequest): Promise<ForecastResponse> {
    return this.client.call(
      "weather",
      "Forecast",
      request
    ) as Promise<ForecastResponse>;
  }
  // Get the current weather report for a location by postcode, city, zip code, ip address
  now(request: NowRequest): Promise<NowResponse> {
    return this.client.call("weather", "Now", request) as Promise<NowResponse>;
  }
}

export interface Forecast {
  // the average temp in celsius
  avgTempC?: number;
  // the average temp in fahrenheit
  avgTempF?: number;
  // chance of rain (percentage)
  chanceOfRain?: number;
  // forecast condition
  condition?: string;
  // date of the forecast
  date?: string;
  // the URL of forecast condition icon. Simply prefix with either http or https to use it
  iconUrl?: string;
  // max temp in celsius
  maxTempC?: number;
  // max temp in fahrenheit
  maxTempF?: number;
  // minimum temp in celsius
  minTempC?: number;
  // minimum temp in fahrenheit
  minTempF?: number;
  // time of sunrise
  sunrise?: string;
  // time of sunset
  sunset?: string;
  // will it rain
  willItRain?: boolean;
}

export interface ForecastRequest {
  // number of days. default 1, max 10
  days?: number;
  // location of the forecase
  location?: string;
}

export interface ForecastResponse {
  // country of the request
  country?: string;
  // forecast for the next number of days
  forecast?: { [key: string]: any }[];
  // e.g 37.55
  latitude?: number;
  // the local time
  localTime?: string;
  // location of the request
  location?: string;
  // e.g -77.46
  longitude?: number;
  // region related to the location
  region?: string;
  // timezone of the location
  timezone?: string;
}

export interface NowRequest {
  // location to get weather e.g postcode, city
  location?: string;
}

export interface NowResponse {
  // cloud cover percentage
  cloud?: number;
  // the weather condition
  condition?: string;
  // country of the request
  country?: string;
  // whether its daytime
  daytime?: boolean;
  // feels like in celsius
  feelsLikeC?: number;
  // feels like in fahrenheit
  feelsLikeF?: number;
  // the humidity percentage
  humidity?: number;
  // the URL of the related icon. Simply prefix with either http or https to use it
  iconUrl?: string;
  // e.g 37.55
  latitude?: number;
  // the local time
  localTime?: string;
  // location of the request
  location?: string;
  // e.g -77.46
  longitude?: number;
  // region related to the location
  region?: string;
  // temperature in celsius
  tempC?: number;
  // temperature in fahrenheit
  tempF?: number;
  // timezone of the location
  timezone?: string;
  // wind degree
  windDegree?: number;
  // wind direction
  windDirection?: string;
  // wind in kph
  windKph?: number;
  // wind in mph
  windMph?: number;
}
