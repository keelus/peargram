import { MessageReceivedEvent } from "./Structs.js"


const SessionWS : WebSocket = new WebSocket("ws://" + window.location.hostname + "/ws");
console.log("👍 Connected to Peargram WebSocket")


SessionWS.onmessage = e => {
	const MessageData = JSON.parse(e.data)

	console.log(MessageData)
	if (MessageData.type == "MESSAGE") {
		console.log("💬 private message received.")

		const messageReceivedEvent = new MessageReceivedEvent(MessageData.Content)
		document.dispatchEvent(messageReceivedEvent)

		// Check current user panel

		// TODO: Show a notification icon on messages?



	} else if (MessageData.type == "NOTIFICATION") {
		console.log("🔔 notification received.")
	}
}