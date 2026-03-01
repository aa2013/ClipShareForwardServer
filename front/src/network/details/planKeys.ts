import {request} from '../request'
import {PageData, PageParams, PlanKey, PlanType, TableSortOptions} from "@/types";

const controller = '/admin/planKeys'
export const list = (params: PageParams & { planId?: string, sorts: TableSortOptions[] }) => {
  return request({
    url: `${controller}/list`,
    method: 'POST',
    data: params
  }) as unknown as Promise<PageData<PlanKey>>
}

export const updateStatus = (id: number, status: boolean) => {
  return request({
    url: `${controller}/updateStatus`,
    method: 'POST',
    data: {id, status}
  }) as unknown as Promise<boolean>
}

export const verify = (key: string) => {
  return request({
    url: `${controller}/verify`,
    method: 'GET',
    params: {
      key
    }
  }) as unknown as Promise<PlanType | undefined>
}
