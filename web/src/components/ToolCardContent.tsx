import { Fallback } from "@/components/fallback.tsx";
import { CardContent } from "@/components/ui/card.tsx";
import { Skeleton } from "@/components/ui/skeleton.tsx";
import { type ComponentProps, Suspense } from "react";
import { ErrorBoundary } from "react-error-boundary";

export const ToolCardContent = (props: ComponentProps<typeof CardContent>) => {
	return (
		<ErrorBoundary FallbackComponent={Fallback}>
			<Suspense fallback={<Skeleton className="h-12 w-12" />}>
				<CardContent {...props} />
			</Suspense>
		</ErrorBoundary>
	);
};
