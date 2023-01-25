<script lang="ts" setup>
import HeaderMain from "@/components/header/HeaderMain.vue";
</script>

<template>
  <div class="bg-blue-500 h-full py-5 flex flex-col">
    <HeaderMain />

    <slot />
    <footer></footer>
  </div>
</template>
