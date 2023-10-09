import { PaneLoadedEvent, MessageReceivedEvent } from "./Structs.js";

document.addEventListener("DOMContentLoaded", (e) => {
	SetupListeners()
})

document.addEventListener("paneLoaded", (e) => {
	const paneLoadedEvent : PaneLoadedEvent = e as PaneLoadedEvent;

	if (!(paneLoadedEvent.paneTarget === "index" || paneLoadedEvent.paneTarget === "post"))
		return
	SetupListeners()
})

function SetupListeners() {
	console.log("setuping")
	const likeButtons = document.querySelectorAll(".post > .bottom > .actions > .likeButton")
	const saveButtons = document.querySelectorAll(".post > .bottom > .actions > .bookmarkButton")
	likeButtons.forEach((likeButton, i) => {
		likeButton.addEventListener("click", (e) => {
			const postID = likeButton.getAttribute("post-id")
			if(!postID) return;

			const likeAmount = likeButton.parentElement?.parentElement?.children[1] // TODO: Better
			if(!likeAmount) return;

			LikePost(likeButton, likeAmount, postID)
		})
	})
	saveButtons.forEach((saveButton, i) => {
		saveButton.addEventListener("click", (e) => {
			const postID = saveButton.getAttribute("post-id")
			if(!postID) return;

			SavePost(saveButton, postID)
		})
	})
}

async function LikePost(likeButton : Element, likeAmount : Element, postID : string) {
	let response : Response = await fetch(`/api/posts/toggleLike?postID=${postID}`);
	if( response.ok ) {
		if( response.status == 200) {
			let result = await response.json()

			const likeAmountInt = result.likes;
			likeAmount.innerHTML = likeAmountInt + (likeAmountInt == 1 ? " like" : " likes");

			if(likeButton.children[0].classList.contains("ti-heart-filled")){
				likeButton.children[0].classList.remove("ti-heart-filled")
				likeButton.children[0].classList.add("ti-heart")
			}
			else {
				likeButton.children[0].classList.remove("ti-heart")
				likeButton.children[0].classList.add("ti-heart-filled")
			}
		} else {
			console.log(`Error!`);
		}
	}
}

async function SavePost(sameButton : Element, postID : string) {
	let response : Response = await fetch(`/api/posts/toggleBookmark?postID=${postID}`);
	if( response.ok ) {
		if( response.status == 200) {
			if(sameButton.children[0].classList.contains("ti-bookmark-filled")){
				sameButton.children[0].classList.remove("ti-bookmark-filled")
				sameButton.children[0].classList.add("ti-bookmark")
			}
			else {
				sameButton.children[0].classList.remove("ti-bookmark")
				sameButton.children[0].classList.add("ti-bookmark-filled")
			}
		} else {
			console.log(`Error!`);
		}
	}
}
