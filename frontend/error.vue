<template>
  <div id="header" class="header errorHeader">
    <h1 v-if="error && error.statusCode === 404">
      <img src="/sp_logo.svg" alt="storePredictor" class="logo" /><br />
      {{ pageNotFound }}<br /><br />
      <NuxtLink to="/">
        <v-btn color="error">Home page</v-btn>
      </NuxtLink>
    </h1>
    <h1 v-else>
      <img src="/sp_logo.svg" alt="storePredictor" class="logo" /><br />
      {{ otherError }}<br /><br />
      <NuxtLink to="/">
        <v-btn color="error">Home page</v-btn>
      </NuxtLink>
    </h1>
  </div>
</template>

<script setup lang="ts">
import { useHead } from "#imports";

const props = defineProps<{
  error?: { statusCode?: number; message?: string };
}>();

const pageNotFound = "404 Not Found";
const otherError = "An error occurred";

useHead({
  title: props.error?.statusCode === 404 ? pageNotFound : otherError,
});
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
.logo {
  filter: brightness(0) invert(1);
}
a {
  text-decoration: none;
}
</style>
