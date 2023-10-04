import { MessageReceivedEvent } from "./Structs.js"


const SessionWS : WebSocket = new WebSocket("ws://" + window.location.hostname + "/ws");
console.log("👍 Connected to Peargram WebSocket")


SessionWS.onmessage = e => {
	const MessageData = JSON.parse(e.data)

	console.log(MessageData)
	if (MessageData.type == "MESSAGE") {
		console.log("💬 private message received.")

		// Check current user panel
		
		let currentURL = window.location.pathname;
		let urlParts : string[] = currentURL.split("/")
		urlParts = urlParts.filter((item) => {
			return item !== ""
		})

		console.log(urlParts)
		console.log(urlParts.length)

		if (urlParts[0] == "messages") {
			// UPDATE LEFT BAR TODO
			if(urlParts.length == 2) { // Is in a chat
				const ChatUsername : string = urlParts[1]
				// console.log(ChatUsername)
				// console.log(MessageData.actor)
				if(MessageData.Content.Actor == ChatUsername) {
					// console.log("Update current chat!!!")

					
					const messageReceivedEvent = new MessageReceivedEvent(MessageData.Content)
					document.dispatchEvent(messageReceivedEvent)

					// const ContentMessages : HTMLElement = document.querySelectorAll(".contentMessages")[0] as HTMLElement

					// ContentMessages.innerHTML = `<div class="message incoming">${MessageData.Content.Content}</div>` + ContentMessages.innerHTML
					
				}
			}
			
		}

		// TODO: Show a notification icon on messages?



	} else if (MessageData.type == "NOTIFICATION") {
		console.log("🔔 notification received.")
	}
}