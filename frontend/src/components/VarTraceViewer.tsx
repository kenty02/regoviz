import { selectedSampleAtom } from "@/App.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Input } from "@/components/ui/input.tsx";
import { Textarea } from "@/components/ui/textarea.tsx";
import { usePostVarTrace } from "@/default/default.ts";
import { useAtomValue } from "jotai/index";
import { useEffect, useState } from "react";

export function VarTraceViewer() {
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

	// if seslectedSample is changed, reset input, data, query
	useEffect(() => {
		setInput(selectedSample.default_inputs.default);
		setData(selectedSample.default_data.default);
		setQuery(selectedSample.default_queries.default);
	}, [selectedSample]);

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
		<div>
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
		</div>
	);
}
