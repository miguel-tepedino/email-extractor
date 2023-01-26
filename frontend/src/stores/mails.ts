import { useFetch } from "@vueuse/core";
import { defineStore } from "pinia";
import { computed, ref } from "vue";
import type { EmailResponse, MailsHits } from "./types";

export default defineStore("mailsStore", () => {
  const mails = ref<MailsHits[]>();

  const mailError = ref(null);

  const isFetchingMails = ref<boolean>(false);

  const offset = ref<number>(0);

  const isFetchingMore = ref<boolean>(false);

  const returnMails = computed(() => {
    return mails.value;
  });

  async function getMails() {
    isFetchingMails.value = true;
    const { data, error } = await useFetch<EmailResponse>(
      "http://localhost:3000/getmails"
    ).json<EmailResponse>();

    if (error.value) {
      mailError.value = error.value;
      return;
    }

    mails.value = data.value?.hits.hits;
    isFetchingMails.value = false;
  }

  async function getOffsetEmails() {
    isFetchingMore.value = true;
    offset.value += 10;
    const { data, error } = await useFetch<EmailResponse>(
      `http://localhost:3000/mails/${offset.value}`
    ).json<EmailResponse>();
    if (error.value) {
      mailError.value = error.value;
      return;
    }

    mails.value = data.value?.hits.hits;
    isFetchingMails.value = false;
  }

  async function backOffsetEmails() {
    isFetchingMore.value = true;
    offset.value -= 10;
    const { data, error } = await useFetch<EmailResponse>(
      `http://localhost:3000/mails/${offset.value}`
    ).json<EmailResponse>();
    if (error.value) {
      mailError.value = error.value;
      return;
    }

    mails.value = data.value?.hits.hits;
    isFetchingMails.value = false;
  }

  return {
    returnMails,
    getMails,
    isFetchingMails,
    mailError,
    getOffsetEmails,
    backOffsetEmails,
    offset,
  };
});
