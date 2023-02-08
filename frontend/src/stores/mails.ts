import { defineStore } from "pinia";
import { computed, ref, watch } from "vue";
import type { EmailResponse, MailsHits, MoreLess } from "./types";
import { httpRequest } from "@/http/axios";

export default defineStore("mailsStore", () => {
  const mails = ref<MailsHits[]>([]);

  const mailError = ref<unknown>(null);

  const isFetchingMails = ref<boolean>(false);

  const offset = ref<number>(0);

  const searchword = ref<string>("");

  const totalNumberMails = ref<number>(0);

  const returnMails = computed(() => {
    return mails.value;
  });

  function setSearchWord(word: string) {
    offset.value = 0;
    searchword.value = word;
    searchEmail();
  }

  async function getMails() {
    isFetchingMails.value = true;
    try {
      const response = await httpRequest<EmailResponse>({
        url: "/getmails",
        method: "GET",
        data: undefined,
      });
      if (!(response.data as any).error) {
        totalNumberMails.value = response.data.hits.total.value;
        mails.value = response.data.hits.hits;
      }
    } catch (e) {
      mailError.value = e;
    }
    isFetchingMails.value = false;
  }

  async function loadMoreLess(type: MoreLess) {
    type == "LESS" ? (offset.value -= 10) : (offset.value += 10);
    if (searchword.value.length > 0) {
      searchEmail();
      return;
    }
    try {
      const response = await httpRequest<EmailResponse>({
        url: `/mails/${offset.value}`,
        data: undefined,
        method: "GET",
      });
      if (!(response.data as any).error) {
        mails.value = response.data.hits.hits;
      }
    } catch (e) {
      mailError.value = e;
    }
  }

  async function searchEmail() {
    if (searchword.value.length == 0) {
      getMails();
      return;
    }
    try {
      const response = await httpRequest<EmailResponse>({
        url: "/search",
        method: "POST",
        data: {
          term: searchword.value,
          offset: offset.value,
        },
      });

      if (response.data) {
        mails.value = response.data.hits.hits;
      }
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
    setSearchWord,
    totalNumberMails
  };
});
