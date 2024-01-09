import { policyAtom } from "@/App.tsx";
import { Button } from "@/components/ui/button.tsx";
import {
	useGetAstPrettySuspense,
	useGetAstSuspense,
} from "@/default/default.ts";
import ReactJson from "@microlink/react-json-view";
import { useAtomValue } from "jotai/index";

export function AstViewer() {
	const policy = useAtomValue(policyAtom);
	if (policy === "") {
		return <></>;
	}

	const { data } = useGetAstSuspense({ policy });
	const { data: dataPretty } = useGetAstPrettySuspense({
		policy,
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
