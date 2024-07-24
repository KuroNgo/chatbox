import http from "@/constant/api.axios";
import type { APIResponse } from "@/constant/type";
import type { MessageModel, InputCreateMessage, InputUpdateMessage, InputWS } from "@/model/message_model";

async function getMessages() {
    return await http.get<APIResponse<MessageModel[]>>("messages");
}

async function deleteMessage(id: string) {
    return await http.delete<APIResponse<boolean>>(`messages/${id}`);
}

async function createMessage(input: InputCreateMessage) {
    return await http.post<APIResponse<MessageModel>>("messages", input);
}

async function updateMessage(id: string, input: InputUpdateMessage) {
    return await http.put<APIResponse<boolean>>(`messages/${id}`, input);
}

async function websocket(input: InputWS) {
    return await http.get<APIResponse<MessageModel>>("/ws")
}

export default {
    getMessages,
    deleteMessage,
    createMessage,
    updateMessage,
};
