open-media-closeup =  ->
	it.prevent-default!
	media-element = it.target.clone-node true
	media-element.class-list.add \media
	media-element.set-attribute \src (media-element.dataset.full-src)
	popup = document.get-element-by-id \media-closeup
	target = popup.first-child
	target.replace-with media-element
	popup.show-modal!

update-dimensions-pdf-iframe = ->
	row = it.target.parent-element.parent-element
	iframe = it.query-selector \iframe
	iframe.width = row.get-bounding-client-rect().width
	# 21 × 29.4 are A4 paper dimensions
	# TODO: do not hardcode that ratio, get it from the PDF side at compile time.
	iframe.height = (29.4/21) * iframe.width

document.query-selector-all \figure .for-each ->
	content-type = it.dataset.contentType
	general-type = (content-type.split '/')[0]
	if not (general-type == \video or content-type == "application/pdf")
		it.add-event-listener \click, open-media-closeup
	if content-type == "application/pdf"
		update-dimensions-pdf-iframe it
		window.add-event-listener \resize, -> update-dimensions-pdf-iframe it
