<template>
  <v-app>
    <Nuxt />
  </v-app>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "nuxt-property-decorator";

@Component
export default class EmptyLayout extends Vue {
  @Prop() readonly error!: any;
  @Prop() readonly default!: string;
}
</script>

<style scoped>

</style>
