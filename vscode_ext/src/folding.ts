import * as vscode from 'vscode';
import * as util from './util';
const request = require('request');

export function activate(context: vscode.ExtensionContext, cfg: util.SibylConfig) {
	vscode.languages.registerFoldingRangeProvider(
		{pattern: `**/${cfg.todoFileName}`},
		new MyFoldingRangeProvider(cfg.restUrl)
	);	
}

class MyFoldingRangeProvider implements vscode.FoldingRangeProvider {

	restUrl: string;

	constructor(restUrl: string) {
		this.restUrl = restUrl
	}

	async provideFoldingRanges(document: vscode.TextDocument, context: vscode.FoldingContext, token: vscode.CancellationToken): Promise<vscode.FoldingRange[]> {
		const text = document.getText();
		return new Promise<vscode.FoldingRange[]>((resolve, reject) => {
			request.post({
				headers: {'content-type' : 'text/plain'},
				url:     `${this.restUrl}/folding`,
				body:    new Buffer(text).toString('base64')
			  },
			  function (error, response, body) {
					if (error) {
						reject(error);
						return;
					}
					if (!response || response.statusCode !== 200) {
						reject()
						return;
					}

					var ranges: vscode.FoldingRange[] = [];
					let lines = body.split(/\r?\n/);				
					for (let i = 0; i < lines.length; i+=1) {
						let parts = lines[i].split('-');
						var foldingRange = new vscode.FoldingRange(
							parseInt(parts[0]), parseInt(parts[1]), 
							vscode.FoldingRangeKind.Region)
						ranges.push(foldingRange);  
					}
											
					resolve(ranges);
			  });
		});
	}
}