<template>
  <v-container
    class="col-xl-6 col-lg-6 col-md-8 col-sm-6 col-10 offset-1 offset-lg-3 offset-xl-3 offset-md-2 offset-sm-3"
  >
    <div class="wrapper">
      <img class="logo-image" src="/sp_logo.svg" alt="alternative" />
      <br>
      <div class="slide-controls">
        <label for="Login" class="slide" @click="activeSignin">{{
          $t("signin.login")
        }}</label>
        <div ref="active" class="slide-tab"></div>
      </div>
      <v-form
        class="mt-10 text-center"
        @submit="signin($event)"
        v-if="signinFlag === true"
      >
        <div class="form-group">
          <input
            class="uiInput"
            v-bind:placeholder="$t('signin.mail')"
            type="email"
            id="email"
            v-model="username"
            required
          />
        </div>
        <div class="form-group mt-4">
          <input
            class="uiInput"
            type="password"
            id="password"
            v-model="password"
            v-bind:placeholder="$t('signin.password')"
            required
          />
        </div>
        <input
          type="submit"
          class="butt mt-4"
          v-bind:value="$t('signin.signin1')"
        />
        <div class="text-center">
          <div class="mt-4 click" v-on:click="showForgot = !showForgot">
            {{ $t("signin.forgot") }}
          </div>
          <p class="forgot-password" v-if="showForgot === true">
            <Forgot />
          </p>
        </div>
      </v-form>
      <NuxtLink to="/signup" class="btn btn-primary mt-5">{{ $t('signin.register')}} </NuxtLink>
      <v-col cols="12" align="right" class="mt-5">
        <v-menu bottom origin="center center" transition="scale-transition">
          <template v-slot:activator="{ on, attrs }">
            <v-btn color="secondary" dark v-bind="attrs" v-on="on">
              <country-flag
                :country="$i18n.locale.replace('en', 'us').replace('cs', 'cz')"
                size="small"
              />
              {{ $i18n.locale }}
            </v-btn>
          </template>

          <v-list>
            <v-list-item
              v-for="(item, i) in languages"
              :key="i"
              @click="onLangChange(item)"
            >
              <v-list-item-title
                ><country-flag
                  :country="item.replace('en', 'us').replace('cs', 'cz')"
                  size="small"
                />
                {{ item }}</v-list-item-title
              >
            </v-list-item>
          </v-list>
        </v-menu>
      </v-col>
    </div>
    <v-row class="mt-5 align-items-center">
          <v-col cols="2" class="d-flex">
            <img
              class="img-fluid align-self-center"
              src="https://storepredictor.com/integration/prestashop.png"
              alt="Prestashop"
            />
          </v-col>
          <v-col cols="2" class="d-flex">
            <img
              class="img-fluid wordpress align-self-center"
              src="https://storepredictor.com/integration/wordpress.png"
              alt="Wordpress"
            />
          </v-col>
          <v-col cols="2" class="d-flex">
            <img
              class="img-fluid magento align-self-center"
              src="https://storepredictor.com/integration/magento.png"
              alt="Magento"
            />
          </v-col>
          <v-col cols="2" class="d-flex">
            <img
              class="img-fluid align-self-center"
              src="https://storepredictor.com/integration/shopify.png"
              alt="Shopify"
            />
          </v-col>
          <v-col cols="2" class="d-flex">
            <img
              class="img-fluid align-self-center"
              src="https://storepredictor.com/integration/opencart.png"
              alt="Opencart"
            />
          </v-col>
          <v-col cols="2" class="d-flex">
            <img
              class="img-fluid gtm align-self-center"
              src="https://storepredictor.com/integration/gtm.png"
              alt="Google Tag Manager"
            />
          </v-col>
        </v-row>
  </v-container>
</template>

<script lang="ts">
import { Component, Vue } from "nuxt-property-decorator";
import Forgot from "~/components/Forgot.vue";
import { namespace } from "vuex-class";
import MattermostService from "~/services/MattermostService"

const snackbar = namespace("Snackbar");
const spinner = namespace("Spinner");

@Component({
  components: {
    Forgot
  },
  layout: "login",
})
export default class SignPage extends Vue {
  @snackbar.Action
  public updateText!: (newText: string) => void;

  @snackbar.Action
  public updateColor!: (newColor: string) => void;

  @snackbar.Action
  public updateShow!: (newShow: boolean) => void;

