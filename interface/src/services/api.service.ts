import apiController from "../controllers/api.controller";

class ApiService {
	async getUserChannel(channelId: string) {
		return await apiController.getUserChannel(channelId);
	}
}

export default new ApiService();
