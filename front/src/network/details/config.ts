import {request} from '../request'
import {SysConfig} from "@/types";

export const getConfigs = () => {
  return request({
    url: '/admin/config',
    method: 'get',
  }) as unknown as Promise<SysConfig>
}

export const updateConfig = (config:any) => {
  return request({
    url: '/admin/config/update',
    method: 'post',
    data:config
  }) as unknown as Promise<boolean>
}
