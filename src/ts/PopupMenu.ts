document.addEventListener("click", (e) => {
	const PopupMenu : HTMLElement = document.getElementById("popupMenu") as HTMLButtonElement;
	const Target : HTMLElement = e.target as HTMLElement
	const clickOutsideMenu : boolean = Target.closest(".popupMenu") === null
	const clickItemMenu : boolean = Target.closest(".pmItem") !== null
	const clickOpenMenu : boolean = Target.closest("#moreButton") !== null

	if(clickOutsideMenu) {
		if(clickOpenMenu) {
			if(PopupMenu.classList.contains("opened"))
				return CloseMenu(PopupMenu)
			OpenMenu(PopupMenu)
		} else {
			CloseMenu(PopupMenu)
		}
	} else if(clickItemMenu)
		CloseMenu(PopupMenu)
})

const OpenMenu = (pmenu:HTMLElement) => {pmenu.classList.add("opened")}
const CloseMenu = (pmenu:HTMLElement) => {pmenu.classList.remove("opened")}