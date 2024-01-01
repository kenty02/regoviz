import { selectedSampleAtom } from "@/App.tsx";
import { Heading } from "@/components/Heading.tsx";
import { useGetSamplesSuspense } from "@/default/default.ts";
import { Sample } from "@/model";
import { useAtom } from "jotai/index";
import { useEffect } from "react";

export function SampleFileViewer() {
	const { data } = useGetSamplesSuspense();
	const files = data.data;
	const [selectedSample, setSelectedSample] = useAtom(selectedSampleAtom);
	useEffect(() => {
		// auto select first sample
		if (selectedSample == null && files.length > 0) {
			setSelectedSample(files[0]);
		}
	});
	const onSampleClick = (file: Sample) => {
		setSelectedSample(file);
	};

	return (
		<>
			<Heading>サンプルファイル一覧</Heading>
			<div>選択中：{selectedSample?.file_name ?? "なし"}</div>
			<div className={"outline"}>
				{files.map((file) => {
					return (
						<div
							key={file.file_name}
							onClick={() => onSampleClick(file)}
							onKeyDown={() => onSampleClick(file)}
						>
							{file.file_name}
						</div>
					);
				})}
			</div>
			{selectedSample && (
				<>
					<div>サンプルファイルの内容</div>
					<div
						className={
							"font-mono whitespace-pre-wrap bg-gray-100 p-2 w-full h-96 overflow-auto" +
							" border-2 border-gray-300 rounded-md outline-none"
						}
					>
						{selectedSample.content}
					</div>
				</>
			)}
		</>
	);
}
