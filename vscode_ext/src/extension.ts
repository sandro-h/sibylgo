import * as vscode from 'vscode';
import * as formatting from './formatting';
import * as folding from './folding';
import * as commands from './commands';
import * as preview from './preview';
import { SibylConfig } from './util';

const cfg: SibylConfig = {
	todoFileName: 'todo.txt',
	restUrl: 'http://localhost:8082'
}

// this method is called when vs code is activated
export function activate(context: vscode.ExtensionContext) {
	formatting.activate(context, cfg);
	folding.activate(context, cfg);
	commands.activate(context, cfg);
	preview.activate(context, cfg);
}
