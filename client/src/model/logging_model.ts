export type Logging = {
    id: string,
    name: string,
    userID: string,
    roomID: string,
    method:string,
    statusCode: number,
    bodySize: number,
    path: string,
    latency: string,
    error: string,
    activityTime: Date,
    expireAt: Date,
}