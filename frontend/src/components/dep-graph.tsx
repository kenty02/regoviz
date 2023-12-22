import { useRef, useState } from "react";
import { GraphView, SelectionT } from "react-digraph";

const GraphConfig = {
	NodeTypes: {
		empty: {
			// required to show empty nodes
			typeText: "None",
			shapeId: "#empty", // relates to the type property of a node
			shape: (
				<symbol viewBox="0 0 100 100" id="empty" key="0">
					<circle cx="50" cy="50" r="45"></circle>
				</symbol>
			),
		},
		custom: {
			// required to show empty nodes
			typeText: "Custom",
			shapeId: "#custom", // relates to the type property of a node
			shape: (
				<symbol viewBox="0 0 50 25" id="custom" key="0">
					<ellipse cx="50" cy="25" rx="50" ry="25"></ellipse>
				</symbol>
			),
		},
	},
	NodeSubtypes: {},
	EdgeTypes: {
		emptyEdge: {
			// required to show empty edges
			shapeId: "#emptyEdge",
			shape: (
				<symbol viewBox="0 0 50 50" id="emptyEdge" key="0">
					<circle cx="25" cy="25" r="8" fill="currentColor">
						{" "}
					</circle>
				</symbol>
			),
		},
	},
};

const NODE_KEY = "id"; // Allows D3 to correctly update DOM

const sample = {
	nodes: [
		{
			id: 1,
			title: "Node A",
			x: 258.3976135253906,
			y: 331.9783248901367,
			type: "empty",
		},
		{
			id: 2,
			title: "Node B",
			x: 593.9393920898438,
			y: 260.6060791015625,
			type: "empty",
		},
		{
			id: 3,
			title: "Node C",
			x: 237.5757598876953,
			y: 61.81818389892578,
			type: "custom",
		},
		{
			id: 4,
			title: "Node C",
			x: 600.5757598876953,
			y: 600.81818389892578,
			type: "custom",
		},
	],
	edges: [
		{
			source: 1,
			target: 2,
			type: "emptyEdge",
		},
		{
			source: 2,
			target: 4,
			type: "emptyEdge",
		},
	],
};
export const DepGraph = () => {
	const [graph, setGraph] = useState(sample);
	const [selected, setSelected] = useState<SelectionT | null>(null);
	const ref = useRef(null);

	const nodes = graph.nodes;
	const edges = graph.edges;

	const NodeTypes = GraphConfig.NodeTypes;
	const NodeSubtypes = GraphConfig.NodeSubtypes;
	const EdgeTypes = GraphConfig.EdgeTypes;

	return (
		<div id="graph" style={{ width: 1000, height: 1000 }}>
			<GraphView
				ref={ref}
				nodeKey={NODE_KEY}
				nodes={nodes}
				edges={edges}
				selected={selected}
				nodeTypes={NodeTypes}
				nodeSubtypes={NodeSubtypes}
				edgeTypes={EdgeTypes}
				allowMultiselect={true}
				// onSelect={onSelect}
				// onCreateNode={onCreateNode}
				// onUpdateNode={onUpdateNode}
				// onDeleteNode={onDeleteNode}
				// onCreateEdge={onCreateEdge}
				// onSwapEdge={onSwapEdge}
				// onDeleteEdge={onDeleteEdge}
			/>
		</div>
	);
};
