document.addEventListener("click", (e) => {
	const PopupMenu : HTMLElement = document.getElementById("popupMenu") as HTMLButtonElement;
	const Target : HTMLElement = e.target as HTMLElement
	const clickOutsideMenu : boolean = Target.closest(".popupMenu") === null
	const clickOpenMenu : boolean = Target.closest("#moreButton") !== null

	console.log("CLICK")
	console.log(Target.closest(".popupMenu"))
	console.log(Target.closest("#moreButton"))

	
	if(clickOutsideMenu) {
		if(clickOpenMenu) {
			if(PopupMenu.classList.contains("opened"))
				return PopupMenu.classList.remove("opened")
			PopupMenu.classList.add("opened")
		} else {
			PopupMenu.classList.remove("opened")
		}
	}
})