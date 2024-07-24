export type UserModel = {
    id: string,
    fullName: string,
    email: string,
    password: string,
    role: string,
    coverURL: string,
    avatarURL: string,
    assetID: string,
    phoneNumber: string,
    provider: string,
    verified: boolean,
    verificationCode: string,
    createdAt: Date,
    updatedAt: Date,
}

export type InputSignIn = {
    email: string,
    password: string,
}

export type InputSignUp = {
    email: string,
    fullName: string,
    password: string,
    avatarURL: string,
    phone: string,
}

export type InputVerificate = {
    verificationCode: string,
}

export type InputForgetPassword = {
    email: string,
}

export type Response = {
    accessToken: string,
    refreshToken: string,
    isLogin: boolean,
}