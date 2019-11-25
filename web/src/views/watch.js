import React, { useRef, useEffect, useState } from 'react'
import useToggleVideo from '../hooks/useToggleVideo'
import { RoomStore } from "../store/roomStore"
import { useSocket } from "../hooks/useSocket"

import "webrtc-adapter"

const videos = ["https://cdn.discordapp.com/attachments/593820938446045194/645379605980512256/e37840dae01f4a43fe9c61cc78e5c4d4.mp4", "https://cdn.discordapp.com/attachments/593820938446045194/645414309026856970/2d5e9a72b9493ed7920ea00bf8308fde.mp4"]
const video = videos[Math.floor((Math.random() * 100) % videos.length)]



export default () => {
	const videoEl = useRef(null)
	const chatEl = useRef(null)
	const inputRef = useRef(null)

	const [isPlaying, toggleVideo] = useToggleVideo(videoEl)

	// const messages = RoomStore.useStoreState(state => state.chat.messages)
	const id = RoomStore.useStoreState(state => state.id)
	const srcObject = RoomStore.useStoreState(state => state.player.srcObject)

	useEffect(() => {
		if (srcObject) {
			videoEl.current.srcObject = srcObject
			videoEl.current.play();
		}
	}, [srcObject])

	// const sendMessage = RoomStore.useStoreActions(state => state.chat.sendMessage)
	const { sendMessage } = useSocket();

	// useEffect(() => {
	// 	let id = setInterval(() => {
	// 		sendMessage("Woop")
	// 	}, 5000)
	// 	return () => {
	// 		clearInterval(id)
	// 	}
	// })

	// const [chatMessage, setChatMessage] = useState("");

	// useEffect(() => {
	// 	if (chatEl.current) {
	// 		chatEl.current.scrollTop = chatEl.current.scrollHeight
	// 	}
	// }, [messages])

	// useEffect(() => {
	// 	function ev(e) {
	// 		if (e.key == "Enter") {
	// 			sendMessage(chatMessage)
	// 			setChatMessage("")
	// 		}
	// 	}
	// 	inputRef.current.addEventListener('keydown', ev)
	// 	return () => {
	// 		inputRef.current.removeEventListener('keydown', ev)
	// 	}
	// }, [inputRef.current, chatMessage])

	return (
		<div className="bg-indigo-900 text-white h-screen flex">
			<div className="w-3/4">
				<video controls autoPlay={true} style={{ maxHeight: `calc(100vh)` }} loop className={`w-full`} ref={videoEl} src={video}></video>
				<button className={`bg-indigo-700 text-white p-3 py-1 m-2 rounded`} onClick={() => {
					toggleVideo()
				}}>{isPlaying ? "Pause" : "Play"}</button>
				RoomID: {id}
			</div>
			<div className="flex flex-col w-1/4">
				{/* <ul className="flex-grow overflow-scroll" ref={chatEl}>
					{messages.map((msg, index) => (
						<li key={index} className="my-1">{msg.nickname}: {msg.message}</li>
					))}
				</ul>
				<div className="mb-4 flex">
					<input ref={inputRef} className="text-black flex-grow border-box" type="text" value={chatMessage} onChange={(e) => setChatMessage(e.target.value)} />
					<button className="py-2 px-6" onClick={() => {
						// setNickName(chatMessage)
						setChatMessage("")
					}}>Set nickname</button>
				</div> */}
			</div>
		</div >
	)
}