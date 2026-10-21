# storePredictor Frontend (Nuxt 4)

The storePredictor client application built with **Nuxt 4** (`^4.6.2`), **Vue 3**, **Vuetify 3**, **Vite 8**, and **`@nuxtjs/i18n`**.

## Build Setup

```bash
# install dependencies
$ npm install

# serve with hot reload at http://localhost:3000
$ npm run dev

# build for production
$ npm run build

# preview production server
$ npm run preview

# generate static site in .output/public
$ npm run generate
```

For detailed explanation on how Nuxt 4 works, check out the [Nuxt documentation](https://nuxt.com/docs).

## Environment Configuration

Create a `.env` file based on `env.example`:

```env
API=http://localhost:9999/v1/
NUXT_PUBLIC_API=http://localhost:9999/v1/
```

Runtime configuration is defined in `nuxt.config.ts` under `runtimeConfig.public` (`api`, `internalApi`, `partnerApi`, `cmsApi`, `token`, `setting`).

## Directory Structure

### `assets/`
Contains uncompiled SCSS variables and styling (`assets/variables.scss`).

### `components/`
Contains reusable Vue components (`Forgot.vue`, `Help.vue`, `Snackbar.vue`, `Spinner.vue`, `Upgrade.vue`).

### `compat/`
Contains Nuxt 4 / Vue 3 class-component and reactive store compatibility modules (`nuxt-property-decorator.ts`, `vuex-class.ts`, `store.ts`).

### `layouts/`
Contains application layouts (`default.vue`, `login.vue`, `empty.vue`, `error.vue`).

### `locales/`
Contains i18n translation dictionaries (`en.json`, `cz.json`, `sk.json`, `de.json`, `pl.json`, `ro.json`, `hu.json`) loaded via `i18n.config.ts`.

### `model/`
Contains TypeScript interfaces and data models for API requests and responses.

### `pages/`
Contains application views and routes (`index.vue`, `stores.vue`, `products.vue`, `tracking.vue`, `profile.vue`, `users.vue`, `sign.vue`, `signup.vue`, `restore.vue`, `help.vue`, `token.vue`, `callback.vue`, `storetest.vue`).

### `plugins/`
Contains Nuxt 4 plugins automatically registered at startup:
- `01.vuetify.ts` — Vuetify 3 setup and theme configuration (`light` & `dark` themes).
- `02.app-bridge.ts` — Axios instance (`$axios`) with Bearer token & refresh token interceptors, authentication state (`$auth`), reactive store (`$store`), event bus (`$nuxt.$on` / `$nuxt.$emit`), and `$i18n` helpers.
- `03.components.client.ts` — Client components (`<apexchart>`, `<country-flag>`, `<quill-editor>`, `<vue-mermaid-string>`, `<paypal-checkout>`).

### `public/`
Contains static files mapped directly to the server root `/` (e.g., `/public/sp_logo.svg` -> `/sp_logo.svg`, `/public/countries.json` -> `/countries.json`).
