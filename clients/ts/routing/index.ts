import * as m3o from "@m3o/m3o-node";

export class RoutingService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Turn by turn directions from a start point to an end point including maneuvers and bearings
  directions(request: DirectionsRequest): Promise<DirectionsResponse> {
    return this.client.call(
      "routing",
      "Directions",
      request
    ) as Promise<DirectionsResponse>;
  }
  // Get the eta for a route from origin to destination. The eta is an estimated time based on car routes
  eta(request: EtaRequest): Promise<EtaResponse> {
    return this.client.call("routing", "Eta", request) as Promise<EtaResponse>;
  }
  // Retrieve a route as a simple list of gps points along with total distance and estimated duration
  route(request: RouteRequest): Promise<RouteResponse> {
    return this.client.call(
      "routing",
      "Route",
      request
    ) as Promise<RouteResponse>;
  }
}

export interface Direction {
  // distance to travel in meters
  distance?: number;
  // duration to travel in seconds
  duration?: number;
  // human readable instruction
  instruction?: string;
  // intersections on route
  intersections?: Intersection[];
  // maneuver to take
  maneuver?: { [key: string]: any };
  // street name or location
  name?: string;
  // alternative reference
  reference?: string;
}

export interface DirectionsRequest {
  // The destination of the journey
  destination?: Point;
  // The staring point for the journey
  origin?: Point;
}

export interface DirectionsResponse {
  // Turn by turn directions
  directions?: Direction[];
  // Estimated distance of the route in meters
  distance?: number;
  // Estimated duration of the route in seconds
  duration?: number;
  // The waypoints on the route
  waypoints?: Waypoint[];
}

export interface EtaRequest {
  // The end point for the eta calculation
  destination?: Point;
  // The starting point for the eta calculation
  origin?: Point;
  // speed in kilometers
  speed?: number;
  // type of transport. Only "car" is supported currently.
  type?: string;
}

export interface EtaResponse {
  // eta in seconds
  duration?: number;
}

export interface Intersection {
  bearings?: number[];
  location?: Point;
}

export interface Maneuver {
  action?: string;
  bearingAfter?: number;
  bearingBefore?: number;
  direction?: string;
  location?: Point;
}

export interface Point {
  // Lat e.g 52.523219
  latitude?: number;
  // Long e.g 13.428555
  longitude?: number;
}

export interface RouteRequest {
  // Point of destination for the trip
  destination?: Point;
  // Point of origin for the trip
  origin?: Point;
}

export interface RouteResponse {
  // estimated distance in meters
  distance?: number;
  // estimated duration in seconds
  duration?: number;
  // waypoints on the route
  waypoints?: Waypoint[];
}

export interface Waypoint {
  // gps point coordinates
  location?: Point;
  // street name or related reference
  name?: string;
}
