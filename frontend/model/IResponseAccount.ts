import Account from "~/model/Account";

export default interface IResponseAccount {
  data: {
    success: boolean;
    account: Account;
    message: string;
    error: string;
  };
}
