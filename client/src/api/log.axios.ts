import http from "@/constant/api.axios"
import type { APIResponse } from "@/constant/type"
import type { Logging } from "@/model/logging_model"

async function getLogging() {
    return await http.get<APIResponse<Logging[]>>("logging")
}

export default {
    getLogging,
};