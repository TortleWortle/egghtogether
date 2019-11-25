import { useRef, useEffect } from 'react'
import { RoomStore } from "../store/roomStore"
import { RoomSocket } from '../ws'

export function useSocket() {
	const id = RoomStore.useStoreState(state => state.id)
	const sendMessage = RoomStore.useStoreActions(state => state.chat.sendMessage)
	const addMessage = RoomStore.useStoreActions(state => state.chat.addMessage)
	const setSrcObject = RoomStore.useStoreActions(state => state.player.setSrcObject)
	const wsRef = useRef(null)
	const pcRef = useRef(null)

	useEffect(() => {
		const ws = new RoomSocket(id);
		const pc = new RTCPeerConnection({ iceServers: [{ urls: 'stun:stun.l.google.com:19302' }] })
		wsRef.current = ws;
		pcRef.current = pc;

		pc.addEventListener('track', function (event) {
			console.log("Track added")
			setSrcObject(event.streams[0])
		})

		ws.on("SEND_ANSWER", (e) => {
			pc.setRemoteDescription(JSON.parse(e.answer))
		})

		ws.on("ICE_CANDIDATE", (e) => {
			pc.addIceCandidate(JSON.parse(e.candidate))
		})


		pc.oniceconnectionstatechange = function () {
			if (pc.iceConnectionState == 'disconnected') {
				console.log("Peer Connection Disconnected")
			}
		}
		ws.connect().then(async () => {
			const offer = await pc.createOffer({
				offerToReceiveAudio: true,
				offerToReceiveVideo: true,
			})
			pc.setLocalDescription(offer)

			ws.sendJSON({
				op: "CREATE_OFFER",
				offer: JSON.stringify(offer)
			})
		})

		return () => {
			wsRef.current.close()
			wsRef.current = null;
			pcRef.current.close();
			pcRef.current = null;
		}
	}, [id])


	return {
		sendMessage: (message) => {
			wsRef.current.sendJSON({
				op: "SEND_MESSAGE",
				data: {
					message
				}
			})
		}
	}
}