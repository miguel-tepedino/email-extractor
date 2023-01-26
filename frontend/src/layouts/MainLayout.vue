<script lang="ts" setup>
import HeaderMain from "@/components/header/HeaderMain.vue";
</script>

<template>
  <div
    class="bg-gradient-to-tr from-cyan-500 to-green-600 w-full py-5 flex flex-col"
  >
    <HeaderMain />

    <slot />
    <footer></footer>
  </div>
</template>
