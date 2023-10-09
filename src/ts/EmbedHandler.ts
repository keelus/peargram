import { PaneLoadedEvent } from "./Structs.js"

export let ACTIVE_PANE = -1;
export let ACTIVE_PANE_DETAIL = "";

const MOUSE_CLICK = {
	LEFT: 0,
	MIDDLE: 1,
	RIGHT: 2
}

export const PANE :{[key:string]: number} = {
	UNDEFINED:-1,
	INDEX:0,
	SEARCH:1,
	MESSAGES:2,
	NOTIFICATIONS:3,
	PROFILE:4,
	SETTINGS:5,
	ACTIVITY:6,
	SAVED:7,
}

document.addEventListener("mousedown", (e) => { 
	if(!e.target)
		return

	const targetElement = e.target as HTMLElement;
	if(targetElement.closest(".EmbeddedLoad") && e.button === MOUSE_CLICK.MIDDLE)
		e.preventDefault();
})

document.addEventListener("mouseup", (e) => {
	if(!e.target)
		return

	const targetButton = e.target as HTMLElement | null
	const embeddedLoadButton = targetButton?.closest(".EmbeddedLoad")

	if(embeddedLoadButton) {
		const target : string = embeddedLoadButton.getAttribute("target") || ""
		const detail : string = embeddedLoadButton.getAttribute("detail") || ""

		LoadEmbedded(target, detail, e, true)
	}
})

async function LoadEmbedded(target: string, detail: string, event: MouseEvent, pushToHistory: boolean) {
	let backupTarget = target
	
	if(target === "index")
		target=""

	if (target === "profile")
		target="profile/" + detail

	if (target === "post")
		target="post/" + detail
	
	if (target === "messages")
		target="messages/" + detail

	if (event.button === MOUSE_CLICK.LEFT && ! event.ctrlKey) {
		let response : Response = await fetch(`/${target}?type=short`);
		let result : string = "";
		if( response.ok ) {
			if( response.status == 200) {
				result = await response.text();
			} else {
				console.log(`Error!`);
			}
		}

		document.querySelector("head > title")?.remove()
		document.querySelectorAll("body > .content")[0].remove()
		document.querySelectorAll("body")[0].innerHTML += result

		const loadEvent = new PaneLoadedEvent(backupTarget, detail, target)
		document.dispatchEvent(loadEvent)

		if(pushToHistory)
			window.history.pushState("/" + (target || ""), "/" + (target || ""), "/" + target)

		UpdateActivePane()
	}

	if (event.button === MOUSE_CLICK.LEFT && event.ctrlKey || event.button === MOUSE_CLICK.MIDDLE)
		window.open(`/${target}`, "_blank")
	
}

addEventListener("popstate", function (e) {
	e.preventDefault();

    let previousUrl = location.pathname;
	CheckPanel(previousUrl)
	// window.location.href = previousUrl
});

function CheckPanel(url : string) {
	let urlParts : string[] = url.split("/")

	urlParts = urlParts.filter((item) => {
		return item !== ""
	})

	const clickEvent = new MouseEvent("mousedown")

	if(urlParts.length === 0) 
		LoadEmbedded("index", "", clickEvent, false)
	else if (urlParts.length === 1) {
		LoadEmbedded(urlParts[0], "", clickEvent, false)
	} else if (urlParts.length === 2) {
		LoadEmbedded(urlParts[0], urlParts[1], clickEvent, false)
	}
}

function UpdateActivePane() {
	console.log("Updating pane")
	let paneStr = ""
	ACTIVE_PANE = PANE.UNDEFINED;
	ACTIVE_PANE_DETAIL = "";

	let currentURL = window.location.pathname;
	let urlParts : string[] = currentURL.split("/")
	urlParts = urlParts.filter((item) => {
		return item !== ""
	})

	if(urlParts.length > 0) {
		const paneStr = urlParts[0].toUpperCase();
		const paneVal = PANE[paneStr];

		if(paneVal !== undefined) {
			ACTIVE_PANE = paneVal;
			ACTIVE_PANE_DETAIL = urlParts[1];

			if(ACTIVE_PANE == PANE.MESSAGES && urlParts.length > 1)
				ACTIVE_PANE_DETAIL = urlParts[1]
		}
	} else {
		ACTIVE_PANE = PANE.INDEX;
	}

}
function print(){
	console.log(ACTIVE_PANE)
	setTimeout(print, 500);
}
document.addEventListener("DOMContentLoaded", () => {UpdateActivePane()/*, print()*/})