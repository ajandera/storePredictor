<template></template>

<script lang="ts">
import { Component, Vue } from "nuxt-property-decorator";

@Component({
  layout: "login",
})
export default class TokenPage extends Vue {
  $route: any;
  $nuxt: any;

  head() {
    return {
      title: "Token login",
    };
  }

  mounted() {
    if (this.$route.query.token !== undefined) {
      this.$auth.setUserToken(
        this.$route.query.token,
        this.$route.query.refresh
      );
      this.$auth
        .fetchUser()
        .then((res) => this.$nuxt.$options.router.push("/"));
    }
  }
}
</script>
