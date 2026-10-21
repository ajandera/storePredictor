import { fileURLToPath } from "node:url";
import vuetify, { transformAssetUrls } from "vite-plugin-vuetify";

const defaultApi =
  process.env.NUXT_PUBLIC_API ||
  process.env.API ||
  process.env.PARTNER_API ||
  "http://localhost:9999/v1/";

export default defineNuxtConfig({
  compatibilityDate: "2025-01-01",
  ssr: false,
  srcDir: ".",
  dir: {
    public: "public",
  },

  app: {
    head: {
      titleTemplate: "%s | storePredictor",
      title: "Dashboard",
      htmlAttrs: {
        lang: "en",
      },
      meta: [
        { charset: "utf-8" },
        { name: "viewport", content: "width=device-width, initial-scale=1" },
        { name: "description", content: "" },
        { name: "format-detection", content: "telephone=no" },
      ],
      link: [{ rel: "icon", type: "image/x-icon", href: "/favicon/favicon.ico" }],
    },
  },

  css: ["vuetify/styles", "@mdi/font/css/materialdesignicons.css"],

  build: {
    transpile: ["vuetify"],
  },

  components: true,

  modules: [
    "@nuxtjs/i18n",
    (_options, nuxt) => {
      nuxt.hooks.hook("vite:extendConfig", (config) => {
        config.plugins = config.plugins || [];
        config.plugins.push(vuetify({ autoImport: true }));
      });
    },
  ],

  i18n: {
    restructureDir: false,
    strategy: "prefix_except_default",
    locales: ["en", "cz", "cs", "sk", "de", "pl", "ro", "hu"],
    defaultLocale: "en",
    vueI18n: "./i18n.config.ts",
    bundle: {
      optimizeTranslationDirective: false,
    },
  },

  alias: {
    "nuxt-property-decorator": fileURLToPath(
      new URL("./compat/nuxt-property-decorator.ts", import.meta.url)
    ),
    "vuex-class": fileURLToPath(new URL("./compat/vuex-class.ts", import.meta.url)),
  },

  vite: {
    vue: {
      template: {
        transformAssetUrls,
      },
    },
  },

  runtimeConfig: {
    public: {
      api: defaultApi,
      internalApi:
        process.env.NUXT_PUBLIC_INTERNAL_API ||
        process.env.INTERNAL_API ||
        defaultApi,
      partnerApi:
        process.env.NUXT_PUBLIC_PARTNER_API ||
        process.env.PARTNER_API ||
        defaultApi,
      cmsApi: process.env.NUXT_PUBLIC_CMS_API || process.env.CMS_API || "",
      token: process.env.NUXT_PUBLIC_TOKEN || process.env.TOKEN || "",
      $paypalSandboxClientId: process.env.PAYPAL_SANDBOX_CLIENT_ID || "",
      $paypalClientId: process.env.PAYPAL_CLIENT_ID || "",
      $paypalEnv: process.env.PAYPAL_ENV || "sandbox",
      setting: {
        cz: { date: "cs-CS", currency: "CZK", symbol: "Kč" },
        eu: { date: "sk-SK", currency: "EUR", symbol: "€" },
        sk: { date: "sk-SK", currency: "EUR", symbol: "€" },
      },
    },
  },
});
