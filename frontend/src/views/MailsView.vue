<script lang="ts" setup>
import { ref } from "vue";
import mailsStore from "@/stores/mails";
import { storeToRefs } from "pinia";
import type { MoreLess } from "@/stores/types";

import ModalComponent from "@/components/modal/ModalComponent.vue";
import EmailComponent from "@/components/emails/EmailComponent.vue";
import type { Email } from "@/components/emails/type";

const openModal = ref<boolean>(false);

const emailToShow = ref<Email | null>(null);

const store = mailsStore();

const {
  returnMails: mails,
  mailError,
  isFetchingMails,
  offset,
} = storeToRefs(store);

store.getMails();

function pagination(r: MoreLess) {
  store.loadMoreLess(r);
}
</script>

<template>
  <main class="mt-16">
    <div v-if="mailError">{{ mailError }}</div>
    <div v-else class="px-5 flex flex-col">
      <div class="max-w-full overflow-auto">
        <table
          v-if="!isFetchingMails"
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
              v-for="email in mails"
              :key="email._id"
            >
              <td class="text-center">{{ email._source.Subject }}</td>
              <td class="text-center">{{ email._source.From }}</td>
              <td class="text-center">{{ email._source.To }}</td>
            </tr>
          </tbody>
        </table>
        <div class="text-center bg-white py-2 mt-5" v-if="mails?.length === 0">
          <span class="text-2xl">No emails to show</span>
        </div>
      </div>
      <div v-if="mails?.length !== 0" class="mt-2 flex flex-row gap-5 self-end">
        <button @click="pagination('LESS')" v-if="offset > 0" type="button">
          &lt; Back
        </button>
        <button @click="pagination('MORE')" type="button">Next ></button>
        <ModalComponent v-model="openModal">
          <EmailComponent v-if="emailToShow != null" :email="emailToShow" />
        </ModalComponent>
      </div>
    </div>
  </main>
</template>

<style scoped></style>
