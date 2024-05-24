<template>
  <div class="preloaderWrapper" v-if="show === true">
    <v-progress-circular
      :size="50"
      color="amber"
      indeterminate
    ></v-progress-circular>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "nuxt-property-decorator";
import { namespace } from "vuex-class";
const snackbar = namespace("Spinner");

@Component
export default class Snackbar extends Vue {
  @snackbar.State
  public show!: boolean;

  @snackbar.Action
  public toggleSpinner!: (newShow: boolean) => void;
}
</script>
<style scoped>
.v-progress-circular {
  padding-top: 45%;
  margin-left: 48%;
}
.preloaderWrapper {
  z-index: 999;
  width: 100%;
  height: 100%;
  position: fixed;
  background: #2828285e;
  top: 0;
  right: 0;
}
</style>
