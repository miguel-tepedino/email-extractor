<script setup lang="ts">
import { computed, ref } from "vue";
import { Field, useForm } from "vee-validate";
import * as yup from "yup";
import authStore from "../stores/auth";
import { useRouter } from "vue-router";

const authstore = authStore();

const router = useRouter();

const authError = ref<string | null>(null);

const schema = computed(() => {
  return yup.object({
    email: yup
      .string()
      .email("Invalidad email address")
      .required("Password field is required"),
    password: yup
      .string()
      .min(8, "Min 8 characters")
      .max(16, "Max 16 characters")
      .required(),
  });
});

const { errors, handleSubmit } = useForm<{
  email: string;
  password: string;
}>({
  validationSchema: schema,
  initialValues: {
    email: "",
    password: "",
  },
});

const onSubmit = handleSubmit((values) => {
  try {
    authstore.login(values);
    router.push("/mails");
  } catch (e: any) {
    authError.value = e.message;
  }
});
</script>

<template>
  <main
    class="w-screen h-screen flex flex-row justify-center items-center bg-gradient-to-tr from-cyan-800 to-cyan-400"
  >
    <div class="flex bg-white w-11/12 md:w-auto p-10 rounded-3xl">
      <form @submit.prevent="onSubmit" class="flex flex-col">
        <div class="flex flex-row justify-center">
          <h3 class="text-xl">Welcome</h3>
        </div>
        <div class="mt-5">
          <p>Hint: admin@gmail.com, 12345678</p>
        </div>
        <Field
          class="border px-2 py-2 rounded-lg mt-5"
          name="email"
          :class="!!errors.email ? 'border-red-600' : 'border-blue-500'"
          validate-on-input
          type="text"
        />
        <div class="text-xs text-red-600" v-if="errors.email">
          {{ errors.email }}
        </div>
        <Field
          class="border px-2 py-2 rounded-lg mt-5"
          name="password"
          :class="!!errors.password ? 'border-red-600' : 'border-blue-500'"
          validate-on-input
          type="password"
        />
        <div class="text-xs text-red-600" v-if="errors.password">
          {{ errors.password }}
        </div>
        <div class="text-xs mt-2 text-red-600" v-if="authError">
          {{ authError }}
        </div>
        <span
          v-if="errors.email || errors.password"
          class="mt-5 self-center bg-gray-400 px-3 py-2 rounded-xl"
        >
          Submit
        </span>
        <button
          v-else
          class="mt-5 bg-cyan-500 self-center px-3 py-2 rounded-xl text-white"
          type="submit"
        >
          Submit
        </button>
      </form>
    </div>
  </main>
</template>
