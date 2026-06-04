import { defineStore } from "pinia";
import {
  type userType,
  store,
  router,
  resetRouter,
  routerArrays,
  storageLocal
} from "../utils";
import {
  type UserResult,
  type RefreshTokenResult,
  getLogin,
  refreshTokenApi
} from "@/api/admin/user";
import { useMultiTagsStoreHook } from "./multiTags";
import { type DataInfo, setToken, removeToken, getUserKey, isUserMode } from "@/utils/auth";

export const useUserStore = defineStore("pure-user", {
  state: (): userType => ({
    // 头像
    avatar: storageLocal().getItem<DataInfo<number>>(getUserKey())?.avatar ?? "",
    // 用户名
    username: storageLocal().getItem<DataInfo<number>>(getUserKey())?.username ?? "",
    // 昵称
    nickname: storageLocal().getItem<DataInfo<number>>(getUserKey())?.nickname ?? "",
    // 页面级别权限
    roles: storageLocal().getItem<DataInfo<number>>(getUserKey())?.roles ?? [],
    // 按钮级别权限
    permissions:
      storageLocal().getItem<DataInfo<number>>(getUserKey())?.permissions ?? [],
    // 是否勾选了登录页的免登录
    isRemembered: false,
    // 登录页的免登录存储几天，默认7天
    loginDay: 7
  }),
  actions: {
    /** 存储头像（同步刷新本地缓存，避免页面刷新后回退） */
    SET_AVATAR(avatar: string) {
      this.avatar = avatar;
      const userInfo =
        storageLocal().getItem<DataInfo<number>>(getUserKey()) ?? ({} as any);
      userInfo.avatar = avatar;
      storageLocal().setItem(getUserKey(), userInfo);
    },
    /** 存储用户名 */
    SET_USERNAME(username: string) {
      this.username = username;
      const userInfo =
        storageLocal().getItem<DataInfo<number>>(getUserKey()) ?? ({} as any);
      userInfo.username = username;
      storageLocal().setItem(getUserKey(), userInfo);
    },
    /** 存储昵称 */
    SET_NICKNAME(nickname: string) {
      this.nickname = nickname;
      const userInfo =
        storageLocal().getItem<DataInfo<number>>(getUserKey()) ?? ({} as any);
      userInfo.nickname = nickname;
      storageLocal().setItem(getUserKey(), userInfo);
    },
    /** 存储角色 */
    SET_ROLES(roles: Array<string>) {
      this.roles = roles;
    },
    /** 存储按钮级别权限 */
    SET_PERMS(permissions: Array<string>) {
      this.permissions = permissions;
    },
    /** 存储是否勾选了登录页的免登录 */
    SET_ISREMEMBERED(bool: boolean) {
      this.isRemembered = bool;
    },
    /** 设置登录页的免登录存储几天 */
    SET_LOGINDAY(value: number) {
      this.loginDay = Number(value);
    },
    /** 登入 */
    async loginByUsername(data) {
      return new Promise<UserResult>((resolve, reject) => {
        getLogin(data)
          .then(data => {
            if (data?.success) setToken(data.data);
            resolve(data);
          })
          .catch(error => {
            reject(error);
          });
      });
    },
    /** 前端登出（不调用接口） */
    logOut() {
      this.username = "";
      this.roles = [];
      this.permissions = [];
      removeToken();
      useMultiTagsStoreHook().handleTags("equal", [...routerArrays]);
      resetRouter();
      const loginPath = isUserMode() ? "/user/login" : "/admin/login";
      router.push(loginPath);
    },
    /** 刷新`token` */
    async handRefreshToken(data) {
      return new Promise<RefreshTokenResult>((resolve, reject) => {
        refreshTokenApi(data)
          .then(data => {
            if (data) {
              setToken(data.data);
              resolve(data);
            }
          })
          .catch(error => {
            reject(error);
          });
      });
    }
  }
});

export function useUserStoreHook() {
  return useUserStore(store);
}
