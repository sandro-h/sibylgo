import * as vscode from 'vscode';

export interface SibylConfig {
	getRestUrl(): string;
	getTodoFileName(): string;
	getTicketPattern(): string;
	getTicketUrl(): string;
}

export const VSCodeSibylConfig: SibylConfig = {
	getRestUrl: configGetter('restUrl'),
	getTodoFileName: configGetter('todoFileName'),
	getTicketPattern: configGetter('ticketPattern'),
	getTicketUrl: configGetter('ticketUrl'),
};

function configGetter<T>(key: string): () => T {
	return () => vscode.workspace.getConfiguration('sibyl').get<T>(key);
}
