import { selectedSampleAtom } from "@/App.tsx";
import { useGetDepTreeTextSuspense } from "@/default/default.ts";
import { useAtomValue } from "jotai/index";

export function DepTreeViewer() {
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
