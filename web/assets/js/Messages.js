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
    if (!(paneLoadedEvent.paneTarget === "messages" && paneLoadedEvent.paneDetail !== ""))
        return;
    SetupListeners();
});
function SetupListeners() {
    console.log("Setting up");
    const MessageField = document.querySelectorAll(".content.messages #MessageField")[0];
    if (!MessageField)
        return;
    const MediaButton = document.querySelectorAll(".content.messages #MediaButton")[0];
    const SendButton = document.querySelectorAll(".content.messages #SendButton")[0];
    ["keyup", "keydown", "value", "focus"].forEach((listener, i) => {
        MessageField.addEventListener(listener, (e) => {
            const MessageContent = MessageField.value;
            console.log(MessageContent);
            if (MessageContent.length > 0) {
                SendButton.classList.add("show");
                MediaButton.classList.remove("show");
            }
            else {
                SendButton.classList.remove("show");
                MediaButton.classList.add("show");
            }
        });
    });
    // SendButton.addEventListener()
    SendButton.addEventListener("click", () => {
        const MessageContent = MessageField.value;
        SendMessage(MessageContent);
    });
    MessageField.addEventListener("keydown", (e) => {
        if (e.key !== "Enter")
            return;
        const MessageContent = MessageField.value;
        SendMessage(MessageContent);
    });
}
function SendMessage(content) {
    return __awaiter(this, void 0, void 0, function* () {
        const TargetUser = window.location.href.split("/messages/")[1];
        let response = yield fetch(`/api/messages/sendMessage/${TargetUser}/${content}`);
        let result = "";
        if (response.ok) {
            if (response.status == 200) {
                console.log("📭 Message sent.");
                const ContentMessages = document.querySelectorAll(".contentMessages")[0];
                ContentMessages.innerHTML = `<div class="message">${content}</div>` + ContentMessages.innerHTML;
                document.querySelectorAll(".content.messages #SendButton")[0].classList.remove("show");
                document.querySelectorAll(".content.messages #MediaButton")[0].classList.add("show");
                const MessageField = document.querySelectorAll(".content.messages #MessageField")[0];
                MessageField.value = "";
                MessageField.focus();
            }
            else {
                console.log(`Error!`);
            }
        }
    });
}
document.addEventListener("messageReceived", (e) => {
    const messageReceivedEvent = e;
    console.log(messageReceivedEvent.messageContent);
    const ContentMessages = document.querySelectorAll(".contentMessages")[0];
    ContentMessages.innerHTML = `<div class="message incoming">${messageReceivedEvent.messageContent.Content}</div>` + ContentMessages.innerHTML;
    let beat = new Audio('/assets/media/audio/message.mp3');
    beat.play();
});
export {};
