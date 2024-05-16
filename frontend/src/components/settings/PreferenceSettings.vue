<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { Icon } from '@iconify/vue';

const { t, locale } = useI18n()
const selectedLanguage = ref(locale.value)
const selectedTheme = ref(localStorage.getItem('theme') || 'light')

watch(selectedLanguage, (newLocale: string) => {
  locale.value = newLocale
  localStorage.setItem('locale', newLocale)
})

watch(selectedTheme, (newTheme: string) => {
  let theme = newTheme
  if (newTheme === 'auto') {
    const isDarkMode = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;
    theme = isDarkMode ? 'dark' : 'light'
  }
  localStorage.setItem('theme', theme)
  document.documentElement.className = theme
})
</script>

<template>
  <form>
    <label for="language">{{ t('Settings.PreferenceSettings.language') }}</label>
    <select id="language" name="language" v-model="selectedLanguage">
      <option value="en-GB">
        <Icon icon="flag:sh-4x3" /><span>English(Traditional)</span>
      </option>
      <option value="en-US">
        <Icon icon="flag:us-4x3" /><span>English(Simplified)</span>
      </option>
      <option value="de-DE">
        <Icon icon="flag:de-4x3" /><span>Deutsch</span>
      </option>
      <option value="fr-FR">
        <Icon icon="flag:fr-4x3" /><span>Français</span>
      </option>
      <option value="es-ES">
        <Icon icon="flag:es-4x3" /><span>Español</span>
      </option>
      <option value="zh-CN">
        <Icon icon="flag:cn-4x3" /><span>简体中文</span>
      </option>
      <option value="zh-TW">
        <Icon icon="flag:tw-4x3" /><span>繁體中文</span>
      </option>
      <option value="ja-JP">
        <Icon icon="flag:jp-4x3" /><span>日本語</span>
      </option>
      <option value="ko-KR">
        <Icon icon="flag:kr-4x3" /><span>한국어</span>
      </option>
    </select>
    <label for="theme">{{ t('Settings.PreferenceSettings.theme') }}</label>
    <select id="theme" name="theme" v-model="selectedTheme">
      <option value="light">{{ t('Settings.PreferenceSettings.light') }}</option>
      <option value="dark">{{ t('Settings.PreferenceSettings.dark') }}</option>
      <option value="auto">{{ t('Settings.PreferenceSettings.auto') }}</option>
    </select>
  </form>
</template>

<style lang="scss" scoped>
@import '../../styles/settings/PreferenceSettings.scss';
</style>