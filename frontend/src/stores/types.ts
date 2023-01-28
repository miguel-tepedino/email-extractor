import type { Email } from "@/components/emails/type";

export type EmailResponse = {
  hits: Hits;
  timed_out: boolean;
  took: number;
};

export type Hits = {
  hits: MailsHits[];
  total: {
    value: number;
  };
};

export type MailsHits = {
  _id: string;
  _index: string;
  _source: Email;
};

export type MoreLess = "MORE" | "LESS";
