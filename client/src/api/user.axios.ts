import http from "@/constant/api.axios"
import type { APIResponse } from "@/constant/type";
import type {  Response, InputSignIn, UserModel, InputVerificate, InputSignUp, InputForgetPassword } from "@/model/user_model";

async function login(input:InputSignIn) {
    return await http.post<APIResponse<Response>>("user", input)
}

async function getMe() {
    return await http.get<APIResponse<UserModel>>("user/info")
}

async function logout() {
    return await http.get<APIResponse<string>>("user/logout")
}

async function verify(input:InputVerificate) {
    return await http.patch<APIResponse<string>>("user/verify", input)
}

async function signup(input:InputSignUp) {
    return await http.post<APIResponse<string>>("user", input)
}

async function google() {
    return await http.get<APIResponse<UserModel>>("google/callback")
}

async function forgetpassword(input: InputForgetPassword) {
return await http.post<APIResponse<string>>("user/forget", input)
}

export default {
    login,
    getMe,
    logout,
    verify,
    signup,
    google,
    forgetpassword,
}

