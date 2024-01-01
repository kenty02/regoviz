import { Suspense } from "react";

import { AstViewer } from "@/components/AstViewer.tsx";
import { CallTreeViewer } from "@/components/CallTreeViewer.tsx";
import { DepTreeViewer } from "@/components/DepTreeViewer.tsx";
import { FlowchartViewer } from "@/components/FlowchartViewer.tsx";
import { IrViewer } from "@/components/IrViewer.tsx";
import { SampleFileViewer } from "@/components/SampleFileViewer.tsx";
import { VarTraceViewer } from "@/components/VarTraceViewer.tsx";
import { Fallback } from "@/components/fallback.tsx";
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "@/components/ui/card.tsx";
import {
	Tabs,
	TabsContent,
	TabsList,
	TabsTrigger,
} from "@/components/ui/tabs.tsx";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { atom, useAtomValue } from "jotai";
import { ErrorBoundary } from "react-error-boundary";
import { Sample } from "./model";

export const selectedSampleAtom = atom<Sample | null>(null);

export function App() {
	const queryClient = new QueryClient();
	return (
		<Suspense fallback={"Loading..."}>
			<ErrorBoundary FallbackComponent={Fallback}>
				<QueryClientProvider client={queryClient}>
					<AppInner />
				</QueryClientProvider>
			</ErrorBoundary>
		</Suspense>
	);
}

function AppInner() {
	const selectedSample = useAtomValue(selectedSampleAtom);
	const isSampleSelected = !!selectedSample;

	return (
		<>
			<div className={"my-6"}>
				<SampleFileViewer />
			</div>
			{isSampleSelected && (
				<>
					<Tabs defaultValue="callTree" className="w-screen mx-4">
						<TabsList className="grid w-full grid-cols-5">
							<TabsTrigger value={"callTree"}>コールツリー(WIP)</TabsTrigger>
							<TabsTrigger value={"varTrace"}>変数トレース</TabsTrigger>
							<TabsTrigger value={"depTree"}>依存関係木</TabsTrigger>
							<TabsTrigger value={"flowchart"}>フローチャート</TabsTrigger>
							<TabsTrigger value={"ast"}>AST</TabsTrigger>
							<TabsTrigger value={"ir"}>IR</TabsTrigger>
						</TabsList>
						<TabsContent value={"callTree"}>
							<Card>
								<CardHeader>
									<CardTitle>コールツリー</CardTitle>
									<CardDescription>
										関数の呼び出し関係を木構造で表示します
									</CardDescription>
								</CardHeader>
								<CardContent>
									<CallTreeViewer />
								</CardContent>
							</Card>
						</TabsContent>
						<TabsContent value={"varTrace"}>
							<Card>
								<CardHeader>
									<CardTitle>変数トレース</CardTitle>
									<CardDescription>
										特定の変数の値を表示・固定できます
									</CardDescription>
								</CardHeader>
								<CardContent>
									<VarTraceViewer />
								</CardContent>
							</Card>
						</TabsContent>
						<TabsContent value={"depTree"}>
							<Card>
								<CardHeader>
									<CardTitle>依存関係木</CardTitle>
									<CardDescription>
										変数の依存関係を木構造で表示します
									</CardDescription>
								</CardHeader>
								<CardContent>
									<DepTreeViewer />
								</CardContent>
							</Card>
						</TabsContent>
						<TabsContent value={"flowchart"}>
							<Card>
								<CardHeader>
									<CardTitle>フローチャート</CardTitle>
									<CardDescription>
										変数の依存関係をフローチャートで表示します
									</CardDescription>
								</CardHeader>
								<CardContent>
									<FlowchartViewer />
								</CardContent>
							</Card>
						</TabsContent>
						<TabsContent value={"ast"}>
							<Card>
								<CardHeader>
									<CardTitle>AST</CardTitle>
									<CardDescription>抽象構文木を表示します</CardDescription>
								</CardHeader>
								<CardContent>
									<AstViewer />
								</CardContent>
							</Card>
						</TabsContent>
						<TabsContent value={"ir"}>
							<Card>
								<CardHeader>
									<CardTitle>IR</CardTitle>
									<CardDescription>中間表現を表示します</CardDescription>
								</CardHeader>
								<CardContent>
									<IrViewer />
								</CardContent>
							</Card>
						</TabsContent>
					</Tabs>
				</>
			)}
		</>
	);
}

export default App;
