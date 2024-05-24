import Vue from "vue";
import Vuex from "vuex";
import Snackbar from "@/store/modules/snackbar";
import Spinner from "@/store/modules/spinner";

Vue.use(Vuex);
export default () =>
  new Vuex.Store({
    modules: {
      Snackbar,
      Spinner,
    },
  });
