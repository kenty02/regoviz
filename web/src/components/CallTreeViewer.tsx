import { policyAtom } from "@/App.tsx";
import { ReactECharts } from "@/components/React-ECharts.tsx";
import { Input } from "@/components/ui/input.tsx";
import { useGetCallTreeSuspense } from "@/default/default.ts";
import { RuleChild, RuleChildElse, RuleParent, RuleStatement } from "@/model";
import { TreeSeriesNodeItemOption } from "echarts/types/src/chart/tree/TreeSeries";
import { useAtomValue } from "jotai/index";
import { useMemo, useState } from "react";

const convertRules = (
	node: RuleParent | RuleChild | RuleChildElse,
): TreeSeriesNodeItemOption => {
	return {
		id: node.uid,
		name: node.name,
		children:
			node.type === "parent" || node.type === "child-else"
				? node.children.map((c) => convertRules(c))
				: node.statements.map((c) => convertStatements(c)),
	};
};
const convertStatements = (node: RuleStatement): TreeSeriesNodeItemOption => {
	return {
		id: node.uid,
		name: node.name,
		children: node.dependencies.map((c) => {
			if (typeof c === "string") {
				return {
					id: `dep-${c}-${node.uid}`,
					name: c,
				};
			}
			return convertRules(c);
		}),
	};
};
export function CallTreeViewer() {
	const policy = useAtomValue(policyAtom);
	if (policy === "") {
		return <></>;
	}
	const [entrypoint, setEntrypoint] = useState("allow"); //todo
	const { data } = useGetCallTreeSuspense({
		policy,
		entrypoint: entrypoint,
	});
	const first = data.data.entrypoint;
	const chartData: TreeSeriesNodeItemOption = useMemo(
		() => convertRules(first),
		[first],
	);
	return (
		<>
			<div>
				<Input
					placeholder={"Entrypoint"}
					value={entrypoint}
					onChange={(e) => setEntrypoint(e.target.value)}
				/>

				<ReactECharts
					style={{ height: "50%" }}
					option={{
						tooltip: {
							trigger: "item",
							triggerOn: "mousemove",
						},
						series: [
							{
								type: "tree",

								data: [chartData],

								top: "1%",
								left: "7%",
								bottom: "1%",
								right: "20%",

								symbolSize: 12,

								edgeShape: "polyline",

								roam: true,

								label: {
									position: "left",
									verticalAlign: "middle",
									align: "right",
								},

								leaves: {
									label: {
										position: "right",
										verticalAlign: "middle",
										align: "left",
									},
								},

								emphasis: {
									focus: "descendant",
								},

								expandAndCollapse: true,
								animationDuration: 550,
								animationDurationUpdate: 750,
							},
						],
					}}
				/>
			</div>
		</>
	);
}
