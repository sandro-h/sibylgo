import * as vscode from 'vscode';
const request = require('request');

// this method is called when vs code is activated
export function activate(context: vscode.ExtensionContext) {

	const catDecorationType = vscode.window.createTextEditorDecorationType({
		color: 'orange; font-weight: bold'
	});

	const momDecorationType = vscode.window.createTextEditorDecorationType({
		color: 'inherit; font-weight: bold',		
	});

	const momPriorityDecorationType = vscode.window.createTextEditorDecorationType({
		color: 'inherit; font-weight: bold',	
		border: 'solid 1px red'	
	});

	const momDoneDecorationType = vscode.window.createTextEditorDecorationType({	
		color: '#1e420f; font-weight: bold'
	});
	
	const comDoneDecorationType = vscode.window.createTextEditorDecorationType({	
		color: '#1e420f;'
	});

	const dateDecorationType = vscode.window.createTextEditorDecorationType({
		after: {
			contentIconPath: context.asAbsolutePath('cal.png')
		}
	});

	const timeDecorationType = vscode.window.createTextEditorDecorationType({
		after: {
			contentIconPath: context.asAbsolutePath('time.png')
		}
	});

	const dueStyles = [
		{textDecoration: 'none; font-weight: bold', color: '#ff0000'},
		{textDecoration: 'none; font-weight: bold', color: '#ff4040'},
		{textDecoration: 'none; font-weight: bold', color: '#ff7d7d'},
		{textDecoration: 'none; font-weight: bold', color: '#fea4a4'},
		{textDecoration: 'none; font-weight: bold', color: '#fec7c7'}
	];
	function flatCloneStyle(st, additional) {
		var cloned = {};
		for (let key in st) {
			cloned[key] = st[key];
		}
		for (let key in additional) {
			cloned[key] = additional[key];
		}
		return cloned;
	}
	const momUntilDecorationTypes = {};
	for (let i = 0; i <= 11; i += 1) {
		let styleIndex = -1;
		if (i <= 1) styleIndex = 0;
		else if (i <= 2) styleIndex = 1;
		else if (i <= 4) styleIndex = 2;
		else if (i <= 7) styleIndex = 3;
		else if (i <= 11) styleIndex = 4;

		if (styleIndex > -1) {
			momUntilDecorationTypes['mom.until' + i] = vscode.window.createTextEditorDecorationType(dueStyles[styleIndex]);
			momUntilDecorationTypes['mom.priority.until' + i] = vscode.window.createTextEditorDecorationType(
				flatCloneStyle(dueStyles[styleIndex], {border: 'solid 1px red'}));
		}
	}

	let activeEditor = null;
	
	setActiveEditor(vscode.window.activeTextEditor);
	triggerUpdateDecorations();

	vscode.window.onDidChangeActiveTextEditor(editor => {
		setActiveEditor(editor);
		// editor.document.fileName
		triggerUpdateDecorations();
	}, null, context.subscriptions);

	vscode.workspace.onDidChangeTextDocument(event => {
		if (activeEditor && event.document === activeEditor.document) {
			triggerUpdateDecorations();
		}
	}, null, context.subscriptions);

	var timeout = null;
	function triggerUpdateDecorations() {
		if (!activeEditor) return;
		if (timeout) {
			clearTimeout(timeout);
		}
		timeout = setTimeout(updateDecorations, 250);
	}

	function isTodoEditor(editor) {
		if (!editor || !editor.document) return false;
		console.log(editor.document.fileName);
		return editor.document.fileName.indexOf('todo.txt') === editor.document.fileName.length - 8;
	}

	function setActiveEditor(editor) {
		activeEditor = isTodoEditor(editor) ? editor : null;
	}

	function updateDecorations() {
		if (!activeEditor) return;

		const text = activeEditor.document.getText();
		request.post({
			headers: {'content-type' : 'text/plain'},
			url:     'http://localhost:8082/format',
			body:    new Buffer(text).toString('base64')
		  },
		  function (error, response, body) {
			  	if (error) return;
				let lines = body.split(/\r?\n/);
				const fmts = {
					'cat': {dec: catDecorationType, list: []},
					'mom': {dec: momDecorationType, list: []},
					'mom.priority': {dec: momPriorityDecorationType, list: []},
					'mom.done': {dec: momDoneDecorationType, list: []},
					'date': {dec: dateDecorationType, list: [], hoverMessage: 'Date'},
					'time': {dec: timeDecorationType, list: [], hoverMessage: 'Time'},
					'com.done': {dec: comDoneDecorationType, list: []},
				};
				for (let key in momUntilDecorationTypes) {
					fmts[key] = {dec: momUntilDecorationTypes[key], list: []};
				}
				for (let i = 0; i < lines.length; i+=1) {
					
					let parts = lines[i].split(',');
					if (parts.length !== 3) continue;
					let startPos = activeEditor.document.positionAt(parseInt(parts[0]));
					let endPos = activeEditor.document.positionAt(parseInt(parts[1]));
					let decoration = { range: new vscode.Range(startPos, endPos), hoverMessage: null };
					let fmt = fmts[parts[2]];
					if (fmt) {
						if (fmt.hoverMessage) decoration.hoverMessage = fmt.hoverMessage;
						fmt.list.push(decoration);
					}
				}
				for (let key in fmts) {
					activeEditor.setDecorations(fmts[key].dec, fmts[key].list);
				}
			});
	}
}


