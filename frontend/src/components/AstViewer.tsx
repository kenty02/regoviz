import { selectedSampleAtom } from "@/App.tsx";
import { Button } from "@/components/ui/button.tsx";
import {
	useGetAstPrettySuspense,
	useGetAstSuspense,
} from "@/default/default.ts";
import ReactJson from "@microlink/react-json-view";
import { useAtomValue } from "jotai/index";

export function AstViewer() {
	const selectedSample = useAtomValue(selectedSampleAtom);
	if (!selectedSample) {
		throw new Error("サンプルファイルを選択してください");
	}
	const { data } = useGetAstSuspense({ sampleName: selectedSample.file_name });
	const { data: dataPretty } = useGetAstPrettySuspense({
		sampleName: selectedSample.file_name,
	});
	const astText = data.data.result;
	const astPrettyText = dataPretty.data.result;
	const ast = JSON.parse(astText);
	const onCopyClick = () => {
		void navigator.clipboard.writeText(astText);
	};

	return (
		<>
			json version is shown below
			<div className={"font-mono whitespace-pre overflow-x-auto"}>
				{astPrettyText}
			</div>
			<Button onClick={onCopyClick}>Copy</Button>
			<ReactJson src={ast} theme={"monokai"} />
		</>
	);
}
