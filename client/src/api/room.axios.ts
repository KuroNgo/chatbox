import http from "@/constant/api.axios"
import type { APIResponse } from "@/constant/type"
import type { InputCreateRoom, RoomModel } from "@/model/room_model"

async function getManyRoom() {
    return await http.get<APIResponse<RoomModel[]>>("room")
}

async function getOneRoom(id: string) {
    return await http.get<APIResponse<RoomModel>>(`room/${id}`)
}

async function getOneRoomByName(name: string) {
    return await http.get<APIResponse<RoomModel>>(`room/${name}`)
}

async function createRoom(input: InputCreateRoom){
    return await http.post<APIResponse<string>>("room", input)
}

async function deleteRoom(id: string) {
    return await http.delete<APIResponse<boolean>>(`room/${id}`)
}

async function updateRoom(id: string, room: RoomModel) {
    return await http.patch<APIResponse<string>>(`room/${id}`, room)
}

export default {
    getManyRoom,
    getOneRoom,
    getOneRoomByName,
    createRoom,
    updateRoom,
    deleteRoom,
};