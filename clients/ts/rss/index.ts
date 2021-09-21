import * as m3o from "@m3o/m3o-node";

export class RssService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Add a new RSS feed with a name, url, and category
  add(request: AddRequest): Promise<AddResponse> {
    return this.client.call("rss", "Add", request) as Promise<AddResponse>;
  }
  // Get an RSS feed by name. If no name is given, all feeds are returned. Default limit is 25 entries.
  feed(request: FeedRequest): Promise<FeedResponse> {
    return this.client.call("rss", "Feed", request) as Promise<FeedResponse>;
  }
  // List the saved RSS fields
  list(request: ListRequest): Promise<ListResponse> {
    return this.client.call("rss", "List", request) as Promise<ListResponse>;
  }
  // Remove an RSS feed by name
  remove(request: RemoveRequest): Promise<RemoveResponse> {
    return this.client.call(
      "rss",
      "Remove",
      request
    ) as Promise<RemoveResponse>;
  }
}

export interface AddRequest {
  // category to add e.g news
  category?: string;
  // rss feed name
  // eg. a16z
  name?: string;
  // rss feed url
  // eg. http://a16z.com/feed/
  url?: string;
}

export interface AddResponse {}

export interface Entry {
  // article content
  content?: string;
  // data of the entry
  date?: string;
  // the rss feed where it came from
  feed?: string;
  // unique id of the entry
  id?: string;
  // rss feed url of the entry
  link?: string;
  // article summary
  summary?: string;
  // title of the entry
  title?: string;
}

export interface Feed {
  // category of the feed e.g news
  category?: string;
  // unique id
  id?: string;
  // rss feed name
  // eg. a16z
  name?: string;
  // rss feed url
  // eg. http://a16z.com/feed/
  url?: string;
}

export interface FeedRequest {
  // limit entries returned
  limit?: number;
  // rss feed name
  name?: string;
  // offset entries
  offset?: number;
}

export interface FeedResponse {
  entries?: Entry[];
}

export interface ListRequest {}

export interface ListResponse {
  feeds?: Feed[];
}

export interface RemoveRequest {
  // rss feed name
  // eg. a16z
  name?: string;
}

export interface RemoveResponse {}
