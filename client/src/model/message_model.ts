export type MessageModel = {
    id: string,
    roomID: string,
    userID: string,
    toUserID: string,
    text: string,
    timeStamp: string,
};

export type InputCreateMessage = {
    text: string,
};

export type InputUpdateMessage = {
    text: string,
}

export type InputWS = {
    roomID: string,
    toUserID: string,
    message: string,
}