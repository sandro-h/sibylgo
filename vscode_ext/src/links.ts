import * as vscode from 'vscode';
import { SibylConfig } from './config';

export function activate(cfg: SibylConfig) {

	vscode.languages.registerDocumentLinkProvider(
		{pattern: `**/${cfg.getTodoFileName()}`},
		new SibylDocumentLinkProvider(cfg)
	);
}

class SibylDocumentLinkProvider implements vscode.DocumentLinkProvider {

	cfg: SibylConfig;

	constructor(cfg: SibylConfig) {
		this.cfg = cfg;
	}

	provideDocumentLinks(document: vscode.TextDocument, token: vscode.CancellationToken): vscode.DocumentLink[] {
		const linkDefs = this.getLinkDefs();
		if (!linkDefs.length) {
			return;
		}

		const text = document.getText();
		return linkDefs.flatMap(def => this.extractLinksForDef(def, text, document));
	}

	getLinkDefs(): SibylLinkDefinition[] {
		if (!this.cfg.getTicketPattern() || !this.cfg.getTicketUrl()) {
			return [];
		}

		return [
			{
				pattern: this.cfg.getTicketPattern(),
				url: this.cfg.getTicketUrl()
			}
		];
	}

	extractLinksForDef(def: SibylLinkDefinition, text: string, document: vscode.TextDocument): vscode.DocumentLink[] {
		const re = new RegExp(def.pattern, 'g');
		const links = [];
		let match: RegExpExecArray;
		while ((match = re.exec(text)) !== null) {
			const uri = vscode.Uri.parse(def.url.replace('$1', text.substr(match.index, match[0].length)));
			const pos = document.positionAt(match.index);
			const link = new vscode.DocumentLink(new vscode.Range(pos, pos.translate(0, match[0].length)), uri);
			links.push(link);
		}

		return links;
	}
}

interface SibylLinkDefinition {
	pattern: string;
	url: string;
}