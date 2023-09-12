import { UserChannel } from "../types";
import "../styles/_header.scss";
export default function Header({ user, channel }: UserChannel) {
	return (
		user &&
		channel && (
			<header className="flex justify-between items-center px-32 py-6 text-3xl">
				<h1>{channel.name}</h1>
				<div className="flex justify-around items-center gap-4 w-fit">
					<img className="w-16 rounded-full" src={user.avatar} alt="avatar" />
					<h1>{user.nick}</h1>
				</div>
			</header>
		)
	);
}
