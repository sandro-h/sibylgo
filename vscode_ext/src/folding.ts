import * as vscode from 'vscode';
import { SibylConfig } from './util';
import { foldTodos } from './client';

export function activate(context: vscode.ExtensionContext, cfg: SibylConfig) {
	vscode.languages.registerFoldingRangeProvider(
		{pattern: `**/${cfg.todoFileName}`},
		new SibylFoldingRangeProvider(cfg.restUrl)
	);
}

class SibylFoldingRangeProvider implements vscode.FoldingRangeProvider {

	restUrl: string;

	constructor(restUrl: string) {
		this.restUrl = restUrl
	}

	async provideFoldingRanges(document: vscode.TextDocument, context: vscode.FoldingContext, token: vscode.CancellationToken): Promise<vscode.FoldingRange[]> {
		const text = document.getText();

		const foldLines = await foldTodos(this.restUrl, text);
		return foldLines.map(line => {
			let parts = line.split('-');
			return new vscode.FoldingRange(
				parseInt(parts[0]),
				parseInt(parts[1]),
				vscode.FoldingRangeKind.Region
			);
		});
	}
}