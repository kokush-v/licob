import { useEffect, useState } from "react";
import CardContainer from "./components/CardContainer";
import Header from "./components/Header";
import { Card, UserChannel } from "./types";
import { useParams } from "react-router-dom";
import apiService from "./services/api.service";
import { AxiosResponse } from "axios";
import { useWebSocket } from "react-use-websocket/dist/lib/use-websocket";
import { SocketEvents } from "./enum";
import captchaEvent from "./ws-events/captcha.event";
import pickEvent from "./ws-events/pick.event";
import winEvent from "./ws-events/win.event";

function App() {
	const [userChannel, setUserChannel] = useState<UserChannel>();
	const { channel_id } = useParams();
	const [cards, setCards] = useState<Card[]>([]);

	const socket = useWebSocket(`ws://${window.location.host}/ws?channel_id=` + channel_id, {
		onOpen: () => {
			// console.log({ ws: "OK" });
		},
		onMessage: (event: MessageEvent) => {
			const jsonData = JSON.parse(event.data);

			switch (jsonData.t) {
				case SocketEvents.Drop:
					const card = {
						drop: jsonData.d,
						active: true,
					} as Card;

					if (!cards.find((findCard) => findCard.drop.id === card.drop.id))
						setCards([...cards, card]);
					break;
				case SocketEvents.Captcha:
					setCards(captchaEvent(jsonData.d, cards));
					break;
				case SocketEvents.Pick:
					setCards(pickEvent(jsonData.d, cards));
					break;
				case SocketEvents.Win:
					setCards(winEvent(jsonData.d, cards));
					break;
				case SocketEvents.Picked:
					setCards(
						cards.map((elem) => {
							if (elem.drop.id === jsonData.d.drop_id) {
								return { ...elem, active: false };
							} else {
								return elem;
							}
						})
					);
					break;
			}
		},
	});

	useEffect(() => {
		const getData = async () => {
			if (channel_id) return apiService.getUserChannel(channel_id);
		};

		getData().then(({ data }: AxiosResponse) => {
			setUserChannel(data);
		});
	}, []);

	return (
		<>
			<Header user={userChannel?.user} channel={userChannel?.channel}></Header>
			{socket?.readyState !== 0 ? (
				cards.length > 0 ? (
					<CardContainer cards={cards} socket={socket}></CardContainer>
				) : (
					<h1 className="text-3xl text-center mt-20">Пусто &#128532;</h1>
				)
			) : (
				<h1 className="text-3xl text-center mt-20">Настраиваем подключение...</h1>
			)}
		</>
	);
}

export default App;
