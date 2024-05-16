<script setup lang="ts">
import { useRouter } from 'vue-router';
import { computed } from 'vue';
import { Icon } from '@iconify/vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const router = useRouter();

const goBack = () => {
  router.push('/');
}

const pageTitle = computed(() => {
  switch (router.currentRoute.value.path) {
    case '/settings/rss':
      return t('Settings.rssSettings');
    case '/settings/preference':
      return t('Settings.preferenceSettings');
    case '/settings/about':
      return t('Settings.aboutSettings');
    default:
      return t('Settings.settings');
  }
});
</script>

<template>
  <div class="settings">
    <div class="nav">
      <div class="title">
        <button @click="goBack" class="btn done" :title="$t('Settings.backHome')">
          <Icon icon="material-symbols:arrow-back" />
        </button>
        <h1>{{ pageTitle }}</h1>
      </div>
      <div class="router">
        <router-link to="/settings/rss" class="btn" :title="$t('Settings.rssSettings')">
          <Icon icon="material-symbols:rss-feed" />
        </router-link>
        <router-link to="/settings/preference" class="btn" :title="$t('Settings.preferenceSettings')">
          <Icon icon="material-symbols:inbox-customize" />
        </router-link>
        <router-link to="/settings/about" class="btn" :title="$t('Settings.about')">
          <Icon icon="material-symbols:page-info" />
        </router-link>
      </div>
    </div>
    <router-view />
  </div>
</template>

<style lang="scss" scoped>
@import '../../styles/settings/Settings.scss';
</style>