import { selectedSampleAtom } from "@/App.tsx";
import { Button } from "@/components/ui/button.tsx";
import { useGetIrSuspense } from "@/default/default.ts";
import { useAtomValue } from "jotai/index";

export function IrViewer() {
	const selectedSample = useAtomValue(selectedSampleAtom);
	if (!selectedSample) {
		throw new Error("サンプルファイルを選択してください");
	}
	const { data } = useGetIrSuspense({ sampleName: selectedSample.file_name });
	const irText = data.data.result;
	const onCopyClick = () => {
		void navigator.clipboard.writeText(irText);
	};

	return (
		<>
			<Button onClick={onCopyClick}>Copy</Button>
			<div className={"font-mono whitespace-pre-wrap"}>{irText}</div>
		</>
	);
}
