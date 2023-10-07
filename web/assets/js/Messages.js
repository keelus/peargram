var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
import { RenderDateShort } from "./Utils.js";
let LeftBarUpdateInitialized = false;
document.addEventListener("DOMContentLoaded", (e) => {
    SetupListeners();
    LeftBarUpdate(true);
    LeftBarUpdateInitialized = true;
});
document.addEventListener("paneLoaded", (e) => {
    const paneLoadedEvent = e;
    if (paneLoadedEvent.paneTarget !== "messages") {
        LeftBarUpdateInitialized = false;
        return;
    }
    if (!LeftBarUpdateInitialized) {
        LeftBarUpdateInitialized = true;
        LeftBarUpdate(true);
    }
    if (paneLoadedEvent.paneDetail === "")
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
    const LoadMoreButton = document.getElementById("LoadMore");
    ["keyup", "keydown", "value", "focus"].forEach((listener, i) => {
        MessageField.addEventListener(listener, (e) => {
            const MessageContent = MessageField.value;
            if (MessageContent.replace(/\s/g, "").length > 0) {
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
    document.addEventListener("click", (e) => {
        const ElementTarget = e.target;
        if (!ElementTarget)
            return;
        const LoadMoreButton = ElementTarget.closest(".loadMore");
        if (!LoadMoreButton)
            return;
        LoadMoreMessages();
    });
    // LoadMoreButton.addEventListener("click", () => {
    // 	LoadMoreMessages()
    // })
}
function SendMessage(Content) {
    var _a, _b, _c;
    return __awaiter(this, void 0, void 0, function* () {
        if (Content.replace(/\s/g, "") == "")
            return;
        const TargetUser = GetChatUser();
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
                LeftBarMessage(TargetUser, Content, false);
                LeftBarUpdate(false);
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
    let currentURL = window.location.pathname;
    let urlParts = currentURL.split("/");
    urlParts = urlParts.filter((item) => {
        return item !== "";
    });
    console.log(urlParts);
    if (urlParts[0] == "messages") {
        if (urlParts.length == 2) { // Is in a chat
            const ChatUsername = urlParts[1];
            if (messageReceivedEvent.messageContent.Actor == ChatUsername) { // If you are in the oncoming message's chat
                LeftBarMessage(messageReceivedEvent.messageContent.Actor, messageReceivedEvent.messageContent.Content, false);
                const ContentMessages = document.querySelectorAll(".contentMessages")[0];
                ContentMessages.innerHTML = `<div class="message incoming">${messageReceivedEvent.messageContent.Content}</div>` + ContentMessages.innerHTML;
            }
            else {
                LeftBarMessage(messageReceivedEvent.messageContent.Actor, messageReceivedEvent.messageContent.Content, true);
            }
        }
        else {
            LeftBarMessage(messageReceivedEvent.messageContent.Actor, messageReceivedEvent.messageContent.Content, true);
        }
        LeftBarUpdate(false);
    }
    let beat = new Audio('/assets/media/audio/message.mp3');
    beat.play();
});
function LeftBarMessage(Username, MessageContentText, Received) {
    const ChatItem = document.querySelector(`.chats > .chat[chat-user='${Username}']`);
    const ChatItemDetails = ChatItem.querySelector(`.details > .message`);
    if (!ChatItemDetails)
        return;
    const ChatText = ChatItemDetails.querySelector(".text");
    const ChatDate = ChatItemDetails.querySelector(".date");
    ChatText.innerHTML = MessageContentText;
    ChatDate.innerHTML = "•  now";
    ChatItem.setAttribute("last-message-date", (new Date().getTime() / 1000).toString());
    if (Received)
        ChatItem.classList.add("hasNotification");
}
function LeftBarUpdate(Repeat) {
    const ChatItemsParent = document.querySelector(".chats");
    if (!ChatItemsParent)
        return;
    const ChatItems = Array.from(ChatItemsParent.querySelectorAll(".chat"), Chat => Chat);
    const SortedChatItems = ChatItems.sort((ChatItemA, ChatItemB) => {
        const DateA = parseInt(ChatItemA.getAttribute("last-message-date") || "0");
        const DateB = parseInt(ChatItemB.getAttribute("last-message-date") || "0");
        if (DateA > DateB)
            return -1;
        if (DateA < DateB)
            return 1;
        return 0;
    });
    ChatItemsParent.innerHTML = "";
    SortedChatItems.forEach((elem) => {
        const UnixDateInt = parseInt(elem.getAttribute("last-message-date") || "0");
        const VisualDate = RenderDateShort(UnixDateInt);
        const ChatDate = elem.querySelector(".details > .message > .date");
        if (!ChatDate)
            return;
        ChatDate.innerText = "• " + VisualDate;
        ChatItemsParent.appendChild(elem);
    });
    if (Repeat)
        setTimeout(() => LeftBarUpdate(true), 1000 * 10); // Update each 10s
}
function LoadMoreMessages() {
    return __awaiter(this, void 0, void 0, function* () {
        const ChatMessagesParent = document.querySelector(".contentMessages");
        if (!ChatMessagesParent)
            return;
        const ChatMessages = Array.from(ChatMessagesParent.querySelectorAll(".message:not(.sending):not(.error)"), Message => Message);
        const OffsetMessages = ChatMessages.length;
        const TargetUser = GetChatUser();
        let response = yield fetch(`/api/messages/getMessages?username=${TargetUser}&offset=${OffsetMessages}`);
        let result;
        if (response.ok) {
            result = yield response.json();
            if (response.status == 200) {
                const StartMessage = document.querySelector(".start");
                const LoadMore = document.getElementById("LoadMore");
                if (!LoadMore || !StartMessage)
                    return;
                if (result.Messages.length == 0)
                    return;
                const pending = result.MessageTotalCount - OffsetMessages - result.Messages.length;
                if (pending == 0) {
                    StartMessage.classList.add("show");
                    LoadMore.classList.remove("show");
                }
                const ContentMessages = document.querySelectorAll(".contentMessages")[0];
                ContentMessages.removeChild(StartMessage);
                ContentMessages.removeChild(LoadMore);
                result.Messages.forEach((message) => {
                    const NewMessageElement = document.createElement('div');
                    NewMessageElement.classList.add("message");
                    NewMessageElement.classList.add(message.Actor == TargetUser ? "incoming" : "sent");
                    NewMessageElement.innerText = `${message.Content}`;
                    ContentMessages.appendChild(NewMessageElement);
                });
                ContentMessages.appendChild(LoadMore);
                ContentMessages.appendChild(StartMessage);
            }
        }
        else {
            const ContentMessages = document.querySelectorAll(".contentMessages")[0];
            ContentMessages.innerHTML = ContentMessages.innerHTML + "<div class='errorMessage'>Error loading more messages. Try again later.</div>";
            const LoadMore = document.getElementById("LoadMore");
            LoadMore === null || LoadMore === void 0 ? void 0 : LoadMore.remove();
            console.log("Error loading more messages.");
        }
    });
}
function GetChatUser() {
    return window.location.href.split("/messages/")[1];
}
