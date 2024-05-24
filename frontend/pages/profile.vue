<template>
  <v-container>
    <v-row>
      <v-col col="12">
        <h2>
          {{ $t("profile.account") }}
          <Help title="Profile" keyString="help-profile" />
        </h2>
        <p>{{ $t("profile.plan") }}: {{ planRefer }}</p>
        <p v-if="planRefer.free === false">
          {{ $t("profile.paid") }}: {{ paidTo }}
        </p>
        <input type="checkbox" id="newsletter" v-model="newsletter" />
        <label for="newsletter">{{ $t("signin.newsletter") }}</label>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <h2>{{ $t("profile.profile") }}</h2>
      </v-col>
    </v-row>
    <v-form ref="form" v-model="valid" lazy-validation>
      <v-row>
        <v-col cols="12" md="4">
          <v-text-field
            v-model="name"
            :rules="nameRules"
            :counter="50"
            :label="$t('profile.name')"
            required
          ></v-text-field>
        </v-col>

        <v-col cols="12" md="4">
          <v-text-field
            v-model="email"
            :rules="emailRules"
            :counter="50"
            :label="$t('profile.email')"
            required
          ></v-text-field>
        </v-col>

        <v-col cols="12" md="4">
          <v-text-field
            v-model="street"
            :rules="streetRules"
            :counter="50"
            :label="$t('profile.street')"
            required
          ></v-text-field>
        </v-col>

        <v-col cols="12" md="4">
          <v-text-field
            v-model="city"
            :rules="cityRules"
            :counter="40"
            :label="$t('profile.city')"
            required
          ></v-text-field>
        </v-col>

        <v-col cols="12" md="4">
          <v-text-field
            v-model="zip"
            :rules="zipRules"
            :counter="6"
            :label="$t('profile.zip')"
            required
          ></v-text-field>
        </v-col>

        <v-col cols="12" md="4">
          <v-autocomplete
            :items="countries"
            :filter="customFilter"
            item-text="Name"
            item-value="Code"
            :rules="countryRules"
            :label="$t('profile.countrycode')"
            v-model="countryCode"
          ></v-autocomplete>`
        </v-col>

        <v-col cols="12" md="4">
          <v-text-field
            v-model="companyNumber"
            :rules="companyNumberRules"
            :counter="15"
            :label="$t('profile.companynumber')"
            required
          ></v-text-field>
        </v-col>

        <v-col cols="12" md="4">
          <v-text-field
            v-model="vatNumber"
            :rules="companyVatRules"
            :counter="15"
            :label="$t('profile.vatnumber')"
            required
          ></v-text-field>
        </v-col>
        <v-col cols="12" md="4">
          <v-text-field
            v-model="password"
            :rules="passwordRules"
            :label="$t('profile.password')"
            type="password"
            required
          ></v-text-field>
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12">
          <v-btn
            depressed
            class="float-right"
            color="primary"
            :disabled="!valid"
            @click="saveProfile"
          >
            {{ $t("profile.save") }}
          </v-btn>
        </v-col>
      </v-row>
    </v-form>
    <v-row>
      <v-col cols-="12">
        <h4>{{ $t("profile.remove") }}</h4>
        <p>{{ $t("profile.removeDesc") }}</p>
        <v-btn depressed color="error" @click="dialogClose = true">
          {{ $t("profile.remove") }}
        </v-btn>
        <v-dialog v-model="dialogClose" max-width="500px">
          <v-card>
            <v-card-title class="text-h5">{{
              $t("profile.delete")
            }}</v-card-title>
            <v-card-text>
              <v-text-field
                v-model="accountConfirmName"
                :counter="50"
                :label="$t('profile.confirmation')"
                required
              ></v-text-field>
              <v-card-actions>
                <v-btn color="blue darken-1" text @click="closeDialogClose">{{
                  $t("stores.cancel")
                }}</v-btn>
                <v-btn
                  color="error"
                  :disabled="accountConfirmName !== email"
                  @click="closeAccountConfirm"
                  >{{ $t("stores.ok") }}</v-btn
                >
              </v-card-actions>
              <v-spacer></v-spacer>
            </v-card-text>
          </v-card>
        </v-dialog>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { Component, Vue } from "nuxt-property-decorator";
import { namespace } from "vuex-class";
import Account from "~/model/Account";
const snackbar = namespace("Snackbar");
const spinner = namespace("Spinner");
import Help from "~/components/Help.vue";
import ICountry from "~/model/ICountry";

@Component({
  components: { Help },
})
export default class ProfilePage extends Vue {
  @snackbar.Action
  public updateText!: (newText: string) => void;

  @snackbar.Action
  public updateColor!: (newColor: string) => void;

  @snackbar.Action
  public updateShow!: (newShow: boolean) => void;

  @spinner.Action
  public toggleSpinner!: (newShow: boolean) => void;

  $i18n: any;
  // @ts-ignore
  title: string = this.$i18n.t("profile.title");
  valid: boolean = true;
  nameRules: Array<any> = [
    (v: any) => !!v || this.$i18n.t("profile.nameis"),
    (v: string | any[]) =>
      (v && v.length <= 50) || this.$i18n.t("profile.namemust"),
  ];
  emailRules: Array<any> = [
    (v: any) => !!v || this.$i18n.t("profile.emailis"),
    (v: string) => /.+@.+\..+/.test(v) || this.$i18n.t("profile.emailmust"),
  ];
  streetRules: Array<any> = [
    (v: any) => !!v || this.$i18n.t("profile.streetis"),
    (v: string | any[]) =>
      (v && v.length <= 50) || this.$i18n.t("profile.streetmust"),
  ];
  cityRules: Array<any> = [
    (v: any) => !!v || this.$i18n.t("profile.cityis"),
    (v: string | any[]) =>
      (v && v.length <= 40) || this.$i18n.t("profile.citymust"),
  ];
  zipRules: Array<any> = [
    (v: any) => !!v || this.$i18n.t("profile.zipis"),
    (v: string | any[]) =>
      (v && v.length <= 6) || this.$i18n.t("profile.zipmust"),
  ];
  companyNumberRules: Array<any> = [
    (v: any) => !!v || this.$i18n.t("profile.companyis"),
    (v: string | any[]) =>
      (v && v.length <= 15) || this.$i18n.t("profile.companymust"),
  ];
  companyVatRules: Array<any> = [
    (v: any) => !!v || this.$i18n.t("profile.companyis1"),
    (v: string | any[]) =>
      (v && v.length <= 15) || this.$i18n.t("profile.companymust"),
  ];
  countryRules: Array<any> = [
    (v: any) => !!v || this.$i18n.t("profile.countrymust"),
  ];
  passwordRules: Array<any> = [];
  accountId: string = "";
  name: string = "";
  email: string = "";
  street: string = "";
  city: string = "";
  zip: string = "";
  countryCode: string = "";
  companyNumber: string = "";
  vatNumber: string = "";
  paidTo: string = "";
  planRefer: string = "";
  password: string = "";
  newsletter: boolean = false;
  countries: Array<ICountry> = [];
  $axios: any;
  plans: any;
  $refs: any;
  dialogClose: boolean = false;
  $nuxt: any;
  accountConfirmName: string = "";
  $t: any;

  head() {
    return {
      title: this.title,
    };
  }

  mounted() {
    this.getProfile();
    this.fillCountries();
  }

  getProfile() {
    this.$axios
      .get(this.$config.internalApi + "plans")
      .then((response: any) => {
        if (response.data.success === true) {
          this.plans = response.data.plans;
          this.$axios
            .get("/accounts/" + this.$auth.$state.user.Id)
            .then((response: any) => {
              if (response.data.success === true) {
                this.name = response.data.account.Name;
                this.email = response.data.account.Email;
                this.street = response.data.account.Street;
                this.city = response.data.account.City;
                this.zip = response.data.account.Zip;
                this.countryCode = response.data.account.CountryCode.toUpperCase();
                this.companyNumber = response.data.account.CompanyNumber;
                this.vatNumber = response.data.account.VatNumber;
                this.newsletter = response.data.account.Newsletter;
                let d = new Date(response.data.account.PaidTo);
                this.paidTo = d.toLocaleDateString("en-US");
                const plan = this.plans.find(
                  (item: any) => item.Id === response.data.account.PlanRefer
                );
                this.planRefer = plan !== undefined ? plan.Name : "";
                this.accountId = response.data.account.Id;
              } else {
                this.updateText(response.data.error);
                this.updateColor("red");
                this.updateShow(true);
              }
            });
        } else {
          this.updateText(response.data.error);
          this.updateColor("red");
          this.updateShow(true);
        }
      });
  }

  customFilter(item: ICountry, queryText: string, itemText: string) {
    const textOne = item.Name.toLowerCase();
    const textTwo = item.Code.toLowerCase();
    const searchText = queryText.toLowerCase();

    return textOne.indexOf(searchText) > -1 || textTwo.indexOf(searchText) > -1;
  }

  fillCountries() {
    this.$axios
      .get(
        window.location.protocol +
          "//" +
          window.location.host +
          "/countries.json"
      )
      .then((res: any) => (this.countries = res.data));
  }

  saveProfile() {
    this.$refs.form.validate();
    this.toggleSpinner(true);
    const userObj: any = {
      ID: this.accountId,
      name: this.name,
      email: this.email,
      street: this.street,
      city: this.city,
      zip: this.zip,
      countryCode: this.countryCode,
      companyNumber: this.companyNumber,
      vatNumber: this.vatNumber,
      newsletter: this.newsletter,
    };

    if (this.password.length > 6) {
      userObj.password = this.password;
    }

    this.$axios
      .put(this.$config.internalApi + "accounts", userObj)
      .then((response: any) => {
        if (response.data.success === true) {
          this.updateText(this.$i18n.t("profile.update"));
          this.updateColor("green");
          this.updateShow(true);
          this.toggleSpinner(false);
        } else {
          this.updateText(response.data.error);
          this.updateColor("red");
          this.updateShow(true);
          this.toggleSpinner(false);
        }
      });
  }

  closeDialogClose() {
    this.dialogClose = false;
  }

  closeAccountConfirm() {
    this.$axios
      .delete(
        this.$config.internalApi + "accounts/" + this.$auth.$state.user.Id
      )
      .then((response: any) => {
        if (response.data.success === true) {
          this.$nuxt.$options.router.push("/sign");
        } else {
          this.updateText(response.data.error);
          this.updateColor("red");
          this.updateShow(true);
        }
      });
  }
}
</script>
<style>
.v-icon {
  color: #000059;
}

.theme--dark.v-icon {
  color: #e1863b !important;
}

.theme--dark.v-application .primary {
  border-color: #e1863b !important;
  background-color: #e1863b !important;
}
.theme--light.v-label,
.theme--light.v-select .v-select__selections,
.theme--light.v-input input,
.theme--light.v-input textarea {
  color: #000 !important;
}
</style>
