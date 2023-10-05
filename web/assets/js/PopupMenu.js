"use strict";
document.addEventListener("click", (e) => {
    const PopupMenu = document.getElementById("popupMenu");
    const Target = e.target;
    const clickOutsideMenu = Target.closest(".popupMenu") === null;
    const clickItemMenu = Target.closest(".pmItem") !== null;
    const clickOpenMenu = Target.closest("#moreButton") !== null;
    if (clickOutsideMenu) {
        if (clickOpenMenu) {
            if (PopupMenu.classList.contains("opened"))
                return CloseMenu(PopupMenu);
            OpenMenu(PopupMenu);
        }
        else {
            CloseMenu(PopupMenu);
        }
    }
    else if (clickItemMenu)
        CloseMenu(PopupMenu);
});
const OpenMenu = (pmenu) => { pmenu.classList.add("opened"); };
const CloseMenu = (pmenu) => { pmenu.classList.remove("opened"); };
