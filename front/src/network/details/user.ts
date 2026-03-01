import {request} from '../request'
import {LoginResp} from "@/types";

export const login = (username: string, password: string) => {
  return request({
    url: '/login',
    method: 'POST',
    data: {username, password},
  }) as unknown as Promise<LoginResp>
}
export const logout = () => {
  return request({
    url: '/admin/logout',
    method: 'POST',
  }) as unknown as Promise<boolean>
}
