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
function SendMessage(Content) {
    var _a, _b, _c;
    return __awaiter(this, void 0, void 0, function* () {
        const TargetUser = window.location.href.split("/messages/")[1];
        const ContentMessages = document.querySelectorAll(".contentMessages")[0];
        const ContentMessagesPreSend = ContentMessages.innerHTML;
        ContentMessages.innerHTML = `<div class="message sending">${Content}</div>` + ContentMessages.innerHTML;
        document.querySelectorAll(".content.messages #SendButton")[0].classList.remove("show");
        document.querySelectorAll(".content.messages #MediaButton")[0].classList.add("show");
        const MessageField = document.querySelectorAll(".content.messages #MessageField")[0];
        MessageField.value = "";
        MessageField.disabled = true;
        (_a = MessageField.parentElement) === null || _a === void 0 ? void 0 : _a.classList.add("disabled");
        // We visually send the message
        let response = yield fetch("/api/messages/sendMessage", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ target: TargetUser, content: Content }),
        });
        let result = "";
        if (response.ok) {
            if (response.status == 200) {
                console.log("📭 Message sent.");
                MessageField.disabled = false;
                (_b = MessageField.parentElement) === null || _b === void 0 ? void 0 : _b.classList.remove("disabled");
                MessageField.focus();
                ContentMessages.innerHTML = `<div class="message sent">${Content}</div>` + ContentMessagesPreSend;
            }
        }
        else {
            console.log(`Error!`);
            MessageField.disabled = false;
            (_c = MessageField.parentElement) === null || _c === void 0 ? void 0 : _c.classList.remove("disabled");
            MessageField.focus();
            ContentMessages.innerHTML = `<div class="message error">${Content}</div>` + ContentMessagesPreSend;
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
