import { policyAtom } from "@/App.tsx";
import { ReactECharts } from "@/components/React-ECharts.tsx";
import {
	useGetCallTreeAvailableEntrypointsSuspense,
	useGetCallTreeSuspense,
} from "@/default/default.ts";
import { RuleChild, RuleChildElse, RuleParent, RuleStatement } from "@/model";
import { TreeSeriesNodeItemOption } from "echarts/types/src/chart/tree/TreeSeries";
import { useAtomValue } from "jotai/index";
import { useMemo, useState } from "react";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu.tsx";
import { Button } from "@/components/ui/button.tsx";

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
		throw new Error("No policy selected");
	}
	const { data: data1 } = useGetCallTreeAvailableEntrypointsSuspense({
		policy,
	});
	const entrypoints = data1.data.entrypoints;
	if (entrypoints.length === 0) {
		throw new Error("No entrypoints found");
	}
	const [entrypoint, setEntrypoint] = useState<string | null>(null);
	if (entrypoint === null || !entrypoints.includes(entrypoint)) {
		setEntrypoint(entrypoints[0]);
		return <></>;
	}
	return (
		<>
			<div>
				<DropdownMenu>
					<DropdownMenuTrigger asChild>
						<Button variant="outline">Tree root: "{entrypoint}"</Button>
					</DropdownMenuTrigger>
					<DropdownMenuContent className="w-56">
						{entrypoints.map((entrypoint) => {
							return (
								<DropdownMenuItem
									key={entrypoint}
									onClick={() => {
										setEntrypoint(entrypoint);
									}}
								>
									{entrypoint}
								</DropdownMenuItem>
							);
						})}
					</DropdownMenuContent>
				</DropdownMenu>
				<CallTreeGraph policy={policy} entrypoint={entrypoint} />
			</div>
		</>
	);
}
const CallTreeGraph = (props: { policy: string; entrypoint: string }) => {
	const { data } = useGetCallTreeSuspense({
		policy: props.policy,
		entrypoint: props.entrypoint,
	});
	const first = data.data.entrypoint;
	const chartData: TreeSeriesNodeItemOption = useMemo(
		() => convertRules(first),
		[first],
	);
	return (
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
	);
};
