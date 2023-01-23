<script lang="ts" setup>
import HeaderMain from "@/components/header/HeaderMain.vue";
</script>

<template>
  <div class="h-screen w-screen bg-blue-500">
    <header>
      <HeaderMain />
    </header>
    <slot />
    <footer></footer>
  </div>
</template>
