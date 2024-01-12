import { AxiosError } from "axios";
import { Button } from "@/components/ui/button.tsx";

export function Fallback({
	error,
	resetErrorBoundary,
}: { error: Error; resetErrorBoundary: () => void }) {
	// Call resetErrorBoundary() to reset the error boundary and retry the render.

	let message = error.message;
	if (error instanceof AxiosError) {
		const responseData = error.response?.data;
		if (
			typeof responseData === "object" &&
			responseData.error_message !== undefined
		) {
			message = responseData.error_message;
		}
	}
	return (
		<div role="alert">
			<p>Something went wrong:</p>
			<pre style={{ color: "red" }}>{message}</pre>
			<Button onClick={resetErrorBoundary}>Try again</Button>
		</div>
	);
}
