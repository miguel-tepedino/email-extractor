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


  return {

  };
});
