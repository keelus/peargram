import { PaneLoadedEvent } from "./Structs.js";


document.addEventListener("DOMContentLoaded", (e) => {
	SetupListeners()
})

document.addEventListener("paneLoaded", (e) => {
	const paneLoadedEvent : PaneLoadedEvent = e as PaneLoadedEvent;

	if (!(paneLoadedEvent.paneTarget === "profile" && paneLoadedEvent.paneDetail !== ""))
		return
	SetupListeners()
})

function SetupListeners(){
	console.log("Setting up")
	const FollowButton : HTMLButtonElement = document.querySelectorAll(".content.profile #FollowButton")[0] as HTMLButtonElement
	if(!FollowButton) return

	FollowButton.addEventListener("click", (e) => {
		ToggleFollow()
	})

}

async function ToggleFollow() {
	const TargetUser : string = window.location.href.split("/profile/")[1]

	let response : Response = await fetch(`/api/toggleFollow/${TargetUser}`);
	if( response.ok ) {
		if( response.status == 200) {
			console.log(response)
			console.log("👍 User (un)followed.")
			let result = await response.json()

			document.querySelectorAll(".content.profile #FollowButton")[0].classList.remove("follow")
			document.querySelectorAll(".content.profile #FollowButton")[0].classList.remove("following")

			const nowString : string = (result.NowFollowing ? "Following" : "Follow")
			document.querySelectorAll(".content.profile #FollowButton")[0].innerHTML = nowString
			document.querySelectorAll(".content.profile #FollowButton")[0].classList.add(nowString.toLowerCase())
		} else {
			console.log(`Error!`);
		}
	}
}