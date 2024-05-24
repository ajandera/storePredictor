<template>
  <v-container
    class="col-lg-4 col-md-6 col-sm-6 offset-lg-4 offset-md-3 offset-sm-2"
  >
    <div class="wrapper">
      <img class="logo-image" src="/sp_logo.svg" alt="alternative" />
      <v-form class="sign text-center" @submit="restore($event)">
        <p>{{ $t("signin.setNewPw") }}</p>
        <div class="gorm-group">
          <input
            class="uiInput"
            :placeholder="$t('signin.password')"
            type="password"
            id="password"
            v-model="password"
            required
          />
        </div>
        <div class="form-group mt-4">
          <input
            class="uiInput"
            type="password"
            id="passwordCheck"
            v-model="passwordCheck"
            :placeholder="$t('signin.password2')"
            required
          />
        </div>
        <input
          type="submit"
          class="butt mt-4"
          :value="$t('signin.restoreSend')"
        />
      </v-form>
      <br />
      <NuxtLink to="sing">Back to login page.</NuxtLink>
    </div>
  </v-container>
</template>

<script lang="ts">
import { Component, Vue } from "nuxt-property-decorator";
import Forgot from "~/components/Forgot.vue";
import { namespace } from "vuex-class";

const snackbar = namespace("Snackbar");
const spinner = namespace("Spinner");

@Component({
  components: {
    Forgot,
  },
  layout: "login",
})
export default class RestorePage extends Vue {
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
  title: string = this.$i18n.t("restore.title");
  token: any;
  password!: string;
  passwordCheck!: string;
  $axios: any;

  head() {
    return {
      title: this.title,
    };
  }

  restore(e: Event) {
    e.preventDefault();
    this.token = this.$route.query.token;
    this.toggleSpinner(true);
    this.$axios
      .put(this.$config.internalApi + "restore/" + this.token, {
        password: this.password,
      })
      .then((response: any) => {
        if (response.data.success === true) {
          this.updateText(this.$i18n.t("restore.success"));
          this.updateColor("green");
          this.updateShow(true);
          this.toggleSpinner(false);
        } else {
          this.updateText(response.data.message);
          this.updateColor("red");
          this.updateShow(true);
          this.toggleSpinner(false);
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
.container {
  position: relative;
  padding-top: 5%;
  z-index: 9;
}
.wrapper {
  position: relative;
  background: #4633af;
  border-radius: 25px;
  height: 100%;
  padding: 10%;
  color: #fff;
  min-height: 630px;
}

.wrapper .uiInput {
  background: linear-gradient(45deg, #fff, #fff);
  width: 230px;
  border-radius: 25px;
  padding: 10px 10px 10px 20px;
  border: none;
  -webkit-box-shadow: 0px 10px 49px -14px rgba(0, 0, 0, 0.7);
  -moz-box-shadow: 0px 10px 49px -14px rgba(0, 0, 0, 0.7);
  box-shadow: 0px 10px 49px -14px rgba(0, 0, 0, 0.7);
}

.wrapper .uiInput:focus {
  outline: none;
}

.sign {
  position: relative;
  margin-top: 120px;
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

.wrapper .butt:focus {
  outline: none;
}

.logo-image {
  filter: brightness(0) invert(1);
  transform: translate(-1%, -40%);
  width: 60%;
  margin: 0 auto;
  position: absolute;
  left: 20%;
  top: 10%;
}

.click {
  cursor: pointer;
}

input[type="radio"] {
  display: none;
}

/* Extra small devices (phones, 600px and down) */
@media only screen and (max-width: 600px) {
  .wrapper {
    min-height: 325px;
  }
}
</style>
