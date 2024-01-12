import { ComponentPreviews, useInitial } from "@/dev";
import { DevSupport } from "@react-buddy/ide-toolbox";
import * as Sentry from "@sentry/react";
import * as axios from "axios";
import React from "react";
import ReactDOM from "react-dom/client";
import { App } from "./App.tsx";
import "./index.css";

Sentry.init({
	dsn: "https://beca4591b31766452455fcee833cd232@o4504839999848448.ingest.sentry.io/4506472334426112",
	integrations: [
		// new Sentry.BrowserTracing(),
		new Sentry.Replay({
			maskAllText: false,
			blockAllMedia: false,
		}),
	],
	// Set 'tracePropagationTargets' to control for which URLs distributed tracing should be enabled
	// tracePropagationTargets: ["localhost", /^https:\/\/yourserver\.io\/api/],
	// Performance Monitoring
	tracesSampleRate: 1.0, //  Capture 100% of the transactions
	// Session Replay
	replaysSessionSampleRate: 0.1, // This sets the sample rate at 10%. You may want to change it to 100% while in development and then sample at a lower rate in production.
	replaysOnErrorSampleRate: 1.0, // If you're not already sampling the entire session, change the sample rate to 100% when sampling sessions where errors occur.
	environment: import.meta.env.VITE_SENTRY_ENVIRONMENT,
	enabled: import.meta.env.PROD,
});
async function enableMocking() {
	if (process.env.NODE_ENV !== "development-mocked") {
		return;
	}

	const { worker } = await import("./mocks/browser");

	// `worker.start()` returns a Promise that resolves
	// once the Service Worker is up and ready to intercept requests.
	return worker.start();
}

enableMocking().then(() => {
	axios.default.defaults.baseURL = import.meta.env.VITE_API_URL;
	// set auth token
	axios.default.defaults.headers.common.Authorization = `Bearer ${
		import.meta.env.VITE_API_TOKEN
	}`;
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
