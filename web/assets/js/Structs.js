export class PaneLoadedEvent extends Event {
    constructor(_paneTarget, _paneDetail, _paneURL) {
        super("paneLoaded");
        this.paneTarget = _paneTarget;
        this.paneDetail = _paneDetail;
        this.paneURL = _paneURL;
    }
}
export class MessageReceivedEvent extends Event {
    constructor(_messageContent) {
        super("messageReceived");
        this.messageContent = _messageContent;
    }
}
export class NotificationReceivedEvent extends Event {
    constructor(_notificationContent) {
        super("notificationReceived");
        this.notificationContent = _notificationContent;
    }
}
