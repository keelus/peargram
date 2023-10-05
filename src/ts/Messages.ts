import { PaneLoadedEvent, MessageReceivedEvent } from "./Structs.js";




document.addEventListener("DOMContentLoaded", (e) => {
	SetupListeners()
})

document.addEventListener("paneLoaded", (e) => {
	const paneLoadedEvent : PaneLoadedEvent = e as PaneLoadedEvent;

	if (!(paneLoadedEvent.paneTarget === "messages" && paneLoadedEvent.paneDetail !== ""))
		return
	SetupListeners()
})

function SetupListeners(){
	console.log("Setting up")
	const MessageField : HTMLInputElement = document.querySelectorAll(".content.messages #MessageField")[0] as HTMLInputElement
	if(!MessageField) return
	const MediaButton : HTMLElement = document.querySelectorAll(".content.messages #MediaButton")[0] as HTMLButtonElement
	const SendButton : HTMLElement = document.querySelectorAll(".content.messages #SendButton")[0] as HTMLButtonElement

	["keyup", "keydown", "value", "focus"].forEach((listener, i) => {
		MessageField.addEventListener(listener, (e) => {
			const MessageContent : string = MessageField.value
			console.log(MessageContent)
			if(MessageContent.length > 0) {
				SendButton.classList.add("show")
				MediaButton.classList.remove("show")
			} else {
				SendButton.classList.remove("show")
				MediaButton.classList.add("show")
			}
		})
	});

	// SendButton.addEventListener()
	
	SendButton.addEventListener("click", () => {
		const MessageContent : string = MessageField.value
		SendMessage(MessageContent);
	})
	MessageField.addEventListener("keydown", (e) => {
		if(e.key !== "Enter")
			return
		const MessageContent : string = MessageField.value
		SendMessage(MessageContent);
	})

}

async function SendMessage(Content : string) {
	
	const TargetUser : string = window.location.href.split("/messages/")[1]
	
	const ContentMessages = document.querySelectorAll(".contentMessages")[0]
	const ContentMessagesPreSend = ContentMessages.innerHTML
	ContentMessages.innerHTML = `<div class="message sending">${Content}</div>` + ContentMessages.innerHTML
	document.querySelectorAll(".content.messages #SendButton")[0].classList.remove("show")
	document.querySelectorAll(".content.messages #MediaButton")[0].classList.add("show")
	const MessageField : HTMLInputElement = document.querySelectorAll(".content.messages #MessageField")[0] as HTMLInputElement
	MessageField.value = ""
	MessageField.disabled = true
	MessageField.parentElement?.classList.add("disabled")
	// We visually send the message

	let response : Response = await fetch("/api/messages/sendMessage", {
		method:'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({ target: TargetUser, content: Content }),
	});
	let result : string = "";
	if( response.ok ) {
		if( response.status == 200) {
			console.log("📭 Message sent.")
			MessageField.disabled = false
			MessageField.parentElement?.classList.remove("disabled")
			MessageField.focus()
			ContentMessages.innerHTML = `<div class="message sent">${Content}</div>` + ContentMessagesPreSend

		}
	} else {
		console.log(`Error!`);
		MessageField.disabled = false
		MessageField.parentElement?.classList.remove("disabled")
		MessageField.focus()
		ContentMessages.innerHTML = `<div class="message error">${Content}</div>` + ContentMessagesPreSend
	}
}


document.addEventListener("messageReceived", (e) => {
	const messageReceivedEvent : MessageReceivedEvent = e as MessageReceivedEvent;

	console.log(messageReceivedEvent.messageContent)
	console.log("✅📨👍⚠️")


		
	let currentURL = window.location.pathname;
	let urlParts : string[] = currentURL.split("/")
	urlParts = urlParts.filter((item) => {
		return item !== ""
	})

	console.log(urlParts)
	if (urlParts[0] == "messages") {
		// UPDATE LEFT BAR TODO
		const ChatItem : HTMLElement = document.querySelector(`.chats > .chat[chat-user='${messageReceivedEvent.messageContent.Actor}']`) as HTMLElement
		const ChatItemDetails : HTMLElement = ChatItem.querySelector(`.details > .message`) as HTMLElement
		if(!ChatItemDetails) return;

		const ChatText : HTMLElement = ChatItemDetails.querySelector(".text") as HTMLElement
		const ChatDate : HTMLElement = ChatItemDetails.querySelector(".date") as HTMLElement

		ChatText.innerHTML = messageReceivedEvent.messageContent.Content;
		ChatDate.innerHTML = "•  now"
		

		console.log("LEFT PART!")
		if(urlParts.length == 2) { // Is in a chat
			const ChatUsername : string = urlParts[1]
			// console.log(ChatUsername)
			// console.log(MessageData.actor)
			if(messageReceivedEvent.messageContent.Actor == ChatUsername) { // If you are in the oncoming message's chat
				console.log("Update current chat!!!")
				const ContentMessages : HTMLElement = document.querySelectorAll(".contentMessages")[0] as HTMLElement
				ContentMessages.innerHTML = `<div class="message incoming">${messageReceivedEvent.messageContent.Content}</div>` + ContentMessages.innerHTML

				

				// const ContentMessages : HTMLElement = document.querySelectorAll(".contentMessages")[0] as HTMLElement

				// ContentMessages.innerHTML = `<div class="message incoming">${MessageData.Content.Content}</div>` + ContentMessages.innerHTML
				
			} else {
				ChatItem.classList.add("hasNotification")
			}
		}
		
	}


	let beat = new Audio('/assets/media/audio/message.mp3');
	beat.play();
})