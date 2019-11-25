const configuration = { iceServers: [{ urls: 'stun:stun.l.google.com:19302' }] };
const EventEmitter = require('events');
const ws_api_base = "ws://localhost:8000/api"

export class RoomSocket extends EventEmitter {
	ws = null;
	interval = null;
	reconnection_attemps = 3;
	connected = false;
	roomID = null

	constructor(roomID) {
		super();
		this.roomID = roomID
	}
	onmsg(event) {
		console.log(event)
		const e = JSON.parse(event.data);
		this.emit(e.op, e.data);
		// switch (e.op) {
		// 	case "IDENTIFY":
		// 		id = e.data.id
		// 		setupRTC()
		// 		break;
		// 	case "SEND_ANSWER":
		// 		processAnswer(JSON.parse(e.data.answer))
		// 		break;
		// 	case "ICE_CANDIDATE":
		// 		console.log("Got ICE Candidate", e)
		// 		rtc.addIceCandidate(JSON.parse(e.data.candidate))
		// 		break
		// }
	}
	sendJSON(data) {
		this.ws && this.connected && this.ws.send(JSON.stringify(data));
	}
	close() {
		this.ws && this.ws.close()
		this.ws = null;
		clearInterval(this.interval)
		this.connected = false;
	}
	connect() {
		return new Promise((resolve, reject) => {

			const ws = new WebSocket(`${ws_api_base}/room/${this.roomID}/ws`)
			ws.onerror = (err) => {
				console.error("Socket Err")
				reject(err)
			}
			ws.onclose = (w) => {
				console.log("Connection closing", w)
				// if (this.reconnection_attemps > 0) {
				// 	console.log("Attempting to reconnect")
				this.ws = null;
				// 	this.connect();
				// 	this.reconnection_attemps--;
				// }
				this.connected = false;
			}
			ws.onopen = () => {
				this.reconnection_attemps = 3;
				this.connected = true;
				resolve(true);
			}
			ws.onmessage = this.onmsg.bind(this)
			this.interval = setInterval(() => {
				ws.send(JSON.stringify({
					op: "HEARTBEAT"
				}))
			}, 30000)
			this.ws = ws
		})
	}
}

// let ws;
// let rtc;
// let stream;
// let timeout;
// const video = document.querySelector('video')

// function processAnswer(answer) {
// 	if (!rtc) return
// 	rtc.setRemoteDescription(answer)
// 	rtc.oniceconnectionstatechange = function () {
// 		if (rtc.iceConnectionState == 'disconnected') {
// 			rtc = null
// 		}
// 	}
// 	ws.send(JSON.stringify({
// 		op: "READY"
// 	}))
// }
// async function setupRTC() {
// 	rtc = new RTCPeerConnection(configuration)
// 	rtc.addEventListener('track', function (event) {
// 		console.log("Track added")
// 		video.srcObject = event.streams[0]
// 	})
// 	rtc.addEventListener('icecandidate', function (event) {
// 		ws.send(JSON.stringify({
// 			op: "ICE_CANDIDATE",
// 			candidate: JSON.stringify(event.candidate)
// 		}))
// 	})
// 	const offer = await rtc.createOffer({
// 		offerToReceiveAudio: true,
// 		offerToReceiveVideo: true,
// 	})
// 	rtc.setLocalDescription(offer);
// 	ws.send(JSON.stringify({
// 		op: "CREATE_OFFER",
// 		offer: JSON.stringify(offer)
// 	}))
// }