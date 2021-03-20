export type SibylConfig = {
    todoFileName: string
    restUrl: string
}


// via https://davidwalsh.name/javascript-debounce-function
export function debounce(func: Function, wait: number, immediate?: boolean): Function {
	var timeout;
	return function() {
		var context = this, args = arguments;
		var later = function() {
			timeout = null;
			if (!immediate) func.apply(context, args);
		};
		var callNow = immediate && !timeout;
		clearTimeout(timeout);
		timeout = setTimeout(later, wait);
		if (callNow) func.apply(context, args);
	};
};