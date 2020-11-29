import * as vscode from 'vscode';
import * as util from './util';
import * as request from 'request';

function handleClean(cfg: util.SibylConfig) {
	request.post(`${cfg.restUrl}/clean`,
		(error, response) => {
			if (error || response.statusCode !== 200) {
				vscode.window.showErrorMessage(`Failed to clean done todos: ${response.statusCode} ${error}`);
			}
			else {
				vscode.window.showInformationMessage('Cleaned done todos!');
			}
		}
	);
}

function handleTrash(cfg: util.SibylConfig) {
	request.post(`${cfg.restUrl}/trash`,
		(error, response) => {
			if (error || response.statusCode !== 200) {
				vscode.window.showErrorMessage(`Failed to trash done todos: ${response.statusCode} ${error}`);
			}
			else {
				vscode.window.showInformationMessage('Trashed done todos!');
			}
		}
	);
}

export function activate(context: vscode.ExtensionContext, cfg: util.SibylConfig) {

	context.subscriptions.push(vscode.commands.registerCommand('sibylgo.clean', () => handleClean(cfg)));
	context.subscriptions.push(vscode.commands.registerCommand('sibylgo.trash', () => handleTrash(cfg)));
}
