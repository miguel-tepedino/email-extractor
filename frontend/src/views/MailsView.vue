<script lang="ts" setup>
import ModalComponent from "@/components/modal/ModalComponent.vue";
import { ref } from "vue";
import EmailComponent from "@/components/emails/EmailComponent.vue";
import type { Email } from "@/components/emails/type";
import type { EmailResponse } from "@/stores/types";
import { useFetch } from "@vueuse/core";

const openModal = ref<boolean>(false);

const emailToShow = ref<Email | null>(null);

const offset = ref<number>(0);

const { data, error, isFetching } = useFetch<EmailResponse>(
  "http://localhost:3000/getmails",
  { immediate: true }
).json<EmailResponse>();

console.log(data);

async function getMoreEmails() {
  offset.value += 10;
  const { data: newData, error: newError } = await useFetch<EmailResponse>(
    `http://localhost:3000/mails/${offset.value}`
  ).json<EmailResponse>();
  if (newError.value) {
    error.value = newError.value;
    return;
  }
  data.value = newData.value;
}

async function backOffsetEmails() {
  offset.value -= 10;
  const { data: newData, error: newError } = await useFetch<EmailResponse>(
    `http://localhost:3000/mails/${offset.value}`
  ).json<EmailResponse>();
  if (newError.value) {
    error.value = newError.value;
    return;
  }
  data.value = newData.value;
}
</script>

<template>
  <main class="mt-16">
    <div v-if="error">Error</div>
    <div v-else class="px-5 flex flex-col">
      <div class="max-w-full overflow-auto">
        <table
          v-if="!isFetching"
          class="w-full table-auto rounded-lg border-none bg-white"
        >
          <thead class="border-none">
            <tr class="border-none py-5">
              <th class="border-slate-300">Subject:</th>
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
      </div>
      <div class="mt-2 flex flex-row gap-5 self-end">
        <button v-if="offset > 0" @click="backOffsetEmails" type="button">
          &lt; Back
        </button>
        <button @click="getMoreEmails" type="button">Next ></button>
        <ModalComponent v-model="openModal">
          <EmailComponent v-if="emailToShow != null" :email="emailToShow" />
        </ModalComponent>
      </div>
    </div>
  </main>
</template>

<style scoped></style>
