import { Card } from "../types";
import CardComponent from "./CardComponent";

import "../styles/_card-container.scss";
import { WebSocketHook } from "react-use-websocket/dist/lib/types";
import { useEffect, useRef } from "react";

interface CardContainerProps {
	cards: Card[];
	socket: WebSocketHook;
}

export default function CardContainer({ cards, socket }: CardContainerProps) {
	const ref = useRef<HTMLDivElement>(null);
	useEffect(() => {
		if (cards.length) {
			ref.current?.scrollIntoView({
				behavior: "smooth",
				block: "end",
			});
		}
	}, [cards.length]);

	return (
		<div className="card-container ">
			<div className="wrapper">
				{cards.map((card) => {
					return (
						<CardComponent key={card.drop.id} card={card} socket={socket}></CardComponent>
					);
				})}
				<div ref={ref}></div>
			</div>
		</div>
	);
}
