<script setup lang="ts">
import { computed } from "vue";
import { useRoute } from "vue-router";

const route = useRoute();

const layoutHandler = computed<string>(() => {
  return (route.meta.layout as string) || "DefaultLayout";
});
</script>

<template>
  <component :is="layoutHandler">
    <slot />
  </component>
</template>

<style></style>
