export interface INotification {
    messageType: string;
    data: never;
}

export enum WebsocketMessageType {
    USER_LEVEL_UPDATE = "userLevelUpdate",
    BOARD_POST_USER_VOTED = "boardPostUserVoted"
}

export function parseNotification(notification: MessageEvent) {
    return JSON.parse(notification.data);
}

function isMessageType(notification: MessageEvent, messageType: WebsocketMessageType): boolean {
    return parseNotification(notification).messageType === messageType;
}

export function isUserLevelUpdate(notification: MessageEvent) {
    return isMessageType(notification, WebsocketMessageType.USER_LEVEL_UPDATE);
}

export function isBoardPostUserVoted(notification: MessageEvent) {
    return isMessageType(notification, WebsocketMessageType.BOARD_POST_USER_VOTED);
}