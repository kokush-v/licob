import { Card, WinCard } from "../types";

export default function winEvent(win: WinCard, cards: Card[]): Card[] {
	return cards.map((elem) => {
		if (elem.drop.id === win.id) {
			const card: Card = {
				drop: {
					id: win.id,
					currency: win.currency,
					timestamp: win.timestamp,
				},
				captcha: elem.captcha,
				since: win.since,
				picks: win.picks,
				winner: win.winner,
				active: false,
			};
			return card;
		} else {
			return elem;
		}
	});
}
