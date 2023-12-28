<template>
    <div class="input-layer">
        <div class="button-layer">
            <div>{{ t("model.title") }}</div>
            <button>{{ t("button.history") }}</button>
        </div>
        <div class="file-layer">
            <div>({{ t("model.no-set") }})</div>
            <input id="modelName" type="file" class="button-open" accept=".pmx" />
        </div>
    </div>
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";

const { t } = useI18n();
</script>