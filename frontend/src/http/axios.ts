import Axios from "axios";
import type { HttpRequest } from "./types";
import type { AxiosResponse } from "axios";

const instance = Axios.create({
  baseURL: "http://localhost:3000",
  timeout: 3000,
  timeoutErrorMessage: "Error connection timeout",
});

export async function httpRequest<T>(
  r: HttpRequest
): Promise<AxiosResponse<T>> {
  return await instance.request<T>(r);
}
