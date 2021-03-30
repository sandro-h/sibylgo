import * as vscode from 'vscode';
import * as formatting from './formatting';
import * as folding from './folding';
import * as commands from './commands';
import * as preview from './preview';
import * as links from './links';
import { VSCodeSibylConfig } from './config';

// this method is called when vs code is activated
export function activate(context: vscode.ExtensionContext) {
	formatting.activate(context, VSCodeSibylConfig);
	folding.activate(VSCodeSibylConfig);
	commands.activate(context, VSCodeSibylConfig);
	preview.activate(context, VSCodeSibylConfig);
	links.activate(VSCodeSibylConfig);
}
