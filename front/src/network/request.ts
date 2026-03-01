import axios, {AxiosRequestConfig, AxiosResponse} from "axios";
import {ReqResponse} from "@/types";
import {local} from "@/utils/user";
import router from "@/router";
import {useGlobalDialog} from "@/stores/dialog";

const devServer = 'http://localhost:8282/'
const prodServer = '/api/'

const {showGlobalDialog} = useGlobalDialog()

//提示： 可参考视频资料：https://www.bilibili.com/video/BV15741177Eh?p=155
export function request(config: AxiosRequestConfig) {
  const isDev = import.meta.env.MODE === "development";
  const instance = axios.create({
    //根路径
    baseURL: isDev ? devServer : prodServer,
    timeout: 30 * 1000,
  })
  //axios 拦截器
  //请求拦截
  instance.interceptors.request.use(request => {
    request.headers.setAuthorization(local.token)
    return request
  }, err => {
    console.log("err", err)
  })
  instance.interceptors.response.use((res: AxiosResponse<ReqResponse>): any => {
    const {code, data, msg} = res.data as ReqResponse;
    console.log(code, msg)
    return data
  }, err => {
    const status = err.status;
    try {
      const {code, msg} = err.response.data as ReqResponse;
      if (config.url?.startsWith("/admin/logout")) {
        router.push("/login");
        return
      }
      if (status == 401) {
        router.push("/login");
      }
      showGlobalDialog({
        iconColor: 'error',
        msg: msg || err.message,
        title: `Error: ${code || status}`,
      })
    } catch (_) {
      showGlobalDialog({
        iconColor: 'error',
        msg: err.message,
        title: `Error`,
      })
    }
    throw err
  })
  // 发送真正的网络请求
  return instance(config)
}
