package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func watchRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["id"]

	w.Header().Set("Content-Type", "text/html")
	watchTemplate.Execute(w, "wss://"+r.Host+"/rooms/"+roomID)
}

var watchTemplate = template.Must(template.New("").Parse(`<video autoplay="autoplay" controls width="1280px" height="720px"></video>
<script src="https://webrtc.github.io/adapter/adapter-latest.js"></script>
<script>
	const configuration = { iceServers: [{ urls: 'stun:stun.l.google.com:19302' }] };
	let ws;
	let rtc;
	let stream;
	let timeout;
	const video = document.querySelector('video')


	function onmsg(evt) {
		const e = JSON.parse(evt.data);
		console.log(e)

		switch (e.op) {
			case "IDENTIFY":
				id = e.data.id
				setupRTC()
				break;
			case "SEND_ANSWER":
				processAnswer(JSON.parse(e.data.answer))
				break;
			case "ICE_CANDIDATE":
				console.log("Got ICE Candidate", e)
				rtc.addIceCandidate(JSON.parse(e.data.candidate))
				break
		}

	}

	function processAnswer(answer) {
		if (!rtc) return

		rtc.setRemoteDescription(answer)

		rtc.oniceconnectionstatechange = function () {
			if (rtc.iceConnectionState == 'disconnected') {
				rtc = null
			}
		}

		ws.send(JSON.stringify({
			op: "READY"
		}))
	}

	async function setupRTC() {
		rtc = new RTCPeerConnection(configuration)

		rtc.addEventListener('track', function (event) {
			console.log("Track added")
			video.srcObject = event.streams[0]
		})

		rtc.addEventListener('icecandidate', function (event) {
			ws.send(JSON.stringify({
				op: "ICE_CANDIDATE",
				candidate: JSON.stringify(event.candidate)
			}))
		})

		const offer = await rtc.createOffer({
			offerToReceiveAudio: true,
			offerToReceiveVideo: true,
		})
		rtc.setLocalDescription(offer);

		ws.send(JSON.stringify({
			op: "CREATE_OFFER",
			offer: JSON.stringify(offer)
		}))
	}

	const close = () => {
		if (!ws) return
		ws.close()
		ws = null
		clearTimeout(timeout)
	}
	const open = (url) => {
		if (ws) return
		ws = new WebSocket(url)

		ws.onopen = function () {
			console.log("Connection opened")
		}
		ws.onerror = function () {
			console.error("Connection opened")
		}
		ws.onclose = function () {
			console.log("Connection closing")
		}
		ws.onmessage = onmsg

		timeout = setInterval(() => {
			ws.send(JSON.stringify({
				op: "HEARTBEAT"
			}))
		}, 30000)
	}
	open("{{.}}")
</script>`))
