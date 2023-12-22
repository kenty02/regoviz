export function Fallback({
	error,
	resetErrorBoundary,
}: { error: Error; resetErrorBoundary: () => void }) {
	// Call resetErrorBoundary() to reset the error boundary and retry the render.

	return (
		<div role="alert">
			<p>Something went wrong:</p>
			<pre style={{ color: "red" }}>{error.message}</pre>
		</div>
	);
}
