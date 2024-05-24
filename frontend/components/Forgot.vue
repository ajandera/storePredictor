<template>
  <v-form @submit.prevent="handleSubmit">
    <div class="form-group mt-4">
      <input
        type="email"
        class="myInput"
        v-model="email"
        :placeholder="$t('signin.forgotEmail')"
      />
    </div>
    <button type="submit" class="butt mt-4">{{ $t("send") }}</button>
  </v-form>
</template>

<script lang="ts">
import { Component, Vue } from "nuxt-property-decorator";
import { namespace } from "vuex-class";
const snackbar = namespace("Snackbar");
const spinner = namespace("Spinner");

@Component
export default class Forgot extends Vue {
  @snackbar.Action
  public updateText!: (newText: string) => void;

  @snackbar.Action
  public updateColor!: (newColor: string) => void;

  @snackbar.Action
  public updateShow!: (newShow: boolean) => void;

  @spinner.Action
  public toggleSpinner!: (newShow: boolean) => void;

  email: string = "";
  $axios: any;
  $i18n: any;

  handleSubmit() {
    this.toggleSpinner(true);
    this.$axios.get(this.$config.internalApi + "account/email/"+this.email).then((response: any) => {
      if (response.data.account.Id !== '') {
          this.sendNewPassword(response.data.account.Id)
      } else {
        this.updateText(this.$i18n.t('signin.notReg'));
        this.updateColor('red');
        this.updateShow(true);
        this.toggleSpinner(true);
      }
    });
  }

  sendNewPassword(id: string) {
    this.$axios.get(this.$config.internalApi + "forgot/" + id+ "/"+this.$i18n.locale).then((response: any) => {
      if (response.data.success === true) {
        this.updateText(this.$i18n.t('signin.pwSend'));
        this.updateColor('green');
      } else {
        this.updateText(response.data.message);
        this.updateColor('red');
        this.updateShow(true);
        this.toggleSpinner(true);
      }
    });
  }
}
</script>

<style scoped>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

.myInput {
  background: linear-gradient(45deg, #fff, #fff);
  width: 230px;
  border-radius: 25px;
  padding: 10px 10px 10px 20px;
  border: none;
  -webkit-box-shadow: 0px 10px 49px -14px rgba(0, 0, 0, 0.7);
  -moz-box-shadow: 0px 10px 49px -14px rgba(0, 0, 0, 0.7);
  box-shadow: 0px 10px 49px -14px rgba(0, 0, 0, 0.7);
}
.butt {
  background: linear-gradient(45deg, #b1ea4e, #35b219);
  color: #fff;
  width: 230px;
  border: none;
  border-radius: 25px;
  padding: 10px;
  -webkit-box-shadow: 0px 10px 41px -11px rgba(0, 0, 0, 0.7);
  -moz-box-shadow: 0px 10px 41px -11px rgba(0, 0, 0, 0.7);
  box-shadow: 0px 10px 41px -11px rgba(0, 0, 0, 0.7);
}
.butt:hover {
  background: linear-gradient(45deg, #04c9db, #02444b);
}
</style>
