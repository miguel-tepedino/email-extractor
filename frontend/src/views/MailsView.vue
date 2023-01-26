<script lang="ts" setup>
import ModalComponent from "@/components/modal/ModalComponent.vue";
import { ref } from "vue";
import EmailComponent from "@/components/emails/EmailComponent.vue";
import type { Email } from "@/components/emails/type";
import mailsStore from "@/stores/mails";
import { storeToRefs } from "pinia";

const mailStore = mailsStore();

const { returnMails, mailError, isFetchingMails, offset } =
  storeToRefs(mailStore);

const openModal = ref<boolean>(false);

const emailToShow = ref<Email | null>(null);

mailStore.getMails();

function getMoreEmails() {
  mailStore.getOffsetEmails();
}

function backOffsetEmails() {
  mailStore.backOffsetEmails();
}
</script>

<template>
  <main class="mt-16">
    <div v-if="!!mailError">
      {{ mailError }}
    </div>
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
              v-for="email in returnMails"
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
