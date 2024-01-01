import { ReactNode } from "react";

export const Heading = ({ children }: { children: ReactNode }) => {
	return (
		<div>
			<h1 className={"text-3xl text-bold"}>{children}</h1>
		</div>
	);
};
