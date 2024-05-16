<script lang="ts" setup>
import { Icon } from '@iconify/vue';

const props = defineProps({
  selectedFeed: {
    type: Object,
    default: null
  }
})
</script>

<template>
  <div class="ContentFunction">
    <button class="btn" :title="$t('Home.RssContentButton.backToTop')">
      <Icon icon="material-symbols:vertical-align-top" />
    </button>
    <button class="btn" :title="$t('Home.RssContentButton.bookmark')">
      <Icon icon="material-symbols:collections-bookmark" />
    </button>
    <button class="btn" @click="selectedFeed && $emit('modifyFeedContentReaded', selectedFeed, false)"
      :title="$t('Home.RssContentButton.unread')">
      <Icon icon="material-symbols:thread-unread" />
    </button>
    <button class="btn" :title="$t('Home.RssContentButton.openInBrowser')" :disabled="!selectedFeed"
      :style="!selectedFeed ? { color: 'gray' } : {}">
      <Icon icon="material-symbols:open-in-browser" />
    </button>
    <button class="btn" :title="$t('Home.RssContentButton.share')">
      <Icon icon="material-symbols:share" />
    </button>
    <button class="btn" :title="$t('Home.RssContentButton.translate')">
      <Icon icon="material-symbols:g-translate" />
    </button>
    <button class="btn" :title="$t('Home.RssContentButton.chatWithAI')">
      <Icon icon="material-symbols:robot-2" />
    </button>
  </div>
</template>

<style lang="scss" scoped>
@import '../../styles/home/RssContentButton.scss';
</style>