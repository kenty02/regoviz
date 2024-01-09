import { dataJsonAtom, inputJsonAtom, policyAtom } from "@/App.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Input } from "@/components/ui/input.tsx";
import { Textarea } from "@/components/ui/textarea.tsx";
import { usePostVarTrace } from "@/default/default.ts";
import { useAtomValue } from "jotai/index";
import { useState } from "react";

export function VarTraceViewer() {
	const policy = useAtomValue(policyAtom);
	if (policy === "") {
		return <></>;
	}
	const [input] = useAtomValue(inputJsonAtom);
	const [data] = useAtomValue(dataJsonAtom);
	const [query, setQuery] = useState("");
	const [commands, setCommands] = useState(`# ここにコマンドを入力してください
# 例：
# showVars 8 role
# fixVar 8 role "hoge"`);

	const mutation = usePostVarTrace();
	const onExecuteClick = () => {
		void mutation.mutateAsync({
			params: {
				policy,
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
