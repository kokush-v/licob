import { Captcha, Card } from "../types";

export default function captchaEvent(captcha: Captcha, cards: Card[]): Card[] {
	const card = cards.find((card) => card.drop.id === captcha.drop_id);
	if (card) card.captcha = captcha;

	return cards.map((elem) => (elem.drop.id === card?.drop.id ? card : elem));
}
