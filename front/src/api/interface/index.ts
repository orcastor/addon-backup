// * 请求响应参数(不包含data)
export interface Result {
  code: string;
  msg: string;
}

// * 请求响应参数(包含data)
export interface ResultData<T = any> extends Result {
  data?: T;
}

// * 登录模块
export namespace Login {
  export interface ReqLoginForm {
    u: string; // username
    p: string; // password
  }
  export interface ResLogin {
    access_token: string;
    u: any; // user
    b: any; // bkts
  }
}

export namespace Device {
  export interface ReqList {
  }
  export interface DeviceInfo {
    device_id: string;
    authorized: boolean;
    serial_no: string;
    connection: string;
    product_name: string;
    brand: string;
    os: string;
    total: string;
    data_available: string;
    sys_available: string;
  }
  export interface ResList {
    count: number;
    devs: DeviceInfo[];
  }
}
