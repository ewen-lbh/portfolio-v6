open-media-closeup =  ->
	it |> console.log
	it.prevent-default!
	media-element = it.target.clone-node true
	media-element.class-list.add \media
	popup = document.get-element-by-id \media-closeup
	target = popup.first-child
	target.replace-with media-element
	popup.show-modal!

document.query-selector-all \figure .for-each ->
	content-type = it.dataset.contentType
	general-type = (content-type / "/")[0]
	if not (general-type == "video" or content-type == "application/pdf")
		it.add-event-listener \click, open-media-closeup