  @spinner.Action
  public toggleSpinner!: (newShow: boolean) => void;

  $t: any;
  $i18n: any;
  $nuxt: any;
  $axios: any;
  $refs: any;
  $auth: any;
  newsletter: boolean = false;
  // @ts-ignore
  title: string = this.$i18n.t("signin.title");
  username: string = "";
  password: string = "";
  showForgot: boolean = false;
  email: string = "";
  passwordReg: string = "";
  passwordCheck: string = "";
  registerFlag: boolean = false;
  signinFlag: boolean = true;
  languages: Array<string> = ["en", "cs", "de", "sk", "pl", "ro", "hu"];
  emailForDemo: string = "";

  head() {
    return {
      title: this.title,
    };
  }

  async signin(e: Event) {
    e.preventDefault();
    await this.loginUser();
  }

  async loginUser() {
    this.toggleSpinner(true);
      try {
        let login = {
          email: this.username,
          password: this.password,
        };
        await this.$auth.loginWith("local", { data: login });
        await this.$nuxt.$options.router.push("/");
        this.toggleSpinner(false);
      } catch (err) {
        this.toggleSpinner(false);
        this.updateText(this.$i18n.t("signin.wrong"));
        this.updateColor("red");
        this.updateShow(true);
      }
  }

  activeSignin() {
    this.$refs.active.style.left = "0";
    this.signinFlag = true;
    this.registerFlag = false;
  }

  activeRegister() {
    this.$refs.active.style.left = "50%";
    this.signinFlag = false;
    this.registerFlag = true;
  }

  onLangChange(event: string) {
    this.$i18n.locale = event;
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
  z-index: 9;
}
.wrapper {
  position: relative;
  background: #4633af;
  border-radius: 25px;
  padding: 10%;
  color: #fff;
  margin-top: 20%;
}

.wrapper .uiInput {
  background: #fff;
  border-radius: 25px;
  padding: 10px 10px 10px 20px;
  border: none;
  width: 100%;
  -webkit-box-shadow: 0px 10px 49px -14px rgba(0, 0, 0, 0.7);
  -moz-box-shadow: 0px 10px 49px -14px rgba(0, 0, 0, 0.7);
  box-shadow: 0px 10px 49px -14px rgba(0, 0, 0, 0.7);
}

.wrapper .uiInput:focus {
  outline: none;
}

.butt {
  background: linear-gradient(45deg, #b1ea4e, #35b219);
  color: #fff;
  width: 100%;
  border: none;
  border-radius: 25px;
  padding: 10px;
  -webkit-box-shadow: 0px 10px 41px -11px rgba(0, 0, 0, 0.7);
  -moz-box-shadow: 0px 10px 41px -11px rgba(0, 0, 0, 0.7);
  box-shadow: 0px 10px 41px -11px rgba(0, 0, 0, 0.7);
}
.butt.email {
  width: 50%;
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
  position: absolute;
  left: 20%;
}

.click {
  cursor: pointer;
}

.wrapper .slide-controls {
  position: relative;
  display: flex;
  height: 50px;
  width: 100%;
  overflow: hidden;
  margin: 10% 0 10px 0;
  justify-content: space-between;
  border-radius: 25px;
}

.slide-controls .slide {
  height: 100%;
  width: 100%;
  color: #fff;
  font-size: 18px;
  font-weight: 500;
  text-align: center;
  line-height: 48px;
  cursor: pointer;
  z-index: 1;
  transition: all 0.6s ease;
}

.slide-controls label.signup {
  color: #000;
}

input[type="radio"] {
  display: none;
}
.slide-controls {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 25px;
}
.white--text /deep/ label {
  color: #fff;
}
.white--text /deep/ .v-icon {
  color: #fff;
}
.groupCenter {
  width: 40%;
  margin: 0 auto;
}
.wrapper[data-v-27c68594] {
  bottom: -51px;

  justify-content: center;
}
.v-list-item__title {
  padding-left: 10px !important;
}
/* Extra small devices (phones, 600px and down) */
@media only screen and (max-width: 600px) {
  .wrapper {
    bottom: 0 !important;
  }
  .wrapper .slide-controls {
    margin: 20% 0 10px 0;
  }
}

.google {
  background-color: #fff;
}
.facebook {
  background-color: #3b5998 !important;
}
.linkedin {
  background-color: #0077b5 !important;
}
</style>
