open-media-closeup =  ->
	it.prevent-default!
	console.log it.target
	media-element = it.target.closest(\figure).query-selector(\img).clone-node true
	media-element.class-list.add \media
	media-element.set-attribute \src (media-element.dataset.full-src)
	media-element.remove-attribute \srcset
	media-element.remove-attribute \sizes
	popup = document.get-element-by-id \media-closeup
	target = popup.first-child
	target.replace-with media-element
	popup.show-modal!

document.query-selector-all \figure .for-each ->
	content-type = it.dataset.contentType
	general-type = (content-type.split '/')[0]
	if general-type == "image"
		it.add-event-listener \click, open-media-closeup
