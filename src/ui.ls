open-media-closeup =  ->
	it.prevent-default!
	media-element = it.target.clone-node true
	media-element.class-list.add \media
	media-element.set-attribute \src (media-element.dataset.full-src)
	popup = document.get-element-by-id \media-closeup
	target = popup.first-child
	target.replace-with media-element
	popup.show-modal!

document.query-selector-all \figure .for-each ->
	content-type = it.dataset.contentType
	general-type = (content-type.split '/')[0]
	if general-type == "image"
		it.add-event-listener \click, open-media-closeup
