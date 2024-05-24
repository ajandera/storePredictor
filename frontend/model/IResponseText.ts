import Text from "~/model/Text";

export default interface IResponseText {
  data: {
    success: boolean;
    text: Text;
    message: string;
    error: string;
  };
}
