"use strict";
document.addEventListener("click", (e) => {
    const PopupMenu = document.getElementById("popupMenu");
    const Target = e.target;
    const clickOutsideMenu = Target.closest(".popupMenu") === null;
    const clickOpenMenu = Target.closest("#moreButton") !== null;
    if (clickOutsideMenu) {
        if (clickOpenMenu) {
            if (PopupMenu.classList.contains("opened"))
                return PopupMenu.classList.remove("opened");
            PopupMenu.classList.add("opened");
        }
        else {
            PopupMenu.classList.remove("opened");
        }
    }
});
