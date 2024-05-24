import ICountry from "~/model/ICountry";

export default interface Store {
  name?: string;
  id: string;
  Id?: string;
  setting?: {
    date: string;
    currency: string;
    symbol: string;
  };
  url?: string;
  country?: string;
  feed?: string;
  window?: number;
}
