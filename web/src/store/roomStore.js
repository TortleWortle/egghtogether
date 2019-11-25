import { createContextStore, action } from 'easy-peasy'
import { RoomSocket } from '../ws'



export const RoomStore = createContextStore(({ id }) => {
	return {
		id,
		player: {
			srcObject: null,
			setSrcObject: action((state, payload) => {
				state.srcObject = payload;
			})
		},
		chat: {
			messages: [],
			sendMessage: action((state, payload) => {
				state.messages.push({
					sender: "System",
					nickname: "Test",
					message: payload
				})
			}),
			addMessage: action((state, payload) => {
				state.messages.push({
					sender: payload.sender,
					nickname: payload.nickname,
					message: payload.message
				})
			})
		},
	}
})