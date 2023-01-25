<script lang="ts" setup>
const props = withDefaults(
  defineProps<{
    modelValue: boolean;
  }>(),
  {
    modelValue: true,
  }
);

const emits = defineEmits<{
  (e: "update:modelValue", value: boolean): void;
}>();
</script>

<template>
  <Teleport v-if="props.modelValue" to="body">
    <div
      class="bg-black/75 backdrop-blur-sm absolute top-0 left-0 bottom-0 right-0 flex flex-row justify-center items-center"
    >
      <div
        class="bg-white overflow-y-auto p-5 flex flex-col gap-8 rounded-lg justify-start items-start max-w-[300px] max-w-screen-lg"
      >
        <slot />
        <button
          class="self-center px-4 py-2 rounded-full text-white bg-blue-600"
          @click="emits('update:modelValue', false)"
          type="button"
        >
          Close
        </button>
      </div>
    </div>
  </Teleport>
</template>
