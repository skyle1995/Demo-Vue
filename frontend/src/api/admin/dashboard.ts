import { http } from "@/utils/http";

export type SystemInfoResult = {
  code: number;
  msg: string;
  data: {
    version: string;
    mode: boolean;
    db_type: string;
    uptime: string;
    uptime_seconds: number;
  };
};

export type LoginLogsResult = {
  code: number;
  msg: string;
  data: {
    total: number;
    list: Array<any>;
  };
};

export const getSystemInfo = () => {
  return http.request<SystemInfoResult>("get", "/api/admin/system/info");
};

export const getLoginLogs = (params?: object) => {
  return http.request<LoginLogsResult>(
    "get",
    "/api/admin/dashboard/login-logs",
    { params }
  );
};
