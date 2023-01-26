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
      class="bg-black/75 backdrop-blur-sm fixed overflow-auto top-0 left-0 h-full w-full flex flex-row justify-center items-center"
    >
      <div
        class="bg-white overflow-auto p-5 flex flex-col gap-8 rounded-lg h-5/6 w-11/12 md:w-3/4 max-w-screen-lg"
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
