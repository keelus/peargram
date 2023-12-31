import { ACTIVE_PANE, ACTIVE_PANE_DETAIL, PANE } from "./EmbedHandler.js";
import { PaneLoadedEvent, MessageReceivedEvent } from "./Structs.js";
import { RenderDateShort } from "./Utils.js";

let LeftBarUpdateInitialized = false;

document.addEventListener("DOMContentLoaded", (e) => {
	SetupListeners()

	if(ACTIVE_PANE == PANE.MESSAGES) {
		LeftBarUpdate(true)
		LeftBarUpdateInitialized = true;
	}
})

document.addEventListener("paneLoaded", (e) => {
	const paneLoadedEvent : PaneLoadedEvent = e as PaneLoadedEvent;

	if (paneLoadedEvent.paneTarget !== "messages") {
		LeftBarUpdateInitialized = false;
		return;
	}
	
	if(!LeftBarUpdateInitialized) {
		LeftBarUpdateInitialized = true;
		LeftBarUpdate(true)
	}

	if (paneLoadedEvent.paneDetail === "") return;

	SetupListeners()
})

function SetupListeners(){
	console.log("Setting up")
	const MessageField : HTMLInputElement = document.querySelectorAll(".content.messages #MessageField")[0] as HTMLInputElement
	if(!MessageField) return
	const MediaButton : HTMLElement = document.querySelectorAll(".content.messages #MediaButton")[0] as HTMLButtonElement
	const SendButton : HTMLElement = document.querySelectorAll(".content.messages #SendButton")[0] as HTMLButtonElement
	
	const LoadMoreButton : HTMLElement = document.getElementById("LoadMore") as HTMLElement

	["keyup", "keydown", "value", "focus"].forEach((listener, i) => {
		MessageField.addEventListener(listener, (e) => {
			const MessageContent : string = MessageField.value

			if(MessageContent.replace(/\s/g, "").length > 0) {
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

	document.addEventListener("click", (e) => {
		const ElementTarget = e.target as HTMLElement
		if(!ElementTarget) return;

		const LoadMoreButton = ElementTarget.closest(".loadMore")
		if(!LoadMoreButton) return;

		
		LoadMoreMessages()
	})
	// LoadMoreButton.addEventListener("click", () => {
	// 	LoadMoreMessages()
	// })
}

async function SendMessage(Content : string) {
	if(Content.replace(/\s/g, "") == "") return;

	const TargetUser : string = GetChatUser()
	
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
			
			LeftBarMessage(TargetUser, Content, false)
			LeftBarUpdate(false)

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

	if(ACTIVE_PANE == PANE.MESSAGES) {
		if(ACTIVE_PANE_DETAIL != "") { // Is in a chat
			const ChatUsername : string = ACTIVE_PANE_DETAIL
			if(messageReceivedEvent.messageContent.Actor == ChatUsername) { // If you are in the oncoming message's chat
				LeftBarMessage(messageReceivedEvent.messageContent.Actor, messageReceivedEvent.messageContent.Content, false)

				const ContentMessages : HTMLElement = document.querySelectorAll(".contentMessages")[0] as HTMLElement
				ContentMessages.innerHTML = `<div class="message incoming">${messageReceivedEvent.messageContent.Content}</div>` + ContentMessages.innerHTML
				
			} else {
				LeftBarMessage(messageReceivedEvent.messageContent.Actor, messageReceivedEvent.messageContent.Content, true)
			}
		} else {
			LeftBarMessage(messageReceivedEvent.messageContent.Actor, messageReceivedEvent.messageContent.Content, true)
		}

		LeftBarUpdate(false);

		
		let beat = new Audio('/assets/media/audio/message.mp3');
		beat.play();
	}
})


function LeftBarMessage(Username : string, MessageContentText : string, Received : boolean) {
	const ChatItem : HTMLElement = document.querySelector(`.chats > .chat[chat-user='${Username}']`) as HTMLElement
	const ChatItemDetails : HTMLElement = ChatItem.querySelector(`.details > .message`) as HTMLElement
	if(!ChatItemDetails) return;

	const ChatText : HTMLElement = ChatItemDetails.querySelector(".text") as HTMLElement
	const ChatDate : HTMLElement = ChatItemDetails.querySelector(".date") as HTMLElement

	ChatText.innerHTML = MessageContentText;
	ChatDate.innerHTML = "•  now"

	ChatItem.setAttribute("last-message-date", (new Date().getTime()/1000).toString())

	if(Received)
		ChatItem.classList.add("hasNotification")
}

function LeftBarUpdate(Repeat : boolean) {
	const ChatItemsParent = document.querySelector(".chats")

	if(!ChatItemsParent) return;

	const ChatItems : HTMLElement[] = Array.from(ChatItemsParent.querySelectorAll(".chat"), Chat => Chat as HTMLElement)

	const SortedChatItems = ChatItems.sort((ChatItemA, ChatItemB) => {
		const DateA : number = parseInt(ChatItemA.getAttribute("last-message-date") || "0");
		const DateB : number = parseInt(ChatItemB.getAttribute("last-message-date") || "0");

		if(DateA > DateB) return -1;
		if(DateA < DateB) return 1;
		return 0;
	})

	ChatItemsParent.innerHTML = "";
	SortedChatItems.forEach((elem) => {
		const UnixDateInt : number = parseInt(elem.getAttribute("last-message-date") || "0")
		
		const VisualDate : string = RenderDateShort(UnixDateInt) !== "0s" ? RenderDateShort(UnixDateInt) : "now"

		const ChatDate : HTMLElement = elem.querySelector(".details > .message > .date") as HTMLElement
		if(!ChatDate) return;

		ChatDate.innerText = "• " + VisualDate
		
		ChatItemsParent.appendChild(elem);
	})


	if(Repeat) setTimeout(() => LeftBarUpdate(true), 1000 * 10); // Update each 10s
}

type Message = {
	ID:        number;
	Actor:     string;
	Target:    string;
	Content:   string;
	Reactions: string;
	Date:      number;
}

async function LoadMoreMessages() {
	const ChatMessagesParent = document.querySelector(".contentMessages")

	if(!ChatMessagesParent) return;

	const ChatMessages : HTMLElement[] = Array.from(ChatMessagesParent.querySelectorAll(".message:not(.sending):not(.error)"), Message => Message as HTMLElement)

	const OffsetMessages = ChatMessages.length

	const TargetUser = GetChatUser()
	
	const LoadMore : HTMLButtonElement = document.getElementById("LoadMore") as HTMLButtonElement
	if(!LoadMore) return;
	LoadMore.disabled = true

	let response : Response = await fetch(`/api/messages/getMessages?username=${TargetUser}&offset=${OffsetMessages}`);
	let result;
	if( response.ok ) {
		result = await response.json()
		if( response.status == 200) {
			LoadMore.disabled = false
			const StartMessage = document.querySelector(".start")
			if(!StartMessage) return;

			if(result.Messages.length == 0) return;
		
			
			
			const pending = result.MessageTotalCount - OffsetMessages - result.Messages.length

			if(pending == 0) {
				StartMessage.classList.add("show")
				LoadMore.classList.remove("show")
			}
			const ContentMessages = document.querySelectorAll(".contentMessages")[0]

			ContentMessages.removeChild(StartMessage)
			ContentMessages.removeChild(LoadMore)
			result.Messages.forEach((message : Message) => {
				const NewMessageElement = document.createElement('div')
				
				NewMessageElement.classList.add("message")
				NewMessageElement.classList.add(message.Actor == TargetUser ? "incoming" : "sent")
				NewMessageElement.innerText = `${message.Content}`;

				ContentMessages.appendChild(NewMessageElement)
			})
			ContentMessages.appendChild(LoadMore)
			ContentMessages.appendChild(StartMessage)
		}
	} else {
		const ContentMessages = document.querySelectorAll(".contentMessages")[0]
		ContentMessages.innerHTML = ContentMessages.innerHTML + "<div class='errorMessage'>Error loading more messages. Try again later.</div>"
		const LoadMore = document.getElementById("LoadMore")
		LoadMore?.remove()
		console.log("Error loading more messages.");
	}
}

function GetChatUser() : string {
	return ACTIVE_PANE_DETAIL
}