<template>
  <img src="/back.png" />
</template>

<script lang="ts">
import { Component, Vue } from "nuxt-property-decorator";
import { namespace } from "vuex-class";

@Component({
  layout: 'empty'
})
export default class StoretestPage extends Vue {

  head() {
    return {
      title: "My Store",
    };
  }
}
</script>

<style scoped>
</style>
