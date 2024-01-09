import { Suspense } from "react";

import { AstViewer } from "@/components/AstViewer.tsx";
import { CallTreeViewer } from "@/components/CallTreeViewer.tsx";
import { DepTreeViewer } from "@/components/DepTreeViewer.tsx";
import { FlowchartViewer } from "@/components/FlowchartViewer.tsx";
import { IrViewer } from "@/components/IrViewer.tsx";
import { Readme } from "@/components/Readme.tsx";
import { VarTraceViewer } from "@/components/VarTraceViewer.tsx";
import { Fallback } from "@/components/fallback.tsx";
import { Button } from "@/components/ui/button.tsx";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu.tsx";
import {
	Tabs,
	TabsContent,
	TabsList,
	TabsTrigger,
} from "@/components/ui/tabs.tsx";
import { Textarea } from "@/components/ui/textarea.tsx";
import { useGetSamplesSuspense } from "@/default/default.ts";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { atom } from "jotai";
import { useAtom } from "jotai/index";
import { atomWithStorage } from "jotai/utils";
import { ErrorBoundary } from "react-error-boundary";
import { Sample } from "./model";

export const selectedSampleAtom = atom<Sample | null>(null);
export const policyAtom = atom<string>("");
export const inputJsonAtom = atom<string>("");
export const dataJsonAtom = atom<string>("");
export const selectedToolKeyAtom = atomWithStorage<string | null>(
	"selectedToolKey",
	null,
	undefined,
	{ getOnInit: true },
);

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
	const { data } = useGetSamplesSuspense();
	const sampleFiles = data.data;
	const [policy, setPolicy] = useAtom(policyAtom);
	if (policy === "") {
		setPolicy(sampleFiles[0].content);
	}
	const [inputJson, setInputJson] = useAtom(inputJsonAtom);
	if (inputJson === "") {
		setInputJson(sampleFiles[0].input_examples.default);
	}
	const [dataJson, setDataJson] = useAtom(dataJsonAtom);
	if (dataJson === "") {
		setDataJson(sampleFiles[0].data_examples.default);
	}
	const tools: { key: string; name: string; component: JSX.Element }[] = [
		{
			key: "readme",
			name: "README",
			component: <Readme />,
		},
		{
			key: "callTree",
			name: "CallTree",
			component: <CallTreeViewer />,
		},
		{
			key: "varTrace",
			name: "VarTracer",
			component: <VarTraceViewer />,
		},
		{
			key: "depTree",
			name: "DepTree",
			component: <DepTreeViewer />,
		},
		{
			key: "flowchart",
			name: "FlowChart",
			component: <FlowchartViewer />,
		},
		{
			key: "ast",
			name: "AST",
			component: <AstViewer />,
		},
		{
			key: "ir",
			name: "IR",
			component: <IrViewer />,
		},
	];
	const [selectedToolKey, setSelectedToolKey] = useAtom(selectedToolKeyAtom);
	if (
		selectedToolKey === null ||
		!tools.some((tool) => tool.key === selectedToolKey)
	) {
		void setSelectedToolKey(tools[0].key);
		return <></>;
	}
	return (
		<main key="1" className="w-full h-full flex flex-col">
			<header className="flex items-center justify-between px-4 py-2 border-b bg-gray-100 dark:bg-gray-800">
				<h1 className="text-lg font-semibold">regoviz</h1>
				<div className="flex items-center gap-4">
					<Button variant="outline">Apply</Button>
					<DropdownMenu>
						<DropdownMenuTrigger asChild>
							<Button variant="outline">Examples</Button>
						</DropdownMenuTrigger>
						<DropdownMenuContent className="w-56">
							<DropdownMenuLabel>Presets</DropdownMenuLabel>
							<DropdownMenuSeparator />
							{sampleFiles.map((file) => {
								return (
									<DropdownMenuItem
										key={file.file_name}
										onClick={() => {
											setPolicy(file.content);
											setInputJson(file.input_examples.default);
											setDataJson(file.data_examples.default);
										}}
									>
										{file.file_name}
									</DropdownMenuItem>
								);
							})}
						</DropdownMenuContent>
					</DropdownMenu>
				</div>
			</header>
			<section className="flex flex-grow overflow-hidden">
				<div className="flex flex-col w-1/2 border-r">
					<h2 className="px-4 py-2 bg-gray-100 dark:bg-gray-800">Policy</h2>
					<div className="flex-1 overflow-auto bg-white dark:bg-gray-900">
						<Textarea
							className="h-full"
							value={policy}
							onChange={(e) => setPolicy(e.target.value)}
						/>
					</div>
				</div>
				<div className="flex flex-col w-1/2">
					<h2 className="px-4 py-2 bg-gray-100 dark:bg-gray-800">Input JSON</h2>
					<div className="flex-1 overflow-auto bg-white dark:bg-gray-900">
						<Textarea
							className="h-full resize-none"
							value={inputJson}
							onChange={(e) => setInputJson(e.target.value)}
						/>
					</div>
					<h2 className="px-4 py-2 bg-gray-100 dark:bg-gray-800">Data JSON</h2>
					<div className="flex-1 overflow-auto bg-white dark:bg-gray-900">
						<Textarea
							className="h-full resize-none"
							value={dataJson}
							onChange={(e) => setDataJson(e.target.value)}
						/>
					</div>
				</div>
			</section>
			<section className="flex flex-col border-t">
				<Tabs
					className="w-full"
					defaultValue={selectedToolKey}
					onValueChange={(value) => {
						setSelectedToolKey(value);
					}}
				>
					<TabsList className="flex justify-start">
						{tools.map((tool) => {
							return (
								<TabsTrigger key={tool.key} value={tool.key}>
									{tool.name}
								</TabsTrigger>
							);
						})}
					</TabsList>

					{tools.map((tool) => {
						return (
							<TabsContent key={tool.key} value={tool.key}>
								<div className="p-4">
									<h3 className="text-lg font-semibold">{tool.name}</h3>
									<div className="mt-2">
										<ErrorBoundary FallbackComponent={Fallback}>
											<Suspense fallback={"Loading..."}>
												{tool.component}
											</Suspense>
										</ErrorBoundary>
									</div>
								</div>
							</TabsContent>
						);
					})}
				</Tabs>
			</section>
		</main>
	);
}

export default App;
