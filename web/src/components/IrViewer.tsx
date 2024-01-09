import { policyAtom } from "@/App.tsx";
import { Button } from "@/components/ui/button.tsx";
import { useGetIrSuspense } from "@/default/default.ts";
import { useAtomValue } from "jotai/index";

export function IrViewer() {
	const policy = useAtomValue(policyAtom);
	if (policy === "") {
		return <></>;
	}
	const { data } = useGetIrSuspense({ policy });
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
