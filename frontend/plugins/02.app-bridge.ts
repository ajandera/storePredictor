import axios from "axios";
import { reactive } from "vue";
import { $store } from "~/compat/store";

export default defineNuxtPlugin(async (nuxtApp) => {
  const runtimeConfig = useRuntimeConfig();
  const router = useRouter();
  const publicConfig = runtimeConfig.public as Record<string, any>;

  const baseURL =
    publicConfig.partnerApi ||
    publicConfig.internalApi ||
    publicConfig.api ||
    "http://localhost:9999/v1/";

  const axiosInstance = axios.create({
    baseURL,
  });

  const getStored = (key: string): string | null => {
    if (typeof window === "undefined") return null;
    return window.localStorage.getItem(key);
  };

  const setStored = (key: string, val: string | null) => {
    if (typeof window === "undefined") return;
    if (val === null || val === undefined || val === "") {
      window.localStorage.removeItem(key);
    } else {
      window.localStorage.setItem(key, val);
    }
  };

  const defaultUser = () => ({
    Id: "",
    Email: "",
    Name: "",
    Role: "Admin",
    Parent: "",
    PlanRefer: "",
    CountryCode: "SK",
  });

  let initialUser = defaultUser();
  if (typeof window !== "undefined") {
    try {
      const raw = window.localStorage.getItem("sp_auth_user");
      if (raw) {
        initialUser = { ...defaultUser(), ...JSON.parse(raw) };
      }
    } catch {
      // ignore corrupted storage
    }
  }

  const authState = reactive({
    loggedIn: Boolean(getStored("sp_access_token")),
    user: initialUser,
  });

  axiosInstance.interceptors.request.use((config) => {
    const token = getStored("sp_access_token");
    if (token) {
      config.headers = config.headers || {};
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  });

  let isRefreshing = false;
  axiosInstance.interceptors.response.use(
    (res) => res,
    async (error) => {
      const originalRequest = error?.config;
      if (
        error?.response?.status === 401 &&
        originalRequest &&
        !originalRequest._retry &&
        !String(originalRequest.url || "").includes("/auth") &&
        !String(originalRequest.url || "").includes("/refresh")
      ) {
        const refreshToken = getStored("sp_refresh_token");
        if (refreshToken && !isRefreshing) {
          isRefreshing = true;
          originalRequest._retry = true;
          try {
            const refreshRes = await axiosInstance.post("/refresh", {
              refresh_token: refreshToken,
            });
            if (refreshRes.data?.jwt?.access_token) {
              setStored("sp_access_token", refreshRes.data.jwt.access_token);
              if (refreshRes.data.jwt.refresh_token) {
                setStored("sp_refresh_token", refreshRes.data.jwt.refresh_token);
              }
              isRefreshing = false;
              originalRequest.headers = originalRequest.headers || {};
              originalRequest.headers.Authorization = `Bearer ${refreshRes.data.jwt.access_token}`;
              return axiosInstance(originalRequest);
            }
          } catch {
            // refresh failed
          }
          isRefreshing = false;
        }
      }
      return Promise.reject(error);
    }
  );

  const $auth = {
    get loggedIn() {
      return authState.loggedIn;
    },
    set loggedIn(val: boolean) {
      authState.loggedIn = val;
    },
    get user() {
      return authState.user;
    },
    set user(val: any) {
      authState.user = val || defaultUser();
    },
    $state: authState,
    strategy: {
      options: {
        name: "local",
      },
    },
    $storage: {
      getUniversal(key: string) {
        if (key === "user") return authState.user;
        return getStored(key);
      },
      setUniversal(key: string, val: any) {
        if (key === "user") {
          authState.user = val;
          setStored("sp_auth_user", JSON.stringify(val));
        } else {
          setStored(key, typeof val === "string" ? val : JSON.stringify(val));
        }
      },
    },
    setUserToken(accessToken: string, refreshToken?: string) {
      setStored("sp_access_token", accessToken);
      if (refreshToken) {
        setStored("sp_refresh_token", refreshToken);
      }
      authState.loggedIn = Boolean(accessToken);
    },
    setUser(user: any) {
      authState.user = { ...defaultUser(), ...(user || {}) };
      setStored("sp_auth_user", JSON.stringify(authState.user));
    },
    async fetchUser() {
      const res = await axiosInstance.get("/account/me");
      if (res.data?.success && res.data?.user) {
        this.setUser(res.data.user);
        authState.loggedIn = true;
      }
      return res;
    },
    async loginWith(_strategy: string, payload: { data: { email: string; password: string } }) {
      const res = await axiosInstance.post("/auth", payload.data);
      if (!res.data?.success || !res.data?.jwt?.access_token) {
        throw new Error(res.data?.error || "Invalid credentials");
      }
      this.setUserToken(
        res.data.jwt.access_token,
        res.data.jwt.refresh_token
      );
      await this.fetchUser();
      return res;
    },
    async logout() {
      try {
        await axiosInstance.post("/logout");
      } catch {
        // ignore logout network errors
      }
      setStored("sp_access_token", null);
      setStored("sp_refresh_token", null);
      setStored("sp_auth_user", null);
      authState.user = defaultUser();
      authState.loggedIn = false;
    },
  };

  if (authState.loggedIn && !authState.user.Id) {
    try {
      await $auth.fetchUser();
    } catch {
      await $auth.logout();
    }
  }

  // Event bus on $nuxt ($on, $emit, $off) + $options.router
  const listeners = new Map<string, Set<Function>>();
  const $nuxtCompat = Object.assign(nuxtApp, {
    $on(event: string, cb: Function) {
      if (!listeners.has(event)) listeners.set(event, new Set());
      listeners.get(event)!.add(cb);
    },
    $off(event: string, cb?: Function) {
      if (!cb) {
        listeners.delete(event);
      } else {
        listeners.get(event)?.delete(cb);
      }
    },
    $emit(event: string, ...args: any[]) {
      listeners.get(event)?.forEach((cb) => cb(...args));
    },
    $options: {
      router,
    },
  });

  // Wrap $i18n so $i18n.locale works as a string getter/setter and $i18n.t works as a function
  const vueI18n: any = (nuxtApp as any).$i18n;
  const i18nCompat = {
    get locale(): string {
      if (!vueI18n) return "en";
      return typeof vueI18n.locale === "string"
        ? vueI18n.locale
        : vueI18n.locale?.value || "en";
    },
    set locale(val: string) {
      if (!vueI18n) return;
      if (typeof vueI18n.setLocale === "function") {
        vueI18n.setLocale(val);
      } else if (vueI18n.locale && typeof vueI18n.locale === "object") {
        vueI18n.locale.value = val;
      } else {
        vueI18n.locale = val;
      }
    },
    t(key: string, ...args: any[]) {
      if (vueI18n && typeof vueI18n.t === "function") {
        return vueI18n.t(key, ...args);
      }
      return key;
    },
  };

  const gp = nuxtApp.vueApp.config.globalProperties;
  gp.$axios = axiosInstance;
  gp.$auth = $auth;
  gp.$store = $store;
  gp.$config = publicConfig;
  gp.$nuxt = $nuxtCompat;
  gp.$i18n = i18nCompat;
  gp.$t = (key: string, ...args: any[]) => i18nCompat.t(key, ...args);
  if (!gp.localePath) {
    gp.localePath = (path: string) => path;
  }
});
