export interface Captcha {
	drop_id: string;
	codes: string[];
}

export interface User {
	id: string;
	nick: string;
	avatar: string;
}

export interface Channel {
	id: string;
	name: string;
}

export interface UserChannel {
	channel?: Channel;
	user?: User;
}

export interface Drop {
	id: string;
	currency: number;
	timestamp: number;
}

export interface Pick {
	id: string;
	drop_id: string;
	timestamp: number;
	code: string;
	since: string;
	user: User;
}

export interface Card {
	drop: Drop;
	captcha: Captcha;
	since?: string;
	picks?: Pick[];
	winner?: User;
	active: boolean;
}

export interface WinCard {
	id: string;
	timestamp: number;
	currency: number;
	since: string;
	picks: Pick[];
	winner: User;
}
