import { dataJsonAtom, inputJsonAtom, policyAtom } from "@/App.tsx";
import {
	ReactECharts,
	ReactEChartsProps,
} from "@/components/React-ECharts.tsx";
import { Button } from "@/components/ui/button.tsx";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu.tsx";
import {
	Form,
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";
import {
	useGetCallTreeAvailableEntrypointsSuspense,
	useGetCallTreeSuspense,
} from "@/default/default.ts";
import { RuleChild, RuleChildElse, RuleParent, RuleStatement } from "@/model";
import { zodResolver } from "@hookform/resolvers/zod";
import { Separator } from "@radix-ui/react-dropdown-menu";
import { TreeSeriesNodeItemOption } from "echarts/types/src/chart/tree/TreeSeries";
import { useAtomValue } from "jotai/index";
import { useCallback, useMemo, useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { ScrollArea } from "./ui/scroll-area";

const APPEND_NODE_TYPE_TO_NAME = true;

const convertRules = (
	node: RuleParent | RuleChild | RuleChildElse,
): TreeSeriesNodeItemOption => {
	let nodeType = "";
	// let nodeLabel = "";
	let pseudoChildName = "";
	if (node.type === "parent") {
		nodeType = "RuleParent";
		// nodeLabel = `クエリ${node.name}が評価されると、同じ名前のルールがOR条件下で評価され、その結果が返されます。`;
		pseudoChildName = "(OR)";
	} else if (node.type === "child") {
		nodeType = "RuleChild";
		// nodeLabel = `ルール${node.name}には、${node.statements.length}つのステートメントが有り、最後までのステートメントが真になったときにルールの値が返されます。`;
		pseudoChildName = "(AND)";
	} else if (node.type === "child-else") {
		nodeType = "RuleChildElse";
		// nodeLabel = `ルール${node.name}は、${node.children.length}つの子ルールを持ち、最初に真になった子ルールの値が返されます。`;
		pseudoChildName = "(First match)";
	} else {
		throw new Error(`Unknown type: ${(node as { type: "__invalid__" }).type}`);
	}
	return {
		id: node.uid,
		name: APPEND_NODE_TYPE_TO_NAME ? ` (${nodeType}) ${node.name}` : node.name,
		itemStyle: {
			color: node.type === "parent" ? "#ff0000" : "#00ff00",
		},
		children: [
			{
				id: `pseudo-child-of-${node.uid}`,
				name: pseudoChildName,
				itemStyle: {
					color: "#808080",
				},
				children:
					node.type === "parent" || node.type === "child-else"
						? node.children.map((c) => convertRules(c))
						: node.statements.map((c) => convertStatements(c)),
			},
		],
	};
};
const convertStatements = (node: RuleStatement): TreeSeriesNodeItemOption => {
	return {
		id: node.uid,
		name: APPEND_NODE_TYPE_TO_NAME ? ` (Statement) ${node.name}` : node.name,
		itemStyle: {
			color: "#ffff00",
		},
		children:
			node.dependencies.length > 0
				? [
						{
							id: `pseudo-child-of-${node.uid}`,
							name: "(depends on)",
							children: node.dependencies.map((c) => {
								if (typeof c === "string") {
									return {
										id: `dep-${c}-${node.uid}`,
										name: c,
									};
								}
								return convertRules(c);
							}),
						},
				  ]
				: [],
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
	const [query, setQuery] = useState("data");
	const form = useForm<z.infer<typeof formSchema>>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			query: "data",
		},
	});
	function onSubmit(values: z.infer<typeof formSchema>) {
		setQuery(values.query);
	}
	if (entrypoint === null || !entrypoints.includes(entrypoint)) {
		setEntrypoint(entrypoints[0]);
		return <></>;
	}
	return (
		<>
			<div>
				<Form {...form}>
					<form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
						<FormField
							control={form.control}
							name="query"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Query</FormLabel>
									<FormControl>
										<Input placeholder="data.example.allow" {...field} />
									</FormControl>
									{query === "data" ? (
										<FormDescription>
											Currently I am evaluating all rules. Edit query to get
											more specific steps.
										</FormDescription>
									) : null}
									<FormMessage />
								</FormItem>
							)}
						/>
						<Button type="submit">Submit Query</Button>
					</form>
				</Form>
				<Separator className="my-2" />
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
				<CallTreeGraph policy={policy} entrypoint={entrypoint} query={query} />
			</div>
		</>
	);
}

const formSchema = z.object({
	query: z.string(),
});
const CallTreeGraph = (props: {
	policy: string;
	entrypoint: string;
	query: string;
}) => {
	const data = useAtomValue(dataJsonAtom);
	const input = useAtomValue(inputJsonAtom);
	const { data: data2 } = useGetCallTreeSuspense({
		policy: props.policy,
		entrypoint: props.entrypoint,
		data: data.length > 0 ? data : undefined,
		input: input.length > 0 ? input : undefined,
		query: props.query.length > 0 ? props.query : undefined,
	});
	const first = data2.data.entrypoint;
	const steps = data2.data.steps;
	const chartData: TreeSeriesNodeItemOption = useMemo(
		() => convertRules(first),
		[first],
	);
	const [hoveredNodeUid, setHoveredNodeUid] = useState<string | null>(null);

	const onReactEChartsMouseOver = useCallback((e: { data: unknown }) => {
		// @ts-ignore
		setHoveredNodeUid(e.data.id as string);
	}, []);

	const onReactEChartsMouseOut = useCallback(() => {
		setHoveredNodeUid(null);
	}, []);

	const option = useMemo<ReactEChartsProps["option"]>(
		() => ({
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

					edgeShape: "curve",

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
		}),
		[chartData],
	);
	return (
		<>
			{props.query ? (
				steps != null && steps.length > 0 ? (
					<ScrollArea className="h-72 w-full rounded-md border">
						<div className="p-4">
							<h4 className="mb-4 text-sm font-medium leading-none">
								Steps{hoveredNodeUid != null ? " for selected node" : ""}
							</h4>
							{steps
								.filter(
									(step) =>
										hoveredNodeUid == null ||
										step.targetNodeUid === hoveredNodeUid,
								)
								.map((step) => (
									<>
										<div key={step.index} className="text-sm">
											[{step.index}] {step.message}
										</div>
										<Separator className="my-2" />
									</>
								))}
						</div>
					</ScrollArea>
				) : (
					<p className="text-sm">No steps found</p>
				)
			) : null}
			<ReactECharts
				onMouseOver={onReactEChartsMouseOver}
				onMouseOut={onReactEChartsMouseOut}
				style={{ height: "50%" }}
				option={option}
			/>
		</>
	);
};
