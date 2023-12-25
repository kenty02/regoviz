import { ReactNode, Suspense, useEffect, useState } from "react";

import { Fallback } from "@/components/fallback.tsx";
import { Button } from "@/components/ui/button.tsx";
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "@/components/ui/card.tsx";
import { Input } from "@/components/ui/input.tsx";
import {
	Tabs,
	TabsContent,
	TabsList,
	TabsTrigger,
} from "@/components/ui/tabs.tsx";
import { Textarea } from "@/components/ui/textarea.tsx";
import ReactJson from "@microlink/react-json-view";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { atom, useAtom, useAtomValue } from "jotai";
import { ErrorBoundary } from "react-error-boundary";
import {
	useGetAstSuspense,
	useGetDepTreeTextSuspense,
	useGetFlowchartSuspense,
	useGetSamplesSuspense,
	usePostVarTrace,
} from "./default/default.ts";
import { Sample } from "./model";

export const Heading = ({ children }: { children: ReactNode }) => {
	return (
		<div>
			<h1 className={"text-3xl text-bold"}>{children}</h1>
		</div>
	);
};

const selectedSampleAtom = atom<Sample | null>(null);
function SampleFileViewer() {
	const { data } = useGetSamplesSuspense();
	const files = data.data;
	const [selectedSample, setSelectedSample] = useAtom(selectedSampleAtom);
	useEffect(() => {
		// auto select first sample
		if (selectedSample == null && files.length > 0) {
			setSelectedSample(files[0]);
		}
	});
	const onSampleClick = (file: Sample) => {
		setSelectedSample(file);
	};

	return (
		<>
			<Heading>サンプルファイル一覧</Heading>
			<div>選択中：{selectedSample?.file_name ?? "なし"}</div>
			<div className={"outline"}>
				{files.map((file) => {
					return (
						<div
							key={file.file_name}
							onClick={() => onSampleClick(file)}
							onKeyDown={() => onSampleClick(file)}
						>
							{file.file_name}
						</div>
					);
				})}
			</div>
			{selectedSample && (
				<>
					<div>サンプルファイルの内容</div>
					<div
						className={
							"font-mono whitespace-pre-wrap bg-gray-100 p-2 w-full h-96 overflow-auto" +
							" border-2 border-gray-300 rounded-md outline-none"
						}
					>
						{selectedSample.content}
					</div>
				</>
			)}
		</>
	);
}

function AstViewer() {
	const selectedSample = useAtomValue(selectedSampleAtom);
	if (!selectedSample) {
		throw new Error("サンプルファイルを選択してください");
	}
	const { data } = useGetAstSuspense({ module: selectedSample.content });
	const astText = data.data.result;
	const ast = JSON.parse(astText);
	const onCopyClick = () => {
		void navigator.clipboard.writeText(astText);
	};

	return (
		<>
			<Button onClick={onCopyClick}>Copy</Button>
			<ReactJson src={ast} theme={"monokai"} />
		</>
	);
}
function VarTraceViewer() {
	const selectedSample = useAtomValue(selectedSampleAtom);
	if (!selectedSample) {
		throw new Error("selectedSample is null");
	}
	const [input, setInput] = useState("");
	const [data, setData] = useState("");
	const [query, setQuery] = useState("");
	const [commands, setCommands] = useState(`# ここにコマンドを入力してください
# 例：
# showVars 8 role
# fixVar 8 role "hoge"`);
	const mutation = usePostVarTrace();
	const onExecuteClick = () => {
		void mutation.mutateAsync({
			params: {
				sampleName: selectedSample.file_name,
				commands,
				input: input.length > 0 ? input : undefined,
				data: data.length > 0 ? data : undefined,
				query,
			},
		});
	};
	return (
		<>
			<div className={"whitespace-pre-wrap"}>
				{mutation.data?.data.result ?? "Press Execute to get output"}
			</div>

			<Textarea
				placeholder={"Input"}
				value={input}
				onChange={(e) => setInput(e.target.value)}
			/>
			<Textarea
				placeholder={"Data"}
				value={data}
				onChange={(e) => setData(e.target.value)}
			/>
			<Input
				placeholder={"Query"}
				value={query}
				onChange={(e) => setQuery(e.target.value)}
			/>
			<Textarea
				placeholder={"Commands"}
				value={commands}
				onChange={(e) => setCommands(e.target.value)}
			/>
			<Button onClick={onExecuteClick}>Execute</Button>
		</>
	);
}
function DepTreeViewer() {
	const selectedSample = useAtomValue(selectedSampleAtom);
	if (!selectedSample) {
		throw new Error("selectedSample is null");
	}
	const { data: depTreeData } = useGetDepTreeTextSuspense({
		sampleName: selectedSample.file_name,
	});

	return (
		<>
			<div
				className={
					"font-mono whitespace-pre-wrap" +
					" bg-gray-100 p-2 w-full" +
					" overflow-auto border-2 border-gray-300 rounded-md outline-none"
				}
			>
				{depTreeData.data.result}
			</div>
		</>
	);
}
function FlowchartViewer() {
	const selectedSample = useAtomValue(selectedSampleAtom);
	if (!selectedSample) {
		throw new Error("selectedSample is null");
	}
	const { data } = useGetFlowchartSuspense({
		sampleName: selectedSample.file_name,
	});
	return (
		<>
			<Button className={"bg-pink-400"} asChild>
				<a href={data.data.result} target={"_blank"} rel={"noreferrer"}>
					Open
				</a>
			</Button>
			<iframe
				className={"w-full h-1/2"}
				src={data.data.result}
				title={"flowchart"}
			/>
		</>
	);
}

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
					<Tabs defaultValue="varTrace" className="w-screen mx-4">
						<TabsList className="grid w-full grid-cols-4">
							<TabsTrigger value={"varTrace"}>変数トレース</TabsTrigger>
							<TabsTrigger value={"depTree"}>依存関係木</TabsTrigger>
							<TabsTrigger value={"flowchart"}>フローチャート</TabsTrigger>
							<TabsTrigger value={"ast"}>AST</TabsTrigger>
						</TabsList>
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
					</Tabs>
				</>
			)}
		</>
	);
}

export default App;
