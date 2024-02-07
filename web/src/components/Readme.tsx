import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";
import "./Readme.css";

export function Readme() {
	return (
		<Markdown
			remarkPlugins={[[remarkGfm]]}
			className={"markdown"}
			components={{
				a: ({ node, ...props }) => (
					<a {...props} target={"_blank"} rel={"noreferrer"} />
				),
			}}
		>{`
このサイトは、[Open Policy Agent](https://www.openpolicyagent.org/)の[Rego言語](https://www.openpolicyagent.org/docs/latest/policy-language/)を様々な形式で可視化するためのものです。

複数のツールがあります。タブによって切り替えることで、それぞれのツールを使うことができます。

上のペインにRego言語によって記述されたポリシー（必要な場合は、JSON形式のInput/Dataも）を入力することで、下のペインにツールの出力が表示されます。
[The Rego Playground](https://play.openpolicyagent.org/)同様にサンプルも用意されています。右上のExamplesから選択してください。

- 今のところエラー内容が不親切です。上手くいかないケースがあったら教えてください。
- Policy入力欄の右下をドラッグすることでなんとかペインをリサイズできます。スクロールが不親切です。

## ツール


### Call Tree Viewer

Regoの評価ツリー及び評価トレース解析を[echarts](https://echarts.apache.org/en/index.html)ライブラリを用いて表示します。ルール名を入力する必要があります。ノードの多くが畳まれているのでクリックで開いてください。

### Variable Tracer

Queryとコマンドを書くと、以下が行えます

#### 変数の値を表示

print文を書くのと等価ですが、その変数を取り合える値を集合として表示してくれます。また、偽のExpressionを挿入することで、その変数が潜在的に取り得る値を表示することもできます。

gdbの\`print\`コマンドのようなものです。
	
#### 変数の値を固定

変数の値を固定することで、その変数が取り得る値を制限することができます。

gdbの\`set variable\`コマンドのようなものです。

### DepTree: 依存関係木

ポリシー間の依存関係を表示しますが、1から書き直したコールツリーの方が性能がいいです

### Flowchart Viewer

RegoのIR（中間表現）を[mermaid](https://mermaid-js.github.io/mermaid/#/)形式に変換して、[mermaid-live-editor](https://mermaid-js.github.io/mermaid-live-editor/)で表示します。

ポリシー評価のフローチャートです。ほとんどできていません

### AST Viewer

RegoのASTをJSON形式で表示します。試験的により見やすくしたものも表示できます。ただしどちらもいくつかの情報が欠落していることに注意してください。

### IR Viewer

RegoのIR(中間表現)をツリー形式で表示します。
	`}</Markdown>
	);
}
