export class PaneLoadedEvent extends Event {
	paneTarget : string;
	paneDetail : string;
	paneURL : string;

	constructor(_paneTarget: string, _paneDetail: string, _paneURL: string) {
		super("paneLoaded");
		this.paneTarget = _paneTarget
		this.paneDetail = _paneDetail
		this.paneURL = _paneURL
	}
}
export class MessageReceivedEvent extends Event {
	messageContent : any;

	constructor(_messageContent: any) {
		super("messageReceived");
		this.messageContent = _messageContent
	}
}
export class NotificationReceivedEvent extends Event {
	notificationContent : any;

	constructor(_notificationContent: any) {
		super("notificationReceived");
		this.notificationContent = _notificationContent
	}
}