<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { useAppStore } from '@/stores/app';
import type { Tag } from '@/types/models';
import { PhCheck } from '@phosphor-icons/vue';

interface Props {
  editingTag: Tag | null;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  close: [];
  save: [name: string, color: string];
}>();

const { t } = useI18n();
const store = useAppStore();

const newTagName = ref('');
const newTagColor = ref('#3B82F6');

// Predefined colors - sorted by color wheel
const predefinedColors = [
  '#EF4444', // red
  '#F97316', // orange
  '#EAB308', // yellow
  '#84CC16', // lime
  '#10B981', // green
  '#06B6D4', // cyan
  '#3B82F6', // blue
  '#6366F1', // indigo
  '#8B5CF6', // violet
  '#EC4899', // pink
];

// Initialize form when editingTag changes
watch(
  () => props.editingTag,
  (tag) => {
    if (tag) {
      newTagName.value = tag.name;
      newTagColor.value = tag.color;
    } else {
      newTagName.value = '';
      newTagColor.value = '#3B82F6';
    }
  },
  { immediate: true }
);

// Computed feeds for the tag being edited
const feedsForTag = computed(() => {
  if (!props.editingTag) return [];
  return store.feeds.filter((f) => f.tags?.some((t) => t.id === props.editingTag!.id));
});

function getFavicon(url: string): string {
  try {
    return `https://www.google.com/s2/favicons?domain=${new URL(url).hostname}`;
  } catch {
    return '';
  }
}

function saveTag() {
  const name = newTagName.value.trim();
  if (!name) {
    return;
  }
  emit('save', name, newTagColor.value);
}

function closeForm() {
  emit('close');
}
</script>

<template>
  <!-- Tag Form Modal -->
  <div
    class="fixed inset-0 z-[110] flex items-center justify-center bg-black/50 backdrop-blur-sm p-2 sm:p-4 animate-fade-in"
    @click.self="closeForm"
  >
    <div
      class="bg-bg-primary rounded-none sm:rounded-2xl shadow-2xl w-full max-w-md max-h-[90vh] flex flex-col"
      @click.stop
    >
      <!-- Header -->
      <div class="p-3 sm:p-5 border-b border-border flex justify-between items-center shrink-0">
        <h3 class="text-base sm:text-lg font-semibold m-0">
          {{ editingTag ? t('modal.tag.editTag') : t('modal.tag.createTag') }}
        </h3>
        <button
          class="text-text-secondary hover:text-text-primary transition-colors text-2xl"
          @click="closeForm"
        >
          &times;
        </button>
      </div>

      <!-- Form -->
      <div class="flex-1 overflow-y-auto p-4 sm:p-6 space-y-4">
        <!-- Name -->
        <div>
          <label class="block mb-1.5 text-sm font-medium text-text-secondary">{{
            t('modal.tag.name')
          }}</label>
          <input
            v-model="newTagName"
            type="text"
            class="input-field w-full"
            :placeholder="t('modal.tag.name')"
            @keyup.enter="saveTag"
          />
        </div>

        <!-- Color -->
        <div>
          <label class="block mb-1.5 text-sm font-medium text-text-secondary">{{
            t('modal.tag.color')
          }}</label>
          <div class="flex gap-2 flex-wrap">
            <button
              v-for="color in predefinedColors"
              :key="color"
              class="w-8 h-8 rounded-md border-2 transition-all hover:scale-110 flex items-center justify-center"
              :class="newTagColor === color ? 'border-accent scale-110' : 'border-transparent'"
              :style="{ backgroundColor: color }"
              @click="newTagColor = color"
            >
              <PhCheck v-if="newTagColor === color" :size="20" weight="bold" class="text-white" />
            </button>
          </div>
        </div>

        <!-- Feeds using this tag (only for editing) -->
        <div v-if="editingTag && feedsForTag.length > 0" class="pt-4 border-t border-border">
          <label class="block mb-2 text-sm font-medium text-text-primary">
            {{ t('modal.tag.assignedFeeds', { count: feedsForTag.length }) }}
          </label>
          <div class="max-h-48 overflow-y-auto space-y-1 pr-1">
            <div
              v-for="feed in feedsForTag"
              :key="feed.id"
              class="flex items-center gap-2 px-2 py-1.5 bg-bg-tertiary rounded text-sm"
            >
              <!-- Favicon -->
              <img
                :src="getFavicon(feed.url)"
                class="w-4 h-4 flex-shrink-0 object-contain"
                @error="
                  (e: Event) => {
                    (e.target as HTMLImageElement).style.display = 'none';
                  }
                "
              />
              <!-- Feed title -->
              <span class="flex-1 min-w-0 truncate text-text-primary">{{ feed.title }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Actions -->
      <div
        class="p-3 sm:p-5 border-t border-border bg-bg-secondary text-right shrink-0 rounded-none sm:rounded-b-2xl"
      >
        <button
          type="button"
          class="text-text-primary hover:text-text-secondary transition-colors text-sm font-medium mr-3"
          @click="closeForm"
        >
          {{ t('common.action.cancel') }}
        </button>
        <button type="button" class="btn-primary text-sm sm:text-base" @click="saveTag">
          {{ editingTag ? t('common.action.save') : t('common.action.add') }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
@reference "../../../../style.css";

.input-field {
  @apply w-full p-2 sm:p-2.5 border border-border rounded-md bg-bg-tertiary text-text-primary text-xs sm:text-sm focus:border-accent focus:outline-none transition-colors;
  box-sizing: border-box;
}

.btn-primary {
  @apply bg-accent text-white border-none px-4 sm:px-5 py-2 sm:py-2.5 rounded-lg cursor-pointer font-semibold hover:bg-accent-hover transition-colors disabled:opacity-70;
}

.animate-fade-in {
  animation: modalFadeIn 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes modalFadeIn {
  from {
    transform: translateY(-20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}
</style>
