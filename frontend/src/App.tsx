import {ReactNode, useState} from "react";

import {DepGraph} from "./components/dep-graph.tsx";
import {
	useGetAstSuspense,
	useGetDepTreeTextSuspense,
	useGetFlowchartSuspense,
	useGetSamplesSuspense,
	usePostVarTrace,
} from "./default/default.ts";
import {atom, useAtom, useAtomValue} from "jotai";
import {Sample} from "./model";
import ReactJson from "@microlink/react-json-view";

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
					<div className={"outline whitespace-pre-wrap"}>{selectedSample.content}</div>
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
			<Heading>AST</Heading>
			<button
				type={"button"}
				className={"btn btn-primary"}
				onClick={onCopyClick}
			>
				Copy
			</button>
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
	const [query, setQuery] = useState("");
	const [commands, setCommands] = useState("");
	const mutation = usePostVarTrace();
	const onExecuteClick = () => {
		mutation.mutateAsync({
			params: {
				sampleName: selectedSample.file_name,
				commands,
				input,
				query,
			},
		});
	};
	return (
		<>
			<Heading>変数トレース</Heading>
			<div className={"whitespace-pre-wrap"}>{mutation.data?.data.result ?? "Press Execute to get output"}</div>

			<label>Input</label>
			<input
				type={"text"}
				value={input}
				onChange={(e) => setInput(e.target.value)}
			/>
			<label>Query</label>
			<input
				type={"text"}
				value={query}
				onChange={(e) => setQuery(e.target.value)}
			/>
			<label>Commands</label>
			<input
				type={"text"}
				value={commands}
				onChange={(e) => setCommands(e.target.value)}
			/>
			<button
				type="button"
				className={"btn btn-primary"}
				onClick={onExecuteClick}
			>
				Execute
			</button>
		</>
	);
}
function DepTreeViewer() {
	const selectedSample = useAtomValue(selectedSampleAtom);
	if (!selectedSample) {
		throw new Error("selectedSample is null");
	}
	const { data: depTreeData } = useGetDepTreeTextSuspense({
		module: selectedSample.content,
	});

	return (
		<>
			<Heading>依存関係木</Heading>
			<div className={"whitespace-pre-wrap"}>{depTreeData.data.result}</div>
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
			<Heading>フローチャート</Heading>
			<a href={data.data.result} target={"_blank"} rel={"noreferrer"}>{data.data.result}</a>
		</>
	);
}
function App() {
	const selectedSample = useAtomValue(selectedSampleAtom);
	const isSampleSelected = !!selectedSample;

	return (
		<>
			<SampleFileViewer />
			{isSampleSelected && (
				<>
					<VarTraceViewer />
					<DepTreeViewer />
					<FlowchartViewer />
					<AstViewer />
				</>
			)}

			<Heading>DepGraph Frontend Example</Heading>
			<DepGraph />
		</>
	);
}

export default App;
