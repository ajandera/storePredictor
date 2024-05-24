import Language from "~/model/Language";

export default interface IResponseLanguages {
  data: {
    success: boolean;
    languages: Language[];
    message: string;
    error: string;
  };
}
