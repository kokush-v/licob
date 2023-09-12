import axios from "axios";
import { UserChannel } from "../types";
import { Routes } from "../enum";

class ApiController {
	async getUserChannel(chanelId: string): Promise<UserChannel | any> {
		try {
			const result = (await axios.get(
				Routes.ApiUrl + Routes.ApiGetChannel + `/${chanelId}`
			)) as UserChannel;
			return result;
		} catch (e) {
			return { error: e };
		}
	}
}

export default new ApiController();
