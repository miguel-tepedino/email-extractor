<script lang="ts" setup>
import { useFetch } from "@vueuse/core";
import ModalComponent from "@/components/modal/ModalComponent.vue";
import { ref } from "vue";
import EmailComponent from "@/components/emails/EmailComponent.vue";
import type { Email } from "@/components/emails/type";

const openModal = ref<boolean>(false);

const emailToShow = ref<Email | null>(null);

const { data, error, isFetching } = useFetch<{ hits: { hits: any[] } }>(
  "http://localhost:3000/getmails"
).json();
</script>

<template>
  <div class="mt-16">
    <div v-if="error">
      {{ error }}
    </div>
    <div class="px-5 flex flex-col" v-if="!isFetching">
      <table class="w-full table-auto rounded-lg border-none bg-white">
        <thead class="border-none">
          <tr class="border-none">
            <th class="py-3 border-slate-300">Subject:</th>
            <th class="border-slate-300">From:</th>
            <th class="border-slate-300">To:</th>
          </tr>
        </thead>
        <tbody>
          <tr
            @click="
              emailToShow = email._source;
              openModal = true;
            "
            class="overflow-hidden border-t-2"
            v-for="email in data?.hits.hits"
            :key="email._id"
          >
            <td class="text-center">{{ email._source.Subject }}</td>
            <td class="text-center">{{ email._source.From }}</td>
            <td class="text-center">{{ email._source.To }}</td>
          </tr>
        </tbody>
      </table>
      <div class="mt-2 self-end">
        <button type="button">Next ></button>
        <ModalComponent v-model="openModal">
          <EmailComponent v-if="emailToShow != null" :email="emailToShow" />
        </ModalComponent>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
