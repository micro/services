import * as m3o from "@m3o/m3o-node";

export class CacheService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Decrement a value (if it's a number)
  decrement(request: DecrementRequest): Promise<DecrementResponse> {
    return this.client.call(
      "cache",
      "Decrement",
      request
    ) as Promise<DecrementResponse>;
  }
  // Delete a value from the cache
  delete(request: DeleteRequest): Promise<DeleteResponse> {
    return this.client.call(
      "cache",
      "Delete",
      request
    ) as Promise<DeleteResponse>;
  }
  // Get an item from the cache by key
  get(request: GetRequest): Promise<GetResponse> {
    return this.client.call("cache", "Get", request) as Promise<GetResponse>;
  }
  // Increment a value (if it's a number)
  increment(request: IncrementRequest): Promise<IncrementResponse> {
    return this.client.call(
      "cache",
      "Increment",
      request
    ) as Promise<IncrementResponse>;
  }
  // Set an item in the cache. Overwrites any existing value already set.
  set(request: SetRequest): Promise<SetResponse> {
    return this.client.call("cache", "Set", request) as Promise<SetResponse>;
  }
}

export interface DecrementRequest {
  // The key to decrement
  key?: string;
  // The amount to decrement the value by
  value?: number;
}

export interface DecrementResponse {
  // The key decremented
  key?: string;
  // The new value
  value?: number;
}

export interface DeleteRequest {
  // The key to delete
  key?: string;
}

export interface DeleteResponse {
  // Returns "ok" if successful
  status?: string;
}

export interface GetRequest {
  // The key to retrieve
  key?: string;
}

export interface GetResponse {
  // The key
  key?: string;
  // Time to live in seconds
  ttl?: number;
  // The value
  value?: string;
}

export interface IncrementRequest {
  // The key to increment
  key?: string;
  // The amount to increment the value by
  value?: number;
}

export interface IncrementResponse {
  // The key incremented
  key?: string;
  // The new value
  value?: number;
}

export interface SetRequest {
  // The key to update
  key?: string;
  // Time to live in seconds
  ttl?: number;
  // The value to set
  value?: string;
}

export interface SetResponse {
  // Returns "ok" if successful
  status?: string;
}
