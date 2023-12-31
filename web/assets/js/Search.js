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
    var _a;
    SetupListeners();
    (_a = document.getElementById("SearchField")) === null || _a === void 0 ? void 0 : _a.focus();
});
document.addEventListener("paneLoaded", (e) => {
    var _a;
    const paneLoadedEvent = e;
    if (!(paneLoadedEvent.paneTarget === "search"))
        return;
    SetupListeners();
    (_a = document.getElementById("SearchField")) === null || _a === void 0 ? void 0 : _a.focus();
});
function SetupListeners() {
    const SearchField = document.querySelectorAll(".content.search #SearchField")[0];
    if (!SearchField)
        return;
    const SearchButton = document.querySelectorAll(".content.search #SearchButton")[0];
    SearchField.addEventListener("keydown", (e) => {
        const SearchContent = SearchField.value;
        if (e.key == "Enter") {
            PerformSearch(SearchContent);
        }
    });
    SearchButton.addEventListener("click", (e) => {
        const SearchContent = SearchField.value;
        PerformSearch(SearchContent);
    });
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
function PerformSearch(content) {
    return __awaiter(this, void 0, void 0, function* () {
        let response = yield fetch(`/api/searchUsers/${content}`);
        if (response.ok) {
            if (response.status == 200) {
                let result = yield response.json();
                console.log("🔎 Searched.");
                console.log(result);
                document.querySelectorAll(".infoResults")[0].classList.remove("show");
                document.querySelectorAll(".noResults")[0].classList.remove("show");
                const ContentResults = document.querySelectorAll(".results")[0];
                const CurrentContents = document.querySelectorAll(".results .result");
                CurrentContents.forEach((elem, i) => {
                    elem.remove();
                });
                if (result.Amount === 0) {
                    document.querySelectorAll(".noResults")[0].classList.add("show");
                }
                else {
                    for (let i = 0; i < result.Amount; i++) {
                        ContentResults.innerHTML = `<div class="result EmbeddedLoad" target="profile" detail="${result.Coincidences[i].Username}">
					<div class="avatar" style="background-image:url(/api/users/getAvatar/${result.Coincidences[i].Username})"></div>
					<div class="details">
						<div class="username">${result.Coincidences[i].Username}</div>
						<div class="followers">${result.Coincidences[i].FollowerAmount} followers</div>
					</div>
				</div>` + ContentResults.innerHTML;
                    }
                }
                // document.querySelectorAll(".content.messages #SendButton")[0].classList.remove("show")
                // document.querySelectorAll(".content.messages #MediaButton")[0].classList.add("show")
                // const MessageField : HTMLInputElement = document.querySelectorAll(".content.messages #MessageField")[0] as HTMLInputElement
                // MessageField.value = ""
                // MessageField.focus()
            }
            else {
                console.log(`Error!`);
            }
        }
    });
}
export {};
