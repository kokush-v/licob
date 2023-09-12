import { WebSocketHook } from "react-use-websocket/dist/lib/types";
import { SocketEvents } from "../enum";

export default function sendCode(socket: WebSocketHook, code: string) {
	const socketMessage = { t: SocketEvents.SendCode, d: code };

	socket.sendJsonMessage(socketMessage);
}
