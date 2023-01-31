<script lang="ts" setup>
import authstore from "../../stores/auth";
import { useRouter } from "vue-router";
import { Field } from "vee-validate";
import mailsStore from "@/stores/mails";

const router = useRouter();

const authst = authstore();

const store = mailsStore();

function logout() {
  authst.logout();
  router.push("/");
}

function handleChange(e: any) {
  store.setSearchWord(e.target.value);
}
</script>

<template>
  <header
    class="fixed flex flex-row justify-between top-0 left-0 right-0 px-10 py-4 bg-white shadow-xl"
  >
    <div>Enron mails</div>
    <div class="flex flex-row gap-4">
      <Field
        class="rounded-full px-3 outline-none border-[3px] focus:border-cyan-500"
        name="search"
        @input.prevent="handleChange"
        type="text"
        placeholder="Search Email"
      />
      <a role="button" @click="logout">Logout</a>
    </div>
  </header>
</template>
