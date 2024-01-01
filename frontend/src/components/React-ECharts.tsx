import { ECOption } from "@/lib/echarts.ts";
import { getInstanceByDom, init } from "echarts";
import type { ECharts, SetOptionOpts } from "echarts";
import { JSX, useEffect, useRef } from "react";
import type { CSSProperties } from "react";

export interface ReactEChartsProps {
	option: ECOption;
	style?: CSSProperties;
	settings?: SetOptionOpts;
	loading?: boolean;
	theme?: "light" | "dark";
}

export function ReactECharts({
	option,
	style,
	settings,
	loading,
	theme,
}: ReactEChartsProps): JSX.Element {
	const chartRef = useRef<HTMLDivElement>(null);

	useEffect(() => {
		// Initialize chart
		let chart: ECharts | undefined;
		if (chartRef.current !== null) {
			chart = init(chartRef.current, theme);
		}

		// Add chart resize listener
		// ResizeObserver is leading to a bit janky UX
		function resizeChart() {
			chart?.resize();
		}
		window.addEventListener("resize", resizeChart);

		// Return cleanup function
		return () => {
			chart?.dispose();
			window.removeEventListener("resize", resizeChart);
		};
	}, [theme]);

	// biome-ignore lint/correctness/useExhaustiveDependencies: Whenever theme changes we need to add option and setting due to it being deleted in cleanup function
	useEffect(() => {
		// Update chart
		if (chartRef.current !== null) {
			const chart = getInstanceByDom(chartRef.current);
			chart?.setOption(option, settings);
		}
	}, [option, settings, theme]);

	// biome-ignore lint/correctness/useExhaustiveDependencies: Same as above
	useEffect(() => {
		// Update chart
		if (chartRef.current !== null) {
			const chart = getInstanceByDom(chartRef.current);
			// eslint-disable-next-line @typescript-eslint/no-unused-expressions
			loading === true ? chart?.showLoading() : chart?.hideLoading();
		}
	}, [loading, theme]);

	return (
		<div ref={chartRef} style={{ width: "100%", height: "100px", ...style }} />
	);
}
