import {
	Category,
	Component,
	Palette,
	Variant,
} from "@react-buddy/ide-toolbox";
import { Fragment } from "react";

export const PaletteTree = () => (
	<Palette>
		<Category name="App">
			<Component name="Loader">
				<Variant>
					<ExampleLoaderComponent />
				</Variant>
			</Component>
		</Category>
	</Palette>
);

export function ExampleLoaderComponent() {
	// biome-ignore lint/complexity/noUselessFragments: <explanation>
	return <Fragment>Loading...</Fragment>;
}
