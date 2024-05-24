export default interface AccountOrders {
  Id: string;
  AccountRefer: string;
  StoreRefer: string;
  PlanRefer: string;
  Amount: number;
  Paid: boolean;
  Name: string;
  Email: string;
  Street: string;
  City: string;
  Zip: string;
  CountryCode: string;
  CompanyNumber: string;
  VatNumber: string;
}
