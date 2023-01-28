import { defineStore } from "pinia";
import { computed, ref } from "vue";
import type { EmailResponse, MailsHits, MoreLess } from "./types";
import { httpRequest } from "@/http/axios";

export default defineStore("mailsStore", () => {
  const mails = ref<MailsHits[]>();

  const mailError = ref<unknown>(null);

  const isFetchingMails = ref<boolean>(false);

  const isFetchingMore = ref<boolean>(false);

  const offset = ref<number>(0);

  const returnMails = computed(() => {
    return mails.value;
  });

  async function getMails() {
    isFetchingMails.value = true;
    try {
      const response = await httpRequest<EmailResponse>({
        url: "/getmails",
        method: "GET",
        data: undefined,
      });
      mails.value = response.data.hits.hits;
    } catch (e) {
      mailError.value = e;
    }
    isFetchingMails.value = false;
  }

  async function loadMoreLess(type: MoreLess) {
    type == "LESS" ? (offset.value -= 10) : (offset.value += 10);
    try {
      const response = await httpRequest<EmailResponse>({
        url: `/mails/${offset.value}`,
        data: undefined,
        method: "GET",
      });
      mails.value = response.data.hits.hits;
    } catch (e) {
      mailError.value = e;
    }
  }

  async function searchEmail(search: string) {
    try {
      const response = await httpRequest<EmailResponse>({
        url: "/search",
        method: "GET",
        data: {
          term: search,
        },
      });
      mails.value = response.data.hits.hits;
    } catch (e) {
      mailError.value = e;
    }
  }

  return {
    getMails,
    mailError,
    isFetchingMails,
    offset,
    loadMoreLess,
    returnMails,
    searchEmail,
  };
});
