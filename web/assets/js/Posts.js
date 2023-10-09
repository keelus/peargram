var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
document.addEventListener("DOMContentLoaded", (e) => {
    SetupListeners();
});
document.addEventListener("paneLoaded", (e) => {
    const paneLoadedEvent = e;
    if (!(paneLoadedEvent.paneTarget === "index" || paneLoadedEvent.paneTarget === "post"))
        return;
    SetupListeners();
});
function SetupListeners() {
    console.log("setuping");
    const likeButtons = document.querySelectorAll(".post > .bottom > .actions > .likeButton");
    const saveButtons = document.querySelectorAll(".post > .bottom > .actions > .bookmarkButton");
    likeButtons.forEach((likeButton, i) => {
        likeButton.addEventListener("click", (e) => {
            var _a, _b;
            const postID = likeButton.getAttribute("post-id");
            if (!postID)
                return;
            const likeAmount = (_b = (_a = likeButton.parentElement) === null || _a === void 0 ? void 0 : _a.parentElement) === null || _b === void 0 ? void 0 : _b.children[1]; // TODO: Better
            if (!likeAmount)
                return;
            LikePost(likeButton, likeAmount, postID);
        });
    });
    saveButtons.forEach((saveButton, i) => {
        saveButton.addEventListener("click", (e) => {
            const postID = saveButton.getAttribute("post-id");
            if (!postID)
                return;
            SavePost(saveButton, postID);
        });
    });
}
function LikePost(likeButton, likeAmount, postID) {
    return __awaiter(this, void 0, void 0, function* () {
        let response = yield fetch(`/api/posts/toggleLike?postID=${postID}`);
        if (response.ok) {
            if (response.status == 200) {
                let result = yield response.json();
                const likeAmountInt = result.likes;
                likeAmount.innerHTML = likeAmountInt + (likeAmountInt == 1 ? " like" : " likes");
                if (likeButton.children[0].classList.contains("ti-heart-filled")) {
                    likeButton.children[0].classList.remove("ti-heart-filled");
                    likeButton.children[0].classList.add("ti-heart");
                }
                else {
                    likeButton.children[0].classList.remove("ti-heart");
                    likeButton.children[0].classList.add("ti-heart-filled");
                }
            }
            else {
                console.log(`Error!`);
            }
        }
    });
}
function SavePost(sameButton, postID) {
    return __awaiter(this, void 0, void 0, function* () {
        let response = yield fetch(`/api/posts/toggleBookmark?postID=${postID}`);
        if (response.ok) {
            if (response.status == 200) {
                if (sameButton.children[0].classList.contains("ti-bookmark-filled")) {
                    sameButton.children[0].classList.remove("ti-bookmark-filled");
                    sameButton.children[0].classList.add("ti-bookmark");
                }
                else {
                    sameButton.children[0].classList.remove("ti-bookmark");
                    sameButton.children[0].classList.add("ti-bookmark-filled");
                }
            }
            else {
                console.log(`Error!`);
            }
        }
    });
}
export {};
