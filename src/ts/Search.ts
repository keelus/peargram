import { PaneLoadedEvent, MessageReceivedEvent } from "./Structs.js";

document.addEventListener("DOMContentLoaded", (e) => {
	SetupListeners()
})

document.addEventListener("paneLoaded", (e) => {
	const paneLoadedEvent : PaneLoadedEvent = e as PaneLoadedEvent;

	if (!(paneLoadedEvent.paneTarget === "search"))
		return
	SetupListeners()
})

function SetupListeners() {
	const SearchField : HTMLInputElement = document.querySelectorAll(".content.search #SearchField")[0] as HTMLInputElement
	if(!SearchField) return
	const SearchButton : HTMLButtonElement = document.querySelectorAll(".content.search #SearchButton")[0] as HTMLButtonElement

	SearchField.addEventListener("keydown", (e) => {
		const SearchContent : string = SearchField.value
		if(e.key == "Enter") {
			PerformSearch(SearchContent);
		}
	})
	
	SearchButton.addEventListener("click", (e) => {
		const SearchContent : string = SearchField.value
		PerformSearch(SearchContent);
	
	})

}

// Listeners3.forEach((listener, i) => {
// 	document.addEventListener(listener ,(event) => {
// 		const CurrentTarget : HTMLElement = event.target as HTMLElement

// 		const SearchField : HTMLInputElement = document.querySelectorAll(".content.search #SearchField")[0] as HTMLInputElement
// 		if(!SearchField) return
// 		const SearchContent : string = SearchField.value

// 		if (CurrentTarget == SearchField) {	// If focused on the input and enter key is pressed
// 			const keyEv = event as KeyboardEvent
// 			if(keyEv.key === "Enter") {
// 				PerformSearch(SearchContent);
// 				return
// 			}
// 		}

// 		// If we get here, it means is not a enter key, then check key:
// 		let elemTarget = event.target as HTMLElement
// 		let targParent = elemTarget.parentElement?.closest(".searchButton")

// 		const SearchButton : HTMLButtonElement = document.querySelectorAll(".content.search #SearchButton")[0] as HTMLButtonElement
// 		if (SearchButton !== event.target && SearchButton !== targParent) return
		
// 		PerformSearch(SearchContent);
// 	})
// })
// TODO: Fix while Embedded
async function PerformSearch(content : string) {
	let response : Response = await fetch(`/api/searchUsers/${content}`);
	if( response.ok ) {
		if( response.status == 200) {
			let result = await response.json()
			console.log("🔎 Searched.")
			console.log(result)

			
			document.querySelectorAll(".infoResults")[0].classList.remove("show")
			document.querySelectorAll(".noResults")[0].classList.remove("show")

			const ContentResults = document.querySelectorAll(".results")[0]
			const CurrentContents = document.querySelectorAll(".results .result")
			CurrentContents.forEach((elem, i) => {
				elem.remove()
			})

			if (result.Amount === 0) {
				document.querySelectorAll(".noResults")[0].classList.add("show")
			} else {
				for(let i = 0; i < result.Amount; i++) {
					ContentResults.innerHTML = `<div class="result EmbeddedLoad" target="profile" detail="${result.Coincidences[i].Username}">
					<div class="avatar" style="background-image:url(/api/users/getAvatar/${result.Coincidences[i].Username})"></div>
					<div class="details">
						<div class="username">${result.Coincidences[i].Username}</div>
						<div class="followers">${result.Coincidences[i].FollowerAmount} followers</div>
					</div>
				</div>` + ContentResults.innerHTML
				}
			}

			
			// document.querySelectorAll(".content.messages #SendButton")[0].classList.remove("show")
			// document.querySelectorAll(".content.messages #MediaButton")[0].classList.add("show")
			// const MessageField : HTMLInputElement = document.querySelectorAll(".content.messages #MessageField")[0] as HTMLInputElement
			// MessageField.value = ""
			// MessageField.focus()
		} else {
			console.log(`Error!`);
		}
	}
}