<template>
  <v-alert
      icon="mdi-shield-lock-outline"
      prominent
      text
      type="info"
    >
    <v-row align="center">
      <v-col class="grow">
        {{ $t('upgrade.text') }}
      </v-col>
      <v-col class="shrink">
        <v-dialog v-model="dialog" max-width="500px">
          <template v-slot:activator="{ on: dialog, attrs }">
            <v-btn v-bind="attrs" v-on="dialog">{{ $t('upgrade.btn') }}</v-btn>
          </template>
          <v-card>
            <v-card-title>
              <span class="text-h5">{{ $t("upgrade.title") }}</span>
            </v-card-title>
            <v-card-text>
              <v-container v-if="pay.account === undefined">
                <h3>{{ $t("upgrade.description") }}</h3>
                <v-select
                  v-model="plan"
                  :items="plans"
                  :label="$t('upgrade.selectPlan')"
                  item-text="text"
                  item-value="Id"
                  @change="handleSubmit(plan)"
                  required
                ></v-select>
                <v-checkbox
                    v-model="year"
                    :label="$t('upgrade.year')"
                  ></v-checkbox>

              </v-container>
              <v-container v-if="pay.account !== undefined">
                <no-ssr>
                  <paypal-checkout
                    class="mt-5"
                    :amount="pay.amount"
                    :currency="pay.currency"
                    :env="payPalEnv"
                    :client="credentials"
                    :invoice-number="pay.account"
                    :braintree="braintreeSdk"
                  >
                  </paypal-checkout>
                </no-ssr>
              </v-container>
            </v-card-text>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="blue darken-1" text @click="dialog = false">
                {{ $t('index.cancel') }}
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
      </v-col>
    </v-row>
    </v-alert>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "nuxt-property-decorator";
import { namespace } from "vuex-class";
import IPlan from "~/model/IPlan";
import Store from "~/model/Store";
const snackbar = namespace("Snackbar");
const spinner = namespace("Spinner");

@Component
export default class Upgrade extends Vue {
  @Prop() readonly store!: Store;
  @Prop() readonly plans!: IPlan[];
  
  @snackbar.Action
  public updateText!: (newText: string) => void;

  @snackbar.Action
  public updateColor!: (newColor: string) => void;

  @snackbar.Action
  public updateShow!: (newShow: boolean) => void;

  @spinner.Action
  public toggleSpinner!: (newShow: boolean) => void;

  email: string = "";
  dialog: boolean = false;
  $axios: any;
  $i18n: any;
  $t: any;
  plan: IPlan = {
    Name: "",
    Price: 0,
    Period: 0,
    Id: ""
  };
  pay: any = {};
  payByPaypal: boolean = true;
  year: boolean = true;
  credentials: object = {
    sandbox: this.$config.$paypalSandboxClientId,
    production: this.$config.$paypalClientId,
  };
  payPalEnv: string = this.$config.$paypalEnv;
  // @ts-ignore
  braintreeSdk: any = window.braintree;
  data: any = {};
  $auth: any;

  handleSubmit(plan: IPlan) {
    this.toggleSpinner(true);
    console.log(plan);
    this.$axios
      .post(this.$config.internalApi + "upgrade/"+ this.$auth.user.Id, {
        PlanRefer: plan,
        Year: this.year
      })
      .then((response: any) => {
        if (response.data.success === true) {
          const period = this.year ? 12 : 1;
          //@ts-ignore
          const price = (this.plans.find((item: IPlan) => item.Id === plan).Price * period);
          const payData = {
            account: this.$auth.user.Id,
            amount: price.toString(),
            currency: "EUR",
          };
          this.data = payData;
          this.pay = payData;
          this.toggleSpinner(false);
        } else {
          this.updateText(response.data.error);
          this.updateColor("red");
          this.updateShow(true);
          this.toggleSpinner(false);
        }
      });
  }
}
</script>

<style scoped>

</style>
