import axios, { AxiosRequestConfig, AxiosResponse, CancelToken } from "axios";

export const baseURL = "http://localhost:8080";

const getPromiseConfig = (
  cancelToken: CancelToken | null = null
): AxiosRequestConfig => {
  let config: AxiosRequestConfig = {
    baseURL: baseURL,
    withCredentials: false,
  };
  if (cancelToken) {
    config = { ...config, cancelToken };
  }
  return config;
};

async function ajaxPromise<ResponseData = unknown, Payload = unknown>(
  url: string,
  type: "GET" | "POST" | "PUT" | "DELETE",
  data?: Payload,
  cancelToken: CancelToken | null = null,
  timeout = 0
) {
  const config = getPromiseConfig(cancelToken);
  if (timeout) {
    config.timeout = timeout;
  }
  const resp: AxiosResponse<ResponseData> = await (type === "POST"
    ? axios.post(url, data, config)
    : type === "GET"
    ? axios.get(url, config)
    : type === "PUT"
    ? axios.put(url, data, config)
    : type === "DELETE"
    ? axios.delete(url, config)
    : axios.get(url, config));

  return resp;
}
export default ajaxPromise;
export { getPromiseConfig as getConfig };
