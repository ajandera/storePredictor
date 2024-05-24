<template>
  <v-app :dark="false">
    <v-row>
      <v-col cols="12" sm="12" md="12" lq="12" xl="12" class="header">
        <v-container>
          <Nuxt />
        </v-container>
      </v-col>
    </v-row>
    <Snackbar />
    <Spinner />
  </v-app>
</template>

<script lang="ts">
import { Component, Vue } from "nuxt-property-decorator";
import Spinner from "~/components/Spinner.vue";
import Snackbar from "~/components/Snackbar.vue";
import IDictionary from "~/model/IDictionary";
import IResponseText from "~/model/IResponseText";

@Component({
  components: {
    Spinner,
    Snackbar,
  }
})
export default class LoginLayout extends Vue {
  fixed: boolean = true;
  miniVariant: boolean = false;
  right: boolean = true;
  title: string = "storePredictor";
  $i18n: any;
  $route: any;
  $nuxt: any;
}
</script>
<style scoped>
.v-select .v-input__slot {
  margin: 0 !important;
}

.header {
  position: relative;
  text-align: center;
  background: linear-gradient(
    60deg,
    rgba(84, 58, 183, 1) 0%,
    rgba(0, 172, 193, 1) 100%
  );
  color: white;
  padding: 0;
}

.waves {
  position: absolute;
  width: 100%;
  margin-bottom: -7px; /*Fix for safari gap*/
  left: 0;
  bottom: 0;
  height: 40%;
}

/* Animation */

.parallax > use {
  animation: move-forever 25s cubic-bezier(0.55, 0.5, 0.45, 0.5) infinite;
}

.parallax > use:nth-child(1) {
  animation-delay: -2s;
  animation-duration: 7s;
}

.parallax > use:nth-child(2) {
  animation-delay: -3s;
  animation-duration: 10s;
}

.parallax > use:nth-child(3) {
  animation-delay: -4s;
  animation-duration: 13s;
}

.parallax > use:nth-child(4) {
  animation-delay: -5s;
  animation-duration: 20s;
}

@keyframes move-forever {
  0% {
    transform: translate3d(-90px, 0, 0);
  }
  100% {
    transform: translate3d(85px, 0, 0);
  }
}

.infoText {
  padding: 50px 100px;
}
.img-fluid {
  max-width: 100%;
}
.wordpress {
  filter: invert(1);
}
</style>
