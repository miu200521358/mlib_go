<template>
    <div class="text-right mt-2 mr-2 px-2 py-1">
        <select class="bg-gray-500" v-model="locale" @change="onchangeLanguageHandle(locale)">
            <option v-for="item in languages" :key="item" :class="{ active: item === locale }" :value="item">
                {{ t("languages." + item) }}
            </option>
        </select>
    </div>
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";

const { t, availableLocales: languages, locale } = useI18n();

const onchangeLanguageHandle = (item: string) => {
    locale.value = item;
};
</script>





