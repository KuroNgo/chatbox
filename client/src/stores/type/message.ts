import type { MessageModel } from "@/model/message_model";

export type MessageState = {
    message: MessageModel | null,
    messages: MessageModel[],
}