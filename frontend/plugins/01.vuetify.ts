import { createVuetify } from "vuetify";
import * as components from "vuetify/components";
import * as directives from "vuetify/directives";

export default defineNuxtPlugin((nuxtApp) => {
  const vuetify = createVuetify({
    components,
    directives,
    theme: {
      defaultTheme: "light",
      themes: {
        dark: {
          dark: true,
          colors: {
            primary: "#FF8F00",
            accent: "#424242",
            secondary: "#FF8F00",
            info: "#26A69A",
            warning: "#FFC107",
            error: "#DD2C00",
            success: "#00E676",
          },
        },
        light: {
          dark: false,
          colors: {
            primary: "#B1EA4E",
            accent: "#000000",
            secondary: "#00C9DB",
            info: "#26A69A",
            warning: "#FFC107",
            error: "#DD2C00",
            success: "#00E676",
          },
        },
      },
    },
  });

  nuxtApp.vueApp.use(vuetify);

  const vuetifyCompat = {
    ...vuetify,
    theme: {
      ...vuetify.theme,
      get dark() {
        return vuetify.theme.global.name.value === "dark";
      },
      set dark(val: boolean) {
        vuetify.theme.global.name.value = val ? "dark" : "light";
      },
    },
  };

  nuxtApp.vueApp.config.globalProperties.$vuetify = vuetifyCompat;
});
