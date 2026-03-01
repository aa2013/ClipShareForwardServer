import {request} from '../request'
import {ChartData, ConnectionStatusResp, ForcedDisconnectionDto} from "@/types";

export const getConnectionStatus = () => {
  return request({
    url: '/admin/connectionStatus',
    method: 'get',
  }) as unknown as Promise<ConnectionStatusResp>
}

export const getChartsData = () => {
  return request({
    url: '/admin/charts',
    method: 'get',
  }) as unknown as Promise<ChartData>
}

export const forcedDisconnection = (params: ForcedDisconnectionDto) => {
  return request({
    url: '/admin/forcedDisconnection',
    method: 'post',
    data: params
  }) as unknown as Promise<boolean>
}
