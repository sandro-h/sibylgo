import * as vscode from 'vscode';
import { debounce } from './util';

export function activate(context: vscode.ExtensionContext) {
	context.subscriptions.push(
		vscode.commands.registerCommand('sibylgo.showPreview', () => {
			SibylPreviewPanel.createOrShow(context.extensionUri, vscode.window.activeTextEditor.document);
		})
	);

	if (vscode.window.registerWebviewPanelSerializer) {
		// Make sure we register a serializer in activation event
		vscode.window.registerWebviewPanelSerializer(SibylPreviewPanel.viewType, {
			async deserializeWebviewPanel(webviewPanel: vscode.WebviewPanel, state: any) {
				console.log(`Got state: ${state}`);
				SibylPreviewPanel.revive(webviewPanel, context.extensionUri, vscode.window.activeTextEditor.document);
			}
		});
	}
}

/**
 * Manages preview webview panels
 */
class SibylPreviewPanel {
	/**
	 * Track the currently panel. Only allow a single panel to exist at a time.
	 */
	public static currentPanel: SibylPreviewPanel | undefined;

	public static readonly viewType = 'sibylPreview';

	private readonly _panel: vscode.WebviewPanel;
	private readonly _extensionUri: vscode.Uri;
	private readonly _textDocument: vscode.TextDocument;
	private _disposables: vscode.Disposable[] = [];

	public static createOrShow(extensionUri: vscode.Uri, textDocument: vscode.TextDocument) {
		const column = vscode.ViewColumn.Two;

		// If we already have a panel, show it.
		if (SibylPreviewPanel.currentPanel) {
			SibylPreviewPanel.currentPanel._panel.reveal(column);
			return;
		}

		// Otherwise, create a new panel.
		const panel = vscode.window.createWebviewPanel(
			SibylPreviewPanel.viewType,
			'Sibyl preview',
			column || vscode.ViewColumn.One,
			{
				// Enable javascript in the webview
				enableScripts: true,

				// And restrict the webview to only loading content from our extension's `media` directory.
				localResourceRoots: [vscode.Uri.joinPath(extensionUri, 'media')]
			}
		);

		SibylPreviewPanel.currentPanel = new SibylPreviewPanel(panel, extensionUri, textDocument);
	}

	public static revive(panel: vscode.WebviewPanel, extensionUri: vscode.Uri, textDocument: vscode.TextDocument) {
		SibylPreviewPanel.currentPanel = new SibylPreviewPanel(panel, extensionUri, textDocument);
	}

	private constructor(panel: vscode.WebviewPanel, extensionUri: vscode.Uri, textDocument: vscode.TextDocument) {
		this._panel = panel;
		this._extensionUri = extensionUri;
		this._textDocument = textDocument;

		// Set the webview's initial html content
		this._updateWebview(this._panel.webview);

		// Listen for when the panel is disposed
		// This happens when the user closes the panel or when the panel is closed programatically
		this._panel.onDidDispose(() => this.dispose(), null, this._disposables);

		// Update the content based on view changes
		this._panel.onDidChangeViewState(
			e => {
				if (this._panel.visible) {
					this._updateWebview(this._panel.webview);
				}
			},
			null,
			this._disposables
		);

		// Handle messages from the webview
		this._panel.webview.onDidReceiveMessage(
			message => {
				switch (message.command) {
					case 'alert':
						vscode.window.showErrorMessage(message.text);
						return;
				}
			},
			null,
			this._disposables
		);

		const debouncedUpdateDocument = debounce(() => this.doUpdateDocument(), 250);
		vscode.workspace.onDidChangeTextDocument((e: vscode.TextDocumentChangeEvent) => {
			if (e.document === this._textDocument) {
				debouncedUpdateDocument();
			}
		},
		null,
		this._disposables);
	}

	public doUpdateDocument() {
		this._panel.webview.postMessage({ command: 'update', text: this._textDocument.getText() });
	}

	public dispose() {
		SibylPreviewPanel.currentPanel = undefined;

		// Clean up our resources
		this._panel.dispose();

		while (this._disposables.length) {
			const x = this._disposables.pop();
			if (x) {
				x.dispose();
			}
		}
	}

	private _updateWebview(webview: vscode.Webview) {
		this._panel.title = 'Sibyl Preview';
		this._panel.webview.html = this._getHtmlForWebview(webview);
		this.doUpdateDocument();
	}

	private _getHtmlForWebview(webview: vscode.Webview) {
		// Local path to main script run in the webview
		const scriptPathOnDisk = vscode.Uri.joinPath(this._extensionUri, 'media', 'main.js');

		// And the uri we use to load this script in the webview
		const scriptUri = webview.asWebviewUri(scriptPathOnDisk);

		// Local path to css styles
		const styleResetPath = vscode.Uri.joinPath(this._extensionUri, 'media', 'reset.css');
		const stylesPathMainPath = vscode.Uri.joinPath(this._extensionUri, 'media', 'vscode.css');

		// Uri to load styles into webview
		const stylesResetUri = webview.asWebviewUri(styleResetPath);
		const stylesMainUri = webview.asWebviewUri(stylesPathMainPath);

		// Use a nonce to only allow specific scripts to be run
		const nonce = getNonce();

		return `<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">

				<!--
					Use a content security policy to only allow loading images from https or from our extension directory,
					and only allow scripts that have a specific nonce.
				-->
				<meta http-equiv="Content-Security-Policy" content="default-src 'none'; style-src ${webview.cspSource}; img-src ${webview.cspSource} https:; script-src 'nonce-${nonce}';">

				<meta name="viewport" content="width=device-width, initial-scale=1.0">

				<link href="${stylesResetUri}" rel="stylesheet">
				<link href="${stylesMainUri}" rel="stylesheet">

				<title>Sibyl Preview</title>
			</head>
			<body>
				<h1 id="lines-of-code-counter">0</h1>

				<script nonce="${nonce}" src="${scriptUri}"></script>
			</body>
			</html>`;
	}
}

function getNonce() {
	let text = '';
	const possible = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
	for (let i = 0; i < 32; i++) {
		text += possible.charAt(Math.floor(Math.random() * possible.length));
	}
	return text;
}
