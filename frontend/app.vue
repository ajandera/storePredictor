<template>
  <NuxtLayout :name="layoutName" />
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useRoute } from "#imports";

const route = useRoute();

const layoutName = computed(() => {
  if (route.meta?.layout) return route.meta.layout as string;
  const cleanPath = route.path.replace(/^\/(en|cz|cs|sk|de|pl|ro|hu)(\/|$)/, "/") || "/";
  if (["/sign", "/signup", "/restore", "/token"].includes(cleanPath)) {
    return "login";
  }
  if (cleanPath === "/storetest") {
    return "empty";
  }
  return "default";
});
</script>
