import { Pick, Card } from "../types";

export default function pickEvent(pick: Pick, cards: Card[]): Card[] {
	const card = cards.find((card) => card.drop.id === pick.drop_id);
	if (card) {
		if (card.picks) {
			if (!card.picks.find((p) => p.timestamp === pick.timestamp)) {
				card.picks.push(pick);
			}
		} else {
			card.picks = [pick];
		}
	}

	return cards.map((elem) => (elem.drop.id === card?.drop.id ? card : elem));
}
