<template>
  <v-dialog v-model="helpDetail">
    <template v-slot:activator="{ on, attrs }">
      <v-btn elevation="2" icon text v-bind="attrs" v-on="on" small>
        <v-icon x-small>mdi-help</v-icon>
      </v-btn>
    </template>
    <v-card>
      <v-card-title class="text-h5">{{ title }}</v-card-title>
      <v-card-text>
        <v-container>
          <v-row>
            <v-col cols="12">
              <div v-html="help[$i18n.locale]"></div>
            </v-col>
          </v-row>
        </v-container>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="blue darken-1" text @click="closeHelp">{{
          $t("stores.close")
        }}</v-btn>
        <v-spacer></v-spacer>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">
import { Component, Vue, Prop, Watch } from "nuxt-property-decorator";
import { namespace } from "vuex-class";
import IResponseTexts from "~/model/IResponseTexts";
const snackbar = namespace("Snackbar");
const spinner = namespace("Spinner");

@Component
export default class Help extends Vue {
  @Prop() readonly title!: string;
  @Prop() readonly keyString!: string;

  @snackbar.Action
  public updateText!: (newText: string) => void;

  @snackbar.Action
  public updateColor!: (newColor: string) => void;

  @snackbar.Action
  public updateShow!: (newShow: boolean) => void;

  @spinner.Action
  public toggleSpinner!: (newShow: boolean) => void;

  $axios: any;
  $i18n: any;
  helpDetail: boolean = false;
  help: any = {};
  $t: any;

  @Watch("helpDetail")
  onPropertyChanged(value: boolean, oldValue: boolean) {
    if (value) {
      this.toggleSpinner(true);
      const endpoint = this.$config.cmsApi + "/" + this.$config.token + "text";
      this.$axios.get(endpoint).then((response: IResponseTexts) => {
        if (response.data.success) {
          this.help = JSON.parse(
            response.data.texts.filter((x) => x.Key === this.keyString)[0].Value
          );
          this.toggleSpinner(false);
        } else {
          this.toggleSpinner(false);
          this.updateText(response.data.error);
          this.updateColor("red");
          this.updateShow(true);
        }
      });
    }
  }

  closeHelp() {
    this.helpDetail = false;
  }
}
</script>

<style scoped>
.container {
  color: white;
}
</style>
