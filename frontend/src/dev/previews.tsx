import { App } from "@/App.tsx";
import { ComponentPreview, Previews } from "@react-buddy/ide-toolbox";
import { PaletteTree } from "./palette";

const ComponentPreviews = () => {
	return (
		<Previews palette={<PaletteTree />}>
			<ComponentPreview path="/AppWithProviders">
				<App />
			</ComponentPreview>
		</Previews>
	);
};

export default ComponentPreviews;
