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
				$('#overview').empty().append(createOverviewBoard(message.preview.overview));
				calEvents = message.preview.calendar;
				$('#calendar').fullCalendar('refetchEvents');
				break;
		}
	});

	function createInstanceList(moments, listEleClassName, showEndDate) {
		return moments.map(m => {
			let text = m.name;
			if (showEndDate) {
				text += ` (${formatDate(m.end)})`;
			}
			return createMomentCell(text).addClass(listEleClassName);
		});
	}

	function createOverviewBoard(overview) {
		return overview.categories.map(createOverviewLane);
	}

	function createOverviewLane(cat) {
		const div = $('<div class="kanban-lane" />');
		if (cat.name !== '_none') {
			div.append($('<h3/>').text(cat.name));
		}
		const cols = {
			'new': {
				title: 'New',
				ele: $('<td/>')
			},
			'waiting': {
				title: 'Waiting',
				ele: $('<td/>')
			},
			'inProgress': {
				title: 'In Progress',
				ele: $('<td/>')
			}
		};

		const header = $('<tr/>').append($.map(cols, c => $('<th>').text(c.title)));

		cat.moments.forEach(m => {
			const col = cols[m.workState];
			col.ele.append(createMomentCell(m.name));
		});
		const body = $('<tr/>').append($.map(cols, c => c.ele));

		const table = $('<table class="kanban-table"></table>')
			.append(header)
			.append(body);
		return div.append(table);
	}

	function createMomentCell(text) {
		return $('<div class="moment-cell"/>')
			.prop('title', text)
			.text(text);
	}

	function formatDate(dtString) {
		return new Date(dtString).toLocaleDateString();
	}
}());