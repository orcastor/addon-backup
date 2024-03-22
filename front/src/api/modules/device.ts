import { Device } from "@/api/interface/index";
import { BAK_API } from "@/api/config/config";

import http from "@/api";

export const listApi = (params: Device.ReqList) => {
  return http.post<Device.ResList>(BAK_API + `/list`, params);
};
