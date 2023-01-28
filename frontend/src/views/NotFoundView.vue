<script lang="ts" setup>
import { useRouter } from "vue-router";

const router = useRouter();

function redirect() {
  router.push("/");
}
</script>

<template>
  <main class="grid grid-flow-row gap-10 place-content-center h-full w-full">
    <h1 class="text-5xl text-white">Sorry, page not found</h1>
    <div class="flex justify-center">
      <button
        type="button"
        @click="redirect"
        class="bg-white px-4 py-2 rounded-full block"
      >
        Return
      </button>
    </div>
  </main>
</template>
