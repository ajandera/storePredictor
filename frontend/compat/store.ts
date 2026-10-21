import { reactive } from "vue";

export const storeState = reactive({
  Snackbar: {
    text: "",
    color: "",
    show: false,
  },
  Spinner: {
    show: false,
  },
});

export const storeActions: Record<string, Record<string, (...args: any[]) => any>> = {
  Snackbar: {
    updateText(newText: string) {
      storeState.Snackbar.text = newText;
    },
    updateColor(newColor: string) {
      storeState.Snackbar.color = newColor;
    },
    updateShow(newVisibility: boolean) {
      storeState.Snackbar.show = newVisibility;
    },
  },
  Spinner: {
    toggleSpinner(newVisibility: boolean) {
      storeState.Spinner.show = newVisibility;
    },
  },
};

export const $store = {
  state: storeState,
  commit(type: string, payload?: any) {
    if (type === "Snackbar/setText" || type === "setText") storeState.Snackbar.text = payload;
    if (type === "Snackbar/setColor" || type === "setColor") storeState.Snackbar.color = payload;
    if (type === "Snackbar/setShow") storeState.Snackbar.show = payload;
    if (type === "Spinner/setShow") storeState.Spinner.show = payload;
  },
  dispatch(type: string, payload?: any) {
    const [ns, act] = type.includes("/") ? type.split("/") : ["", type];
    if (ns && storeActions[ns]?.[act]) {
      return storeActions[ns][act](payload);
    }
  },
};

// Shared metadata queue for both Legacy and Stage 3 decorators
export const pendingMeta = {
  props: {} as Record<string, any>,
  watchers: [] as Array<{ path: string; handler: string; options: any }>,
  vuexActions: {} as Record<string, Function>,
  vuexStates: {} as Record<string, () => any>,
  reset() {
    this.props = {};
    this.watchers = [];
    this.vuexActions = {};
    this.vuexStates = {};
  },
};
