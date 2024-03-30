import type { ECOption } from "@/lib/echarts.ts";
import type { ECElementEvent, ECharts, SetOptionOpts } from "echarts";
import { getInstanceByDom, init } from "echarts";
import type { CSSProperties } from "react";
import { forwardRef, useEffect, useImperativeHandle, useRef } from "react";

export interface ReactEChartsProps {
	option: ECOption;
	style?: CSSProperties;
	settings?: SetOptionOpts;
	loading?: boolean;
	theme?: "light" | "dark";
	className?: string;

	onMouseOver?: (e: ECElementEvent) => void;
	onMouseOut?: (e: ECElementEvent) => void;
}

export interface ReactEChartsRef {
	focusNode: (nodeId: string) => void;
}

export const ReactECharts = forwardRef<ReactEChartsRef, ReactEChartsProps>(
	(
		{
			option,
			style,
			settings,
			loading,
			theme,
			className,
			onMouseOver,
			onMouseOut,
		},
		ref,
	) => {
		const chartRef = useRef<HTMLDivElement>(null);

		useImperativeHandle(
			ref,
			() => ({
				focusNode: (nodeId: string) => {
					if (chartRef.current === null) {
						return;
					}
					const chart = getInstanceByDom(chartRef.current);
					if (chart === undefined) {
						return;
					}
					chart.dispatchAction({
						type: "highlight",
						seriesIndex: 0,
						dataName: nodeId,
					});
				},
			}),
			[],
		);
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

		useEffect(() => {
			if (chartRef.current !== null) {
				const chart = getInstanceByDom(chartRef.current);
				if (chart === undefined) {
					return;
				}
				const onChartMouseOver = (e: ECElementEvent) => {
					onMouseOver?.(e);
				};
				chart.on("mouseover", onChartMouseOver);
				const onChartMouseOut = (e: ECElementEvent) => {
					onMouseOut?.(e);
				};
				chart.on("mouseout", onChartMouseOut);
				return () => {
					chart.off("mouseover", onChartMouseOver);
					chart.off("mouseout", onChartMouseOut);
				};
			}
		}, [onMouseOver, onMouseOut]);

		return (
			<div
				ref={chartRef}
				style={{
					width: "100%",
					height: "100%",
					position: "absolute",
					...style,
				}}
				className={className}
			/>
		);
	},
);
