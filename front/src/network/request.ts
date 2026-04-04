import { t } from '@/i18n'
import router from '@/router'
import { useGlobalDialog } from '@/stores/dialog'
import { ReqResponse } from '@/types'
import { local } from '@/utils/user'
import axios, { AxiosRequestConfig, AxiosResponse } from 'axios'

const devServer = 'http://localhost:8282/'
const prodServer = '/api/'

const { showGlobalDialog } = useGlobalDialog()

export function request (config: AxiosRequestConfig) {
  const isDev = import.meta.env.MODE === 'development'
  const instance = axios.create({
    baseURL: isDev ? devServer : prodServer,
    timeout: 30 * 1000,
  })

  instance.interceptors.request.use(request => {
    request.headers.setAuthorization(local.token)
    return request
  }, err => {
    console.log('err', err)
  })

  instance.interceptors.response.use((res: AxiosResponse<ReqResponse>): any => {
    const { code, data, msg } = res.data as ReqResponse
    console.log(code, msg)
    return data
  }, err => {
    const status = err.status
    try {
      const { code, msg } = err.response.data as ReqResponse
      if (config.url?.startsWith('/admin/logout')) {
        router.push('/login')
        return
      }
      if (status === 401) {
        router.push('/login')
      }
      showGlobalDialog({
        iconColor: 'error',
        msg: msg || err.message,
        title: `${t('errors.title')}: ${code || status}`,
      })
    } catch (_) {
      showGlobalDialog({
        iconColor: 'error',
        msg: err.message,
        title: t('errors.title'),
      })
    }
    throw err
  })

  return instance(config)
}
