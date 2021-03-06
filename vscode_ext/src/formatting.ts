import * as vscode from 'vscode';
import { SibylConfig } from './config';
import { formatTodos } from './client';

type FormatDefinition = {
	dec: vscode.TextEditorDecorationType;
	hoverMessage?: string;
}

type Format = FormatDefinition & {list: vscode.DecorationOptions[]};

function initFormats(context: vscode.ExtensionContext): Record<string, FormatDefinition> {
	const fmts: Record<string, FormatDefinition> = {
		'cat': {
			dec: vscode.window.createTextEditorDecorationType({
				color: 'orange; font-weight: bold'
			})
		},
		'mom': {
			dec: vscode.window.createTextEditorDecorationType({
				color: 'inherit; font-weight: bold',
			})
		},
		'mom.priority': {
			dec: vscode.window.createTextEditorDecorationType({
				color: 'inherit; font-weight: bold',
				border: 'solid 1px red'
			})
		},
		'mom.done': {
			dec: vscode.window.createTextEditorDecorationType({
				color: '#1e420f; font-weight: bold'
			})
		},
		'date': {
			dec: vscode.window.createTextEditorDecorationType({
				after: {
					contentIconPath: context.asAbsolutePath('cal.png')
				}
			}),
			hoverMessage: 'Date'
		},
		'time': {
			dec: vscode.window.createTextEditorDecorationType({
				after: {
					contentIconPath: context.asAbsolutePath('time.png')
				}
			}),
			hoverMessage: 'Time'
		},
		'id': {
			dec: vscode.window.createTextEditorDecorationType({
				color: '#3f679a',
			}),
			hoverMessage: 'ID'
		},
		'com.done': {
			dec: vscode.window.createTextEditorDecorationType({
				color: '#1e420f;'
			}),
		},
	};

	const dueStyles = [
		{textDecoration: 'none; font-weight: bold', color: '#ff0000'},
		{textDecoration: 'none; font-weight: bold', color: '#ff4040'},
		{textDecoration: 'none; font-weight: bold', color: '#ff7d7d'},
		{textDecoration: 'none; font-weight: bold', color: '#fea4a4'},
		{textDecoration: 'none; font-weight: bold', color: '#fec7c7'}
	];

	const momUntilDecorationTypes = {};
	for (let i = 0; i <= 11; i += 1) {
		let styleIndex = -1;
		if (i <= 1) styleIndex = 0;
		else if (i <= 2) styleIndex = 1;
		else if (i <= 4) styleIndex = 2;
		else if (i <= 7) styleIndex = 3;
		else if (i <= 11) styleIndex = 4;

		if (styleIndex > -1) {
			momUntilDecorationTypes[`mom.until${i}`] = vscode.window.createTextEditorDecorationType(dueStyles[styleIndex]);
			momUntilDecorationTypes[`mom.until${i}.priority`] = vscode.window.createTextEditorDecorationType({
				...dueStyles[styleIndex],
				border: 'solid 1px red'
			});
		}
	}

	for (let key in momUntilDecorationTypes) {
		fmts[key] = { dec: momUntilDecorationTypes[key] };
	}

	return fmts;
}

function parseFormatting(formattingLines: string[], formats: Record<string, FormatDefinition>, document: vscode.TextDocument): Record<string, Format> {
	const res: Record<string, Format> = {};
	for (let key in formats) {
		res[key] = {
			...formats[key],
			list: []
		};
	}

	formattingLines.forEach(line => {
		const parts = line.split(',');
		if (parts.length !== 3) return;

		const startPos = document.positionAt(parseInt(parts[0]));
		const endPos = document.positionAt(parseInt(parts[1]));
		const fmt = res[parts[2]];

		if (fmt) {
			fmt.list.push({
				range: new vscode.Range(startPos, endPos),
				hoverMessage: fmt.hoverMessage
			});
		}
	});

	return res;
}

export function activate(context: vscode.ExtensionContext, cfg: SibylConfig) {
	const formats = initFormats(context);
	let activeEditor: vscode.TextEditor|null = null;
	let timeout = null;

	setActiveEditor(vscode.window.activeTextEditor);
	triggerUpdateDecorations();

	vscode.window.onDidChangeActiveTextEditor(editor => {
		setActiveEditor(editor);
		triggerUpdateDecorations();
	}, null, context.subscriptions);

	vscode.workspace.onDidChangeTextDocument(event => {
		if (activeEditor && event.document === activeEditor.document) {
			triggerUpdateDecorations();
		}
	}, null, context.subscriptions);


	function triggerUpdateDecorations() {
		if (!activeEditor) return;
		if (timeout) {
			clearTimeout(timeout);
		}
		timeout = setTimeout(updateDecorations, 250);
	}

	function isTodoEditor(editor: vscode.TextEditor) {
		if (!editor || !editor.document) return false;
		return editor.document.fileName.indexOf(cfg.getTodoFileName()) === editor.document.fileName.length - cfg.getTodoFileName().length;
	}

	function setActiveEditor(editor: vscode.TextEditor) {
		activeEditor = isTodoEditor(editor) ? editor : null;
	}

	async function updateDecorations() {
		if (!activeEditor) return;

		const text = activeEditor.document.getText();

		let formatLines: string[];
		try {
			formatLines = await formatTodos(cfg.getRestUrl(), text);
		}
		catch (err) {
			vscode.window.showErrorMessage(`Failed format todos: ${err}`);
			return;
		}
		
		const fmts = parseFormatting(formatLines, formats, activeEditor.document);
		for (let key in fmts) {
			const fmt = fmts[key];
			// Note: it's important to also set if the list is empty, to disable old decorations on the line.
			activeEditor.setDecorations(fmt.dec, fmt.list);
		}
	}
}
