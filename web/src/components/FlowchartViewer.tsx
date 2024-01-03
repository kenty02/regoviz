import { selectedSampleAtom } from "@/App.tsx";
import { Button } from "@/components/ui/button.tsx";
import { useGetFlowchartSuspense } from "@/default/default.ts";
import { useAtomValue } from "jotai/index";

export function FlowchartViewer() {
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
