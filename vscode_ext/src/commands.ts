import * as vscode from 'vscode';
import { SibylConfig } from './config';
import { cleanTodos, trashTodos } from './client';

async function handleClean(cfg: SibylConfig) {
	try {
		await cleanTodos(cfg.getRestUrl());
		vscode.window.showInformationMessage('Cleaned done todos!');
	}
	catch (err) {
		vscode.window.showErrorMessage(`Failed to clean done todos: ${err}`);
	}
}

async function handleTrash(cfg: SibylConfig) {
	try {
		await trashTodos(cfg.getRestUrl());
		vscode.window.showInformationMessage('Trashed done todos!');
	}
	catch (err) {
		vscode.window.showErrorMessage(`Failed to trash done todos: ${err}`);
	}
}

export function activate(context: vscode.ExtensionContext, cfg: SibylConfig) {
	context.subscriptions.push(vscode.commands.registerCommand('sibylgo.clean', () => handleClean(cfg)));
	context.subscriptions.push(vscode.commands.registerCommand('sibylgo.trash', () => handleTrash(cfg)));
}
