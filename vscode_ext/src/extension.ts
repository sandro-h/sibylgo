import * as vscode from 'vscode';
import * as formatting from './formatting';
import * as folding from './folding';
import * as util from './util';
const request = require('request');

// this method is called when vs code is activated
export function activate(context: vscode.ExtensionContext) {
	const cfg: util.SibylConfig = {
		todoFileName: 'todo.txt',
		restUrl: 'http://localhost:8082'
	}

	formatting.activate(context, cfg);
	folding.activate(context, cfg);
}
