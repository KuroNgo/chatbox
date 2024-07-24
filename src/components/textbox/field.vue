<template>
    <div class="grid gap-6 mb-6 md:grid-cols-2">
        <label for="small-input" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">{{ text }}</label>
        <input type="text" id="small-input" class="block w-full p-2 text-gray-900 border border-gray-300 rounded-lg bg-gray-50 text-xs focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500">
    </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
export default defineComponent ({
    name:"TextField",
    emits: ["update","delete","create"],
    props: {
        text: {
            type: String,
            require: true,
        }
    }
})
</script>