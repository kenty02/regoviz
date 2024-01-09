import { policyAtom } from "@/App.tsx";
import { Button } from "@/components/ui/button.tsx";
import { useGetFlowchartSuspense } from "@/default/default.ts";
import { useAtomValue } from "jotai/index";

export function FlowchartViewer() {
	const policy = useAtomValue(policyAtom);
	if (policy === "") {
		return <></>;
	}
	const { data } = useGetFlowchartSuspense({
		policy,
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
