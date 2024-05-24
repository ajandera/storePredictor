import axios from 'axios' 

export default class MattermostService {
  static async sendToMarketingChannel(api: string, message: string) {
    const response = await axios
      .post(api + "mattermost", {
        text: message
      });
    if (response.data.success === true) {
      return true;
    } else {
      return false;
    }
  }
}
