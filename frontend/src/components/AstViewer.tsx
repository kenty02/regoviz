import { selectedSampleAtom } from "@/App.tsx";
import { Button } from "@/components/ui/button.tsx";
import { useGetAstSuspense } from "@/default/default.ts";
import ReactJson from "@microlink/react-json-view";
import { useAtomValue } from "jotai/index";

export function AstViewer() {
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
