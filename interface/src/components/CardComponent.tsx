import { Card } from "../types";

import "../styles/_card.scss";
import sendCode from "../ws-events/send.code.event";
import { WebSocketHook } from "react-use-websocket/dist/lib/types";

interface CardComponentProps {
	card: Card;
	socket: WebSocketHook;
}

function formatDate(timestamp: number) {
	const date = new Date(timestamp * 1000);

	const options: Intl.DateTimeFormatOptions = {
		year: "numeric",
		month: "numeric",
		day: "numeric",
		hour: "numeric",
		minute: "numeric",
		second: "numeric",
		hour12: false, // Use 24-hour format
	};

	return date.toLocaleDateString("en-US", options);
}

export default function CardComponent({ card, socket }: CardComponentProps) {
	return (
		<div className="card">
			<div className="title flex justify-between items-center text-2xl">
				<h2>{card.drop.currency}🌸 появились!</h2>
				<h3 className="text-xl">{formatDate(card.drop.timestamp)}</h3>
			</div>
			<div className="content">
				<div className="w-3/5 picks">
					<ul>
						{card.picks ? (
							card.picks.map((pick) => {
								return (
									<li className="flex justify-between" key={pick.timestamp}>
										<p className="w-2/4 text-left">{pick.user.nick}</p>
										<p className="w-1/4 text-center">{pick.code}</p>
										<p className="w-1/4 text-center">{pick.since}</p>
									</li>
								);
							})
						) : (
							<h1 className="text-center mt-20">No picks yet</h1>
						)}
					</ul>
				</div>
				<div className="text-2xl mt-5">
					{card.active
						? card.captcha && (
								<div className="flex justify-around gap-10">
									{card.captcha.codes.map((code) => {
										return (
											<input
												key={code}
												type="button"
												value={code}
												onClick={(e) => {
													sendCode(socket, e.currentTarget.value);
												}}
											/>
										);
									})}
								</div>
						  )
						: card.winner && (
								<div className="winner flex justify-between gap-4">
									<p>{card.winner.nick}</p>
									<p>{card.since}</p>
								</div>
						  )}
				</div>
			</div>
		</div>
	);
}
