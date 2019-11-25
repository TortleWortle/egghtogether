import { useState, useEffect } from 'react'

function useToggleVideo(videoEl) {
	const [isPlaying, setIsPlaying] = useState(false)

	useEffect(() => {
		function onPlayState(event) {
			if (event.type === "play") {
				setIsPlaying(true)
			} else {
				setIsPlaying(false)
			}
		}
		videoEl.current.addEventListener('play', onPlayState)
		videoEl.current.addEventListener('pause', onPlayState)
	}, [videoEl.current])

	function toggle() {
		if (isPlaying) {
			videoEl.current.pause()
		} else {
			videoEl.current.play()
		}
	}

	return [isPlaying, toggle];
}

export default useToggleVideo