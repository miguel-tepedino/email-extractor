<script lang="ts" setup>
import type { Email } from "./type";
const props = defineProps<{
  email: Email;
}>();
</script>

<template>
  <div class="flex flex-col gap-2 bg-white">
    <div>
      <span>From: </span>
      <div class="ml-5">{{ props.email.From }}</div>
    </div>

    <div>
      <span>To: </span>
      <div class="ml-5">
        {{ props.email.To }}
      </div>
    </div>

    <div class="">
      <span>Subject: </span>
      <div class="ml-5">{{ props.email.Subject }}</div>
    </div>

    <div class="bg-gray-300 p-2 rounded-lg break-words">
      {{ props.email.Body }}
    </div>
  </div>
</template>
