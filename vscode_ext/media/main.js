/* eslint-disable */
// This script will be run within the webview itself
// It cannot access the main VS Code APIs directly.
(function () {
	const vscode = acquireVsCodeApi();

	// const oldState = vscode.getState();

	// const counter = document.getElementById('lines-of-code-counter');
	// console.log(oldState);
	// let currentCount = (oldState && oldState.count) || 0;
	// counter.textContent = currentCount;

	const todayList = document.getElementById('due-today');
	const weekList = document.getElementById('due-week');
	const overviewList = document.getElementById('overview');

	let calEvents = [];
	function calendarEvents(start, end, timezone, callback) {
		callback(calEvents);
	}

	$('#calendar').fullCalendar({
		header: {
			left: '',
			center: '',
			right: ''
		},
		defaultView: 'basicWeek',
		events: calendarEvents,
		firstDay: 1,
		height: 150
	});

	// setInterval(() => {
	// 	counter.textContent = currentCount++;

	// 	// Update state
	// 	vscode.setState({ count: currentCount });

	// 	// Alert the extension when the cat introduces a bug
	// 	if (Math.random() < Math.min(0.001 * currentCount, 0.05)) {
	// 		// Send a message back to the extension
	// 		vscode.postMessage({
	// 			command: 'alert',
	// 			text: '🐛  on line ' + currentCount
	// 		});
	// 	}
	// }, 100);



	// Handle messages sent from the extension to the webview
	window.addEventListener('message', event => {
		const message = event.data; // The json data that the extension sent
		switch (message.command) {
			case 'update':
				updateInstanceList(todayList, message.preview.today, 'due-today');
				updateInstanceList(weekList, message.preview.week, 'due-week', true);
				updateOverviewList(overviewList, message.preview.overview);
				calEvents = message.preview.calendar;
				$('#calendar').fullCalendar('refetchEvents');
				break;
		}
	});

	function updateInstanceList(list, moments, listEleClassName, showEndDate) {
		list.innerHTML = '';
		moments.forEach(mom => {
			var text = mom.name;
			if (showEndDate) {
				text += ' (' + formatDate(mom.end) + ')';
			}
			addListEle(list, text, listEleClassName)
		});
	}

	function updateOverviewList(parent, overview) {
		parent.innerHTML = '';
		overview.categories.forEach(cat => {
			var catDiv = document.createElement("DIV");
			catDiv.appendChild(createEleWithText("H2", cat.name));

			var momList = document.createElement("UL"); 
			cat.moments.forEach(mom => addListEle(momList, mom.name, 'moment'));
			catDiv.appendChild(momList);

			parent.appendChild(catDiv);
		});
	}

	function addListEle(list, eleString, className) {
		var node = createEleWithText("LI", eleString);
		if (className) {
			node.className = className;
		}
		list.appendChild(node);
	}

	function createEleWithText(tag, text) {
		var ele = document.createElement(tag);
		var textnode = document.createTextNode(text);
		ele.appendChild(textnode);
		return ele;
	}

	function formatDate(dtString) {
		return new Date(dtString).toLocaleDateString();
	}
}());