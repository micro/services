import * as m3o from "@m3o/m3o-node";

export class ThumbnailService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Create a thumbnail screenshot by passing in a url, height and width
  screenshot(request: ScreenshotRequest): Promise<ScreenshotResponse> {
    return this.client.call(
      "thumbnail",
      "Screenshot",
      request
    ) as Promise<ScreenshotResponse>;
  }
}

export interface ScreenshotRequest {
  // height of the browser window, optional
  height?: number;
  url?: string;
  // width of the browser window. optional
  width?: number;
}

export interface ScreenshotResponse {
  imageURL?: string;
}
