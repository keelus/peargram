// const $2 = (selector:string) => document.querySelectorAll(selector)

// const Listeners : string[] = ["keyup", "keydown", "value", "focus"];
// const Listeners2 : string[] = ["click", "keydown"];

// Listeners.forEach((listener, i) => {
// 	document.addEventListener(listener, (event) => {
// 		const CurrentTarget : HTMLElement = event.target as HTMLElement

// 		const MessageField : HTMLInputElement = $2(".content.messages #MessageField")[0] as HTMLInputElement
// 		if(MessageField !== event.target) return
		
// 		const MessageContent : string = MessageField.value
// 		const MediaButton : HTMLElement = $2(".content.messages #MediaButton")[0] as HTMLButtonElement
// 		const SendButton : HTMLElement = $2(".content.messages #SendButton")[0] as HTMLButtonElement
	
// 		if(MessageContent.length > 0) {
// 			SendButton.classList.add("show")
// 			MediaButton.classList.remove("show")
// 		} else {
// 			SendButton.classList.remove("show")
// 			MediaButton.classList.add("show")
// 		}
// 	})
// })

// Listeners2.forEach((listener, i) => {
// 	document.addEventListener(listener ,(event) => {
// 		const CurrentTarget : HTMLElement = event.target as HTMLElement

// 		const MessageField : HTMLInputElement = $2(".content.messages #MessageField")[0] as HTMLInputElement
// 		if(!MessageField) return
// 		const MessageContent : string = MessageField.value

// 		if (CurrentTarget == MessageField) {	// If focused on the input and enter key is pressed
// 			const keyEv = event as KeyboardEvent
// 			if(keyEv.key === "Enter") {
// 				SendMessage(MessageContent);
// 				return
// 			}
// 		}

// 		const SendButton : HTMLButtonElement = $2(".content.messages #SendButton")[0] as HTMLButtonElement
// 		if (SendButton !== event.target) return
		
// 		SendMessage(MessageContent);
// 	})
// })

// // SendButton.addEventListener("click",(event) => {
// // 	const MessageContent : string = MessageField.value

// // 	SendMessage(MessageContent, TargetUser);
// // })

// async function SendMessage(content : string) {
	
// 	const TargetUser : string = window.location.href.split("/messages/")[1]

// 	let response : Response = await fetch(`/api/messages/sendMessage/${TargetUser}/${content}`);
// 	let result : string = "";
// 	if( response.ok ) {
// 		if( response.status == 200) {
// 			console.log("📭 Message sent.")

// 			const ContentMessages = $2(".contentMessages")[0]

// 			ContentMessages.innerHTML = `<div class="message">${content}</div>` + ContentMessages.innerHTML
			
// 			$2(".content.messages #SendButton")[0].classList.remove("show")
// 			$2(".content.messages #MediaButton")[0].classList.add("show")
// 			const MessageField : HTMLInputElement = $2(".content.messages #MessageField")[0] as HTMLInputElement
// 			MessageField.value = ""
// 			MessageField.focus()
// 		} else {
// 			console.log(`Error!`);
// 		}
// 	}
// }