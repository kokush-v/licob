export enum Routes {
	ApiUrl = "http://127.0.0.1:8888",
	ApiGetChannel = "/api/channels",
	WsUrl = "ws://127.0.0.1:8888",
}

export enum SocketEvents {
	Drop = "sakura.drop",
	Captcha = "sakura.captcha",
	Pick = "sakura.pick",
	Win = "sakura.win",
	Picked = "sakura.picked",
	SendCode = "sakura.send",
}
