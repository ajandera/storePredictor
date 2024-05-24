<template>
  <div id="header" class="header errorHeader">
    <h1 v-if="error.statusCode === 404">
      <img src="/sp_logo.svg" alt="storePredictor" class="logo" /><br>
      {{ pageNotFound }}<br><br>
      <NuxtLink to="/">
        <v-btn color="danger">Home page</v-btn>
      </NuxtLink>
    </h1>
    <h1 v-else>
      <img src="/sp_logo.svg" alt="storePredictor" class="logo" /><br>
      {{ otherError }}<br><br>
      <NuxtLink to="/">
        <v-btn color="danger">Home page</v-btn>
      </NuxtLink>
    </h1>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "nuxt-property-decorator";

@Component({
  layout: 'empty'
})
export default class ErrorLayout extends Vue {
  @Prop() readonly type!: object;
  @Prop() readonly default!: string;
  @Prop() readonly error!: any;

  pageNotFound: string = "404 Not Found";
  otherError: string = "An error occurred";

  head() {
    return {
      title: this.error.statusCode === 404 ? this.pageNotFound : this.otherError
    };
  }
}
</script>

<style scoped>
h1 {
  font-size: 20px;
  text-align: center;
  padding: 20%;
}
.errorHeader {
  height: 100%;
  display: block;
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
.back {
  text-align: center;
}
.logo{
  filter: brightness(0) invert(1);
}
#__nuxt,
#__layout {
  height: 100%;
}
a {
  text-decoration: none;
}
</style>
