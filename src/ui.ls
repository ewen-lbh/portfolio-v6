open-media-closeup =  ->
	it.prevent-default!
	console.log it.target
	media-element = it.target.closest(\figure).query-selector(\img).clone-node true
	media-element.class-list.add \media
	media-element.set-attribute \src (media-element.dataset.full-src)
	media-element.remove-attribute \srcset
	media-element.remove-attribute \sizes
	popup = document.get-element-by-id \media-closeup
	target = popup.query-selector \.media
	loading = popup.query-selector \.loading
	popup.show-modal!
	setup-spinner!
	# play-spinner!
	media-element.style.opacity = 0
	target.replace-with media-element
	media-element.onload = ->
		console.log \loaded-image
		media-element.style.opacity = 1
		stop-spinner!

document.query-selector-all \figure .for-each ->
	content-type = it.dataset.contentType
	general-type = (content-type.split '/')[0]
	if general-type == "image"
		it.add-event-listener \click, open-media-closeup

setup-spinner = ->
	line = -> ".spinner .#{it}"
	lines = -> it.map line

	loading = document.get-element-by-id \media-closeup .query-selector \.loading

	document.query-selector \.spinner .query-selector-all \line .for-each ->
		it.style.stroke-dasharray = 1000
		it.style.stroke-dashoffset = 1010

	loading.style.opacity = 1

	# T = 600

	# document.spinner-timeline = anime.timeline do
	# 	loop:true
	# 	duration: T
	# 	# direction: \alternate
	# 	easing: \easeInOutExpo

	# keyframe = (delay, obj) -> document.spinner-timeline.add obj, delay * T

	# keyframe 0 do
	# 	targets: lines <[E1-middle E2-middle W-middle]>
	# 	stroke-dashoffset: 0
	# keyframe 1 do
	# 	targets: lines <[E1-top E2-bottom W-right N-right E1-left E2-left]>
	# 	stroke-dashoffset: 0
	# keyframe 2 do
	# 	targets: \.spinner
	# 	rotate: \1turn
	# 	duration: 2 * T
	# keyframe 2.5 do
	# 	targets: lines <[E1-top E2-bottom W-right N-right E1-left E2-left]>
	# 	stroke-dashoffset: 1010
	# keyframe 3 do
	# 	targets: lines <[E1-middle E2-middle W-middle]>
	# 	stroke-dashoffset: 1010
	# # keyframe 4.5 do
	# # 	targets: lines <[N-left E1-bottom E1-right W-bottom]>
	# # 	stroke-dashoffset: 1010

	# document.spinner-animation = anime do
	# 	targets: lines <[N-left E1-bottom E1-right W-bottom]>
	# 	stroke-dashoffset: 0
	# 	duration: T
	# 	easing: \easeInOutSine
	# 	complete: ->
	# 		console.log \completed-initial-anim
	# 		document.spinner-timeline.play!

stop-spinner = ->
	loading = document.get-element-by-id \media-closeup .query-selector \.loading
	loading.style.opacity = 0
	# document.spinner-timeline.pause!
	# document.spinner-timeline.seek 0
	# document.spinner-animation.pause!
	# document.spinner-animation.seek 0

play-spinner = ->
	# document.spinner-timeline.pause!
	# document.spinner-timeline.seek 0
	# document.spinner-animation.restart!

queue-next-track = ->
	current-track-number = it.target.get-attribute \title - /\. .+$/ |> parseFloat
	next-track = document.query-selector "[title^=\"0#{current-track-number + 1}. \"]"
	it.target.add-event-listener \ended, ->
		next-track.play! if next-track

document.query-selector-all "audio[title]" .for-each ->
	it.add-event-listener \play, ->
		window.title = "ewen.works: ▶ #{it.get-attribute \title } from #{document.query-selector \h1 .inner-text}"
		queue-next-track it

	
