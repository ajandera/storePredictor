<template></template>

<script lang="ts">
import { Component, Vue } from "nuxt-property-decorator";


@Component
export default class CallbackPage extends Vue {

  $auth: any;

  created() {
    console.log(this.$auth.$storage.getUniversal('user'));
    debugger;
  }
}
</script>

<style scoped>

</style>
