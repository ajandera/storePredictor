export default interface Account {
  Name: string;
  Parent: string;
  Email: string;
  Role: string;
  CountryCode: string;
  Id: string;
  Street?: string;
  City?: string;
  Zip?: string;
  CompanyNumber?: string;
  VatNumber?: string;
  Newsletter?: boolean;
  PaidTo?: string;
}
