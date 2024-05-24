import Text from "~/model/Text";

export default interface IResponseTexts {
  data: {
    success: boolean;
    texts: Text[];
    message: string;
    error: string;
  };
}
