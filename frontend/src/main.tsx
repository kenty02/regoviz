import { ComponentPreviews, useInitial } from "@/dev";
import { DevSupport } from "@react-buddy/ide-toolbox";
import * as axios from "axios";
import React from "react";
import ReactDOM from "react-dom/client";
import { App } from "./App.tsx";
import "./index.css";

async function enableMocking() {
	if (process.env.NODE_ENV !== "development") {
		return;
	}

	const { worker } = await import("./mocks/browser");

	// `worker.start()` returns a Promise that resolves
	// once the Service Worker is up and ready to intercept requests.
	return worker.start();
}

enableMocking().then(() => {
	axios.default.defaults.baseURL = "http://localhost:8080";
	// biome-ignore lint/style/noNonNullAssertion: This is a React thing, not a Biome thing.
	ReactDOM.createRoot(document.getElementById("root")!).render(
		<React.StrictMode>
			<DevSupport
				ComponentPreviews={ComponentPreviews}
				useInitialHook={useInitial}
			>
				<App />
			</DevSupport>
		</React.StrictMode>,
	);
});
