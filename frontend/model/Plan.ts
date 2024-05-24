export default interface Plan {
  Id: string;
  Price: number;
  Period: number;
  Name: string;
  Products: number;
  Enabled: boolean;
  Free: boolean;
  OneTime: boolean;
}
