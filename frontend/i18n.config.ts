import en from "./locales/en.json";
import cz from "./locales/cz.json";
import sk from "./locales/sk.json";
import de from "./locales/de.json";
import pl from "./locales/pl.json";
import ro from "./locales/ro.json";
import hu from "./locales/hu.json";

export default defineI18nConfig(() => ({
  legacy: false,
  locale: "en",
  fallbackLocale: "en",
  messages: {
    en,
    cz,
    cs: cz,
    sk,
    de,
    pl,
    ro,
    hu,
  },
}));
