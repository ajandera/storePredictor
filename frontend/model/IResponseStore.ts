export default interface IResponseStore {
  data: {
    success: boolean;
    stores: any;
    message: string;
    error: string;
  };
}
