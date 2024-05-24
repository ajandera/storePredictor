import colors from "vuetify/es5/util/colors";
import en from "./locales/en.json";
import cz from "./locales/cz.json";
import sk from "./locales/sk.json";
import de from "./locales/de.json";
import pl from "./locales/pl.json";
import ro from "./locales/ro.json";
import hu from "./locales/hu.json";

export default {
  // Disable server-side rendering: https://go.nuxtjs.dev/ssr-mode
  ssr: false,

  // Target: https://go.nuxtjs.dev/config-target
  target: "static",

  // Global page headers: https://go.nuxtjs.dev/config-head
  head: {
    titleTemplate: "%s | storePredictor",
    title: "Dashboard",
    htmlAttrs: {
      lang: "en",
    },
    meta: [
      { charset: "utf-8" },
      { name: "viewport", content: "width=device-width, initial-scale=1" },
      { hid: "description", name: "description", content: "" },
      { name: "format-detection", content: "telephone=no" },
    ],
    link: [{ rel: "icon", type: "image/x-icon", href: "/favicon/favicon.ico" }]
  },

  // Global CSS: https://go.nuxtjs.dev/config-css
  css: [],

  // Plugins to run before rendering page: https://go.nuxtjs.dev/config-plugins
  plugins: [
    { src: "~/plugins/format-date.js" },
    { src: "~/plugins/apexcharts.js" },
    { src: "~/plugins/vue-quill-editor.js" },
    { src: "~/plugins/country-flag.js" }
  ],

  // Auto import components: https://go.nuxtjs.dev/config-components
  components: true,

  // Modules for dev and build (recommended): https://go.nuxtjs.dev/config-modules
  buildModules: [
    // https://go.nuxtjs.dev/typescript
    "@nuxt/typescript-build",
    // https://go.nuxtjs.dev/vuetify
    "@nuxtjs/vuetify",
  ],

  // Modules: https://go.nuxtjs.dev/config-modules
  modules: ["@nuxtjs/i18n", "@nuxtjs/axios", "@nuxtjs/auth-next", "nuxt-mermaid-string", "endent", "tinycolor2"],

  /*
   ** Axios module configuration
   */
  axios: {
    // See https://github.com/nuxt-community/axios-module#options
  },

  i18n: {
    strategy: "prefix_except_default",
    locales: ["en", "cz", "sk", "de", "pl", "ro", "hu"],
    defaultLocale: "en",
    vueI18n: {
      fallbackLocale: "en",
      messages: { en, cz, sk, de, pl, ro, hu },
    },
  },

  // Vuetify module configuration: https://go.nuxtjs.dev/config-vuetify
  vuetify: {
    customVariables: ["~/assets/variables.scss"],
    theme: {
      dark: false,
      themes: {
        dark: {
          primary: colors.amber.darken3,
          accent: colors.grey.darken3,
          secondary: colors.amber.darken3,
          info: colors.teal.lighten1,
          warning: colors.amber.base,
          error: colors.deepOrange.accent4,
          success: colors.green.accent3,
        },
        light: {
          primary: "#B1EA4E",
          accent: colors.shades.black,
          secondary: "#00C9DB",
          info: colors.teal.lighten1,
          warning: colors.amber.base,
          error: colors.deepOrange.accent4,
          success: colors.green.accent3,
        },
      },
    },
  },

  auth: {
    redirect: {
      login: '/sign',
      logout: '/',
      callback: '/callback',
      home: '/'
    },
    strategies: {
      local: {
        scheme: "refresh",
        token: {
          property: "jwt.access_token",
          maxAge: 60,
          global: true,
          type: "Bearer",
        },
        refreshToken: {
          property: "jwt.refresh_token",
          data: "refresh_token",
          maxAge: 60 * 60 * 24,
        },
        user: {
          property: "user",
        },
        endpoints: {
          login: { url: "/auth", method: "post" },
          refresh: { url: "/refresh", method: "post" },
          user: { url: "/account/me", method: "get" },
          logout: { url: "/logout", method: "post" },
        },
      }
    },
  },

  // Build Configuration: https://go.nuxtjs.dev/config-build
  build: {
    vendor: ["vue-apexchart"],
  },

  publicRuntimeConfig: {
    api: process.env.API,
    setting: {
      cz: { date: "cs-CS", currency: "CZK", symbol: "Kč" },
      eu: { date: "sk-SK", currency: "EUR", symbol: "€" },
      sk: { date: "sk-SK", currency: "EUR", symbol: "€" },
    },
    axios: {
      baseURL: process.env.PARTNER_API,
    },
  },
  generate: { fallback: true }
};
