import { policyAtom } from "@/App.tsx";
import { useGetDepTreeTextSuspense } from "@/default/default.ts";
import { useAtomValue } from "jotai/index";

export function DepTreeViewer() {
	const policy = useAtomValue(policyAtom);
	if (policy === "") {
		return <></>;
	}
	const { data: depTreeData } = useGetDepTreeTextSuspense({
		policy,
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
