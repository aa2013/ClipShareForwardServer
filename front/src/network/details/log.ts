import {request} from '../request'
import {LogResp} from "@/types";

export const getLogs = (begin?: string) => {
  return request({
    url: '/admin/logs',
    method: 'get',
    params: {begin}
  }) as unknown as Promise<LogResp[]>
}
