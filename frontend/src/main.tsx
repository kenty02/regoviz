import React, { Suspense } from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";
import "./index.css";
import { ErrorBoundary } from "react-error-boundary";
import { Fallback } from "./components/fallback.tsx";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import * as axios from "axios";

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
	const queryClient = new QueryClient();
	axios.default.defaults.baseURL = "http://localhost:8080";
	// biome-ignore lint/style/noNonNullAssertion: This is a React thing, not a Biome thing.
	ReactDOM.createRoot(document.getElementById("root")!).render(
		<React.StrictMode>
			<Suspense fallback={"Loading..."}>
				<ErrorBoundary FallbackComponent={Fallback}>
					<QueryClientProvider client={queryClient}>
						<App />
					</QueryClientProvider>
				</ErrorBoundary>
			</Suspense>
		</React.StrictMode>,
	);
});
