import {request} from '../request'
import {PlanType} from "@/types";

const controller = '/admin/plan'
export const addPlan = (plan: PlanType) => {
  return request({
    url: `${controller}/add`,
    method: 'post',
    data: plan
  }) as unknown as Promise<boolean>
}
export const editPlan = (plan: PlanType) => {
  return request({
    url: `${controller}/edit`,
    method: 'post',
    data: plan
  }) as unknown as Promise<boolean>
}

export const getPlans = () => {
  return request({
    url: `${controller}/list`,
    method: 'get',
  }) as unknown as Promise<PlanType[]>
}
export const updateStatus = (id: string, status: boolean) => {
  return request({
    url: `${controller}/updateStatus`,
    method: 'post',
    data: {
      id,
      status
    }
  }) as unknown as Promise<boolean>
}

export const generatePlanKeys = (id: string, size: number) => {
  return request({
    url: `${controller}/generatePlanKeys`,
    method: 'post',
    data: {
      id,
      size
    }
  }) as unknown as Promise<boolean>
}
