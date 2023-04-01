export type URLInfo = {
  uid: string;
  url: string;
  createdAt: number;
  html?: string;
};

export type Statistics = {
  total: number;
  urls: URLInfo[];
};
