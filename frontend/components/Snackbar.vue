<template>
  <v-snackbar v-model="showSnackbar" :timeout="snackTimeout" :color="color">
    {{ text }}
  </v-snackbar>
</template>

<script lang="ts">
import { Component, Vue } from "nuxt-property-decorator";
import { namespace } from "vuex-class";
const snackbar = namespace("Snackbar");

@Component
export default class Snackbar extends Vue {
  snackTimeout: number = 2000;

  get showSnackbar() {
    return this.$store.state.Snackbar.show === true;
  }

  set showSnackbar(value: boolean) {
    this.updateShow(value);
  }

  @snackbar.State
  public text!: string;

  @snackbar.State
  public color!: string;

  @snackbar.State
  public show!: boolean;

  @snackbar.Action
  public updateShow!: (newShow: boolean) => void;
}
</script>
