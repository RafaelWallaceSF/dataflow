export function setTitle(title : undefined | string) {
	document.title = (typeof title === 'undefined' || title === '')
		? 'Data Flow'
		: `${title} | Data Flow`
}
