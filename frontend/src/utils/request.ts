/* eslint-disable no-console */
/* eslint-disable @typescript-eslint/no-explicit-any */
/*
 * 通信核心
 * 负责后端接口调用、错误捕获、身份信息绑定、通信配置和异常处理。
 */
import { LOGIN_REDIRECT_KEY } from '@/config'
import { passport } from '@/utils/passport'
import { message } from 'antd'
import axios, { type Method } from 'axios'

declare global {
  interface Window {
    cancelRequest?: (message?: string) => void
  }
}

export const baseURL = window.location.origin + '/api'

axios.defaults.baseURL = baseURL
axios.defaults.timeout = 1000 * 120

axios.interceptors.request.use(
  config => {
    if (passport.isValid()) config.headers.Authorization = 'Bearer ' + passport.getToken()
    return config
  },
  err => Promise.reject(err)
)

axios.interceptors.response.use(
  res => res,
  error => Promise.reject(error)
)

type TArgs = [data: any, url: string, method?: Method, progressCallback?: (e: any) => void]

export const CANCEL_STACK: Map<string, (message?: string) => void> = new Map()

function CommonAjax(...args: TArgs): Promise<any> {
  const [data, url, method, progressCallback] = args
  let payload: any = { data }

  if (!method || method.toLocaleLowerCase() === 'get') payload = { params: data }

  const params: any = {
    method: method || 'get',
    url,
    ...payload,
    cancelToken: new axios.CancelToken(cancelRequest => {
      CANCEL_STACK.set(url, cancelRequest)
      window.cancelRequest = cancelRequest
    }),
  }

  if (progressCallback) params.onUploadProgress = progressCallback

  return new Promise((resolve, reject) => {
    axios(params)
      .then(res => resolve(res.data))
      .catch(err => {
        if (err.message === 'Network Error') {
          reject(Object.assign(err, { message: '网络异常，请检查网络连接' }))
          return
        }

        if (typeof err.response !== 'undefined' && typeof err.response.status !== 'undefined') {
          const code = err.response.status

          switch (code) {
            case 401: {
              sessionStorage.setItem(
                LOGIN_REDIRECT_KEY,
                location.hash.includes('sign-in') ? '/log/top-page' : location.pathname
              )
              passport.clear()
              location.href = '/'
              break
            }
            case 403:
              console.warn('无权限')
              break
            case 404:
              console.error('===========访问的资源不存在，请检查===========')
              break
            case 500:
              console.error('===========服务器罢工============')
              message.error(err?.response?.data?.msg || err?.response?.data?.message)
              console.error(`data:${JSON.stringify(data)}, url:${url}, method:${method}`)
              break
            default: {
              const responseData = err.response.data
              const msg: string =
                responseData && responseData.msg
                  ? responseData.msg[0] || responseData.msg
                  : responseData.message || '未预料的异常'
              message.error(msg)
              console.error(`${code} 未做错误代码的捕获的后端端异常`)
              break
            }
          }
        }

        if (err.response?.status !== 401) {
          reject(err)
        }
      })
      .finally(() => {
        CANCEL_STACK.delete(url)
      })
  })
}

export { CommonAjax as ajax }
