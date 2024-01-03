import { TreeChart } from "echarts/charts";
import type {
	// The series option types are defined with the SeriesOption suffix
	TreeSeriesOption,
} from "echarts/charts";
import {
	// Dataset
	DatasetComponent,
	GridComponent,
	TitleComponent,
	TooltipComponent,
	// Built-in transform (filter, sort)
	TransformComponent,
} from "echarts/components";
import type {
	DatasetComponentOption,
	GridComponentOption,
	// The component option types are defined with the ComponentOption suffix
	TitleComponentOption,
	TooltipComponentOption,
} from "echarts/components";
import * as echarts from "echarts/core";
import type { ComposeOption } from "echarts/core";
import { LabelLayout, UniversalTransition } from "echarts/features";
import { CanvasRenderer } from "echarts/renderers";

// Create an Option type with only the required components and charts via ComposeOption
export type ECOption = ComposeOption<
	| TreeSeriesOption
	| TitleComponentOption
	| TooltipComponentOption
	| GridComponentOption
	| DatasetComponentOption
>;

// Register the required components
echarts.use([
	TitleComponent,
	TooltipComponent,
	GridComponent,
	DatasetComponent,
	TransformComponent,
	TreeChart,
	LabelLayout,
	UniversalTransition,
	CanvasRenderer,
]);
