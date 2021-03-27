/* eslint-disable */
// This script will be run within the webview itself
// It cannot access the main VS Code APIs directly.
(function () {
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

	// Handle messages sent from the extension to the webview
	window.addEventListener('message', event => {
		const message = event.data; // The json data that the extension sent
		switch (message.command) {
			case 'update':
				$('#due-today').empty().append(createInstanceList(message.preview.today, 'due-today'));
				$('#due-week').empty().append(createInstanceList(message.preview.week, 'due-week', true));
				$('#overview').empty().append(createOverviewList(message.preview.overview));
				calEvents = message.preview.calendar;
				$('#calendar').fullCalendar('refetchEvents');
				break;
		}
	});

	function createInstanceList(moments, listEleClassName, showEndDate) {
		const eles = moments.map(mom => {
			let text = mom.name;
			if (showEndDate) {
				text += ` (${formatDate(mom.end)})`;
			}
			return createListEle(text, listEleClassName);
		});

		return $('<ul/>').append(eles);
	}

	function createOverviewList(overview) {
		const eles = overview.categories.map(cat => $('<div/>')
			.append($(`<h2>${cat.name === '_none' ? 'No category' : cat.name}</h2>`))
			.append(createMomentList(cat.moments)));

		return $('<ul/>').append(eles);
	}

	function createMomentList(moments) {
		const eles = moments.map(mom => createListEle(mom.name, 'moment'));
		return $('<ul/>').append(eles);
	}

	function createListEle(text, className) {
		var ele = $('<li/>');
		ele.text(text);
		if (className) {
			ele.addClass(className);
		}
		return ele;
	}

	function formatDate(dtString) {
		return new Date(dtString).toLocaleDateString();
	}
}());