id = -> document.get-element-by-id it

open-media-closeup =  ->
	it.prevent-default!
	window-scroll = window.scrollY

	media-element = it.target.closest(\figure).query-selector(\img).clone-node true
	media-element.class-list.add \media
	media-element.set-attribute \src (media-element.dataset.full-src)
	media-element.remove-attribute \srcset
	media-element.remove-attribute \sizes


	popup = id \media-closeup
	target = popup.query-selector \.media
	loading = popup.query-selector \.loading
	magnifier = popup.query-selector \.magnifier

	if not popup.open
		popup.show-modal!
	setup-spinner!
	media-element.style.opacity = 0
	target.replace-with media-element

	media-element.add-event-listener "load", ->
		media-element.style.opacity = 1
		magnifier.style.background-image = "url(#{media-element.dataset.full-src})"
		stop-spinner!

		zoom-factor = -> parseFloat popup.dataset.zoom-factor
		magnifier-width = -> parseFloat magnifier.style.width - "px"
		magnifier-height = -> parseFloat magnifier.style.height - "px"

		update-zoom-factor = ->
			if not (1 < it <= 10)
				return
			popup.dataset.zoom-factor = it
			[media-height, media-width] = <[height width]>.map -> media-element.getBoundingClientRect![it]
			magnifier.style.background-size = "#{media-width * zoom-factor!}px #{media-height * zoom-factor!}px"
			magnifier-width = magnifier-height = 50 * zoom-factor!
			magnifier.style.width = magnifier-width + "px"
			magnifier.style.height = magnifier-height + "px"

		update-magnifier = ->
			if not it.keep-transition
				magnifier.style.transition = "none"
			for $1 in <[clientX clientY offsetX offsetY]>
				magnifier.dataset[$1] = it[$1]
			magnifier.style.left = (it.clientX - magnifier-width!/2) + "px"
			magnifier.style.top = (it.clientY - magnifier-height!/2) + "px"
			magnifier.style.background-position = "
				left -#{it.offsetX * zoom-factor! - magnifier-width!/2}px top -#{it.offsetY * zoom-factor! - magnifier-height!/2}px
			"

		update-zoom-factor 2
		magnifier.style.display = \block

		media-element.add-event-listener "mouseout",  -> magnifier.style.display = \none
		media-element.add-event-listener "mouseover", -> magnifier.style.display = \block
		media-element.add-event-listener "mousemove", update-magnifier
		media-element.add-event-listener "wheel", (->
			it.prevent-default!
			magnifier.style.transition = "all 250ms ease"
			update-zoom-factor zoom-factor! + 0.5 * Math.sign it.deltaY
			update-magnifier(
				keep-transition: yes
				clientX: magnifier.dataset.client-x
				clientY: magnifier.dataset.client-y
				offsetX: magnifier.dataset.offset-x
				offsetY: magnifier.dataset.offset-y
			)
			window-scroll = window.scrollY
		), passive: false



document.query-selector-all \figure .for-each ->
	content-type = it.dataset.contentType
	general-type = (content-type.split '/')[0]
	if general-type == "image"
		it.add-event-listener \click, open-media-closeup

setup-spinner = ->
	line = -> ".spinner .#{it}"
	lines = -> it.map line

	loading = id \media-closeup .query-selector \.loading

	document.query-selector \.spinner .query-selector-all \line .for-each ->
		it.style.stroke-dasharray = 1000
		it.style.stroke-dashoffset = 1010

	loading.style.opacity = 1

stop-spinner = ->
	loading = id \media-closeup .query-selector \.loading
	loading.style.opacity = 0

queue-next-track = ->
	current-track-number = it.target.get-attribute \title - /\. .+$/ |> parseFloat
	next-track = document.query-selector "[title^=\"0#{current-track-number + 1}. \"]"
	it.target.add-event-listener \ended, ->
		next-track.play! if next-track

document.query-selector-all "audio[title]" .for-each ->
	it.add-event-listener \play, ->
		window.title = "ewen.works: ▶ #{it.get-attribute \title } from #{document.query-selector \h1 .inner-text}"
		queue-next-track it


document.hide-wip-banner = ->
	id \wip-banner .dataset.hidden = yes
	window.local-storage.set-item \wip-banner-hidden, yes

if not window.local-storage.get-item \wip-banner-hidden
	delete id \wip-banner .dataset.hidden
