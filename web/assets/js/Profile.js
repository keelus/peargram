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
    if (!(paneLoadedEvent.paneTarget === "profile" && paneLoadedEvent.paneDetail !== ""))
        return;
    SetupListeners();
});
function SetupListeners() {
    console.log("Setting up");
    const FollowButton = document.querySelectorAll(".content.profile #FollowButton")[0];
    if (!FollowButton)
        return;
    FollowButton.addEventListener("click", (e) => {
        ToggleFollow();
    });
}
function ToggleFollow() {
    return __awaiter(this, void 0, void 0, function* () {
        const TargetUser = window.location.href.split("/profile/")[1];
        let response = yield fetch(`/api/toggleFollow/${TargetUser}`);
        if (response.ok) {
            if (response.status == 200) {
                console.log(response);
                console.log("👍 User (un)followed.");
                let result = yield response.json();
                document.querySelectorAll(".content.profile #FollowButton")[0].classList.remove("follow");
                document.querySelectorAll(".content.profile #FollowButton")[0].classList.remove("following");
                const nowString = (result.NowFollowing ? "Following" : "Follow");
                document.querySelectorAll(".content.profile #FollowButton")[0].innerHTML = nowString;
                document.querySelectorAll(".content.profile #FollowButton")[0].classList.add(nowString.toLowerCase());
            }
            else {
                console.log(`Error!`);
            }
        }
    });
}
export {};
