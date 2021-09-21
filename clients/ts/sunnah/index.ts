import * as m3o from "@m3o/m3o-node";

export class SunnahService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Get a list of books from within a collection. A book can contain many chapters
  // each with its own hadiths.
  books(request: BooksRequest): Promise<BooksResponse> {
    return this.client.call(
      "sunnah",
      "Books",
      request
    ) as Promise<BooksResponse>;
  }
  // Get all the chapters of a given book within a collection.
  chapters(request: ChaptersRequest): Promise<ChaptersResponse> {
    return this.client.call(
      "sunnah",
      "Chapters",
      request
    ) as Promise<ChaptersResponse>;
  }
  // Get a list of available collections. A collection is
  // a compilation of hadiths collected and written by an author.
  collections(request: CollectionsRequest): Promise<CollectionsResponse> {
    return this.client.call(
      "sunnah",
      "Collections",
      request
    ) as Promise<CollectionsResponse>;
  }
  // Hadiths returns a list of hadiths and their corresponding text for a
  // given book within a collection.
  hadiths(request: HadithsRequest): Promise<HadithsResponse> {
    return this.client.call(
      "sunnah",
      "Hadiths",
      request
    ) as Promise<HadithsResponse>;
  }
}

export interface Book {
  // arabic name of the book
  arabicName?: string;
  // number of hadiths in the book
  hadiths?: number;
  // number of the book e.g 1
  id?: number;
  // name of the book
  name?: string;
}

export interface BooksRequest {
  // Name of the collection
  collection?: string;
  // Limit the number of books returned
  limit?: number;
  // The page in the pagination
  page?: number;
}

export interface BooksResponse {
  // A list of books
  books?: Book[];
  // Name of the collection
  collection?: string;
  // The limit specified
  limit?: number;
  // The page requested
  page?: number;
  // The total overall books
  total?: number;
}

export interface Chapter {
  // arabic title
  arabicTitle?: string;
  // the book number
  book?: number;
  // the chapter id e.g 1
  id?: number;
  // the chapter key e.g 1.00
  key?: string;
  // title of the chapter
  title?: string;
}

export interface ChaptersRequest {
  // number of the book
  book?: number;
  // name of the collection
  collection?: string;
  // Limit the number of chapters returned
  limit?: number;
  // The page in the pagination
  page?: number;
}

export interface ChaptersResponse {
  // number of the book
  book?: number;
  // The chapters of the book
  chapters?: Chapter[];
  // name of the collection
  collection?: string;
  // Limit the number of chapters returned
  limit?: number;
  // The page in the pagination
  page?: number;
  // Total chapters in the book
  total?: number;
}

export interface Collection {
  // Arabic title if available
  arabicTitle?: string;
  // Total hadiths in the collection
  hadiths?: number;
  // Name of the collection e.g bukhari
  name?: string;
  // An introduction explaining the collection
  summary?: string;
  // Title of the collection e.g Sahih al-Bukhari
  title?: string;
}

export interface CollectionsRequest {
  // Number of collections to limit to
  limit?: number;
  // The page in the pagination
  page?: number;
}

export interface CollectionsResponse {
  collections?: Collection[];
}

export interface Hadith {
  // the arabic chapter title
  arabicChapterTitle?: string;
  // the arabic text
  arabicText?: string;
  // the chapter id
  chapter?: number;
  // the chapter key
  chapterKey?: string;
  // the chapter title
  chapterTitle?: string;
  // hadith id
  id?: number;
  // hadith text
  text?: string;
}

export interface HadithsRequest {
  // number of the book
  book?: number;
  // name of the collection
  collection?: string;
  // Limit the number of hadiths
  limit?: number;
  // The page in the pagination
  page?: number;
}

export interface HadithsResponse {
  // number of the book
  book?: number;
  // name of the collection
  collection?: string;
  // The hadiths of the book
  hadiths?: Hadith[];
  // Limit the number of hadiths returned
  limit?: number;
  // The page in the pagination
  page?: number;
  // Total hadiths in the  book
  total?: number;
}
