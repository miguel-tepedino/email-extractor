<script lang="ts" setup>
import authstore from "../../stores/auth";
import { useRouter } from "vue-router";

const router = useRouter();

const authst = authstore();

function logout() {
  authst.logout();
  router.push("/");
}
</script>

<template>
  <header
    class="fixed flex flex-row justify-between top-0 left-0 right-0 px-10 py-4 bg-white"
  >
    <div>Email Extractor</div>
    <div class="flex flex-row gap-4">
      <a role="button" @click="logout">Logout</a>
    </div>
  </header>
</template>
