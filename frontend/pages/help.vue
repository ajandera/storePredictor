<template>
  <v-container
    class="col-xl-6 col-lg-6 col-md-8 col-sm-6 col-10 offset-1 offset-lg-3 offset-xl-3 offset-md-2 offset-sm-3"
  >
    <div class="wrapper">
      <img class="logo-image d-block d-sm-none" src="/sp_logo.svg" alt="alternative" />
      <div class="slide-controls">
        <input type="radio" name="slide" id="Login" />
        <input type="radio" name="slide" id="Signup" />
        <label for="Login" class="slide" @click="activeSignin">{{
          $t("signin.login")
        }}</label>
        <label for="Signup" class="slide Signup" @click="activeRegister">{{
          $t("signin.register")
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
      <v-form
        class="mt-2 text-center"
        @submit="register($event)"
        v-if="registerFlag === true"
      >
        <div class="form-group mt-4">
          <input
            class="uiInput"
            v-bind:placeholder="$t('signin.mail')"
            type="email"
            id="email1"
            v-model="email"
            required
          />
        </div>

        <div class="gorm-group mt-4">
          <input
            class="uiInput"
            type="password"
            v-bind:placeholder="$t('signin.password')"
            id="password1"
            v-model="passwordReg"
            required
          />
        </div>

        <div class="form-group mt-4">
          <input
            class="uiInput color-text"
            type="password"
            id="password2"
            v-model="passwordCheck"
            v-bind:placeholder="$t('signin.password2')"
            required
          />
        </div>
        <div class="form-group mt-5 groupCenter">
          <v-checkbox
            v-model="newsletter"
            :label="$t('signin.newsletter')"
            color="success"
            class="white--text"
          ></v-checkbox>
        </div>
        <p
          style="color: red"
          v-if="
            passwordReg !== passwordCheck &&
            passwordCheck !== null &&
            passwordCheck.length > 1
          "
        >
          {{ $t("signin.mismatch") }}
        </p>
        <input
          :disabled="passwordReg !== passwordCheck"
          type="submit"
          class="butt"
          v-bind:value="$t('signin.register1')"
        />
      </v-form>
      <v-btn
        block
        @click="socialLogin()"
        class="google mt-5">{{ $t('signin.social')}}</v-btn>
      <p>Google, Facebook, linkedin</p>
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
  </v-container>
</template>

<script lang="ts">
import { Component, Vue } from "nuxt-property-decorator";
import { namespace } from "vuex-class";

const snackbar = namespace("Snackbar");
const spinner = namespace("Spinner");

@Component
export default class HelpPage extends Vue {
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

  // @ts-ignore
  title: string = this.$i18n.t("help.title");
  languages: Array<string> = ["en", "cs", "de", "sk", "pl", "ro", "hu"];

  head() {
    return {
      title: this.title,
    };
  }

  onLangChange(event: string) {
    this.$i18n.locale = event;
  }

}
</script>

<style scoped>
</style>
