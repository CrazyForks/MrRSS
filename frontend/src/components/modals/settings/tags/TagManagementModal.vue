<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useAppStore } from '@/stores/app';
import type { Tag } from '@/types/models';
import { PhPlus, PhPencil, PhTrash } from '@phosphor-icons/vue';
import TagFormModal from './TagFormModal.vue';

const { t } = useI18n();
const store = useAppStore();

const emit = defineEmits<{
  close: [];
}>();

// State
const tags = ref<Tag[]>([]);
const editingTag = ref<Tag | null>(null);
const showAddForm = ref(false);

// Fetch tags on mount
onMounted(async () => {
  await fetchTags();
});

async function fetchTags() {
  try {
    const res = await fetch('/api/tags');
    tags.value = await res.json();
  } catch (e) {
    console.error('Failed to fetch tags:', e);
  }
}

function openAddForm() {
  editingTag.value = null;
  showAddForm.value = true;
}

function openEditForm(tag: Tag) {
  editingTag.value = tag;
  showAddForm.value = true;
}

async function handleSaveTag(name: string, color: string) {
  try {
    if (editingTag.value) {
      // Update existing tag
      await fetch('/api/tags/update', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          id: editingTag.value.id,
          name,
          color,
          position: editingTag.value.position,
        }),
      });
    } else {
      // Create new tag
      await fetch('/api/tags', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, color }),
      });
    }

    await fetchTags();
    showAddForm.value = false;
    editingTag.value = null;

    // Refresh store tags
    await store.fetchTags();
  } catch (e) {
    console.error('Failed to save tag:', e);
  }
}

async function deleteTag(tag: Tag) {
  if (!window.confirm(t('modal.tag.confirmDelete'))) {
    return;
  }

  try {
    await fetch('/api/tags/delete', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id: tag.id }),
    });

    await fetchTags();

    // Refresh store tags
    await store.fetchTags();
  } catch (e) {
    console.error('Failed to delete tag:', e);
  }
}

function getFeedCount(tag: Tag): number {
  return store.feeds.filter((f) => f.tags?.some((t) => t.id === tag.id)).length;
}

function handleEditTag(tag: Tag) {
  openEditForm(tag);
}

function handleDeleteTag(tag: Tag) {
  deleteTag(tag);
}

function closeForm() {
  showAddForm.value = false;
  editingTag.value = null;
}
</script>

<template>
  <Teleport to="body">
    <!-- Tag Management Modal -->
    <div
      class="fixed inset-0 z-[100] flex items-center justify-center bg-black/50 backdrop-blur-sm p-2 sm:p-4 animate-fade-in"
      @click.self="emit('close')"
    >
      <div
        class="bg-bg-primary rounded-none sm:rounded-2xl shadow-2xl w-full max-w-2xl max-h-[90vh] flex flex-col"
        @click.stop
      >
        <!-- Header -->
        <div class="p-3 sm:p-5 border-b border-border flex justify-between items-center shrink-0">
          <h3 class="text-base sm:text-lg font-semibold m-0">{{ t('modal.tag.manageTags') }}</h3>
          <button
            class="text-text-secondary hover:text-text-primary transition-colors text-2xl"
            @click="emit('close')"
          >
            &times;
          </button>
        </div>

        <!-- Content -->
        <div class="flex-1 overflow-y-auto p-4 sm:p-6">
          <!-- Tag List -->
          <div v-if="tags.length > 0">
            <!-- Tags displayed in a flex wrap layout -->
            <div class="flex flex-wrap gap-2">
              <div
                v-for="tag in tags"
                :key="tag.id"
                class="inline-flex items-center rounded-l rounded-r border border-border hover:border-accent/50 transition-all overflow-hidden"
              >
                <!-- Left part: Tag name and count with colored background -->
                <div
                  class="flex items-center gap-1.5 px-2.5 py-1.5"
                  :style="{ backgroundColor: tag.color }"
                >
                  <!-- Tag name -->
                  <span class="text-sm font-medium text-white">{{ tag.name }}</span>

                  <!-- Feed count badge -->
                  <span class="feed-count-badge">{{ getFeedCount(tag) }}</span>
                </div>

                <!-- Right part: Action buttons -->
                <div
                  class="flex items-center px-2 py-1.5 bg-bg-secondary gap-1 border-l border-border"
                >
                  <!-- Edit button -->
                  <button
                    class="text-text-secondary hover:text-accent transition-colors"
                    @click="handleEditTag(tag)"
                  >
                    <PhPencil :size="16" />
                  </button>

                  <!-- Delete button -->
                  <button
                    class="text-text-secondary hover:text-red-500 transition-colors"
                    @click="handleDeleteTag(tag)"
                  >
                    <PhTrash :size="16" />
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Empty state -->
          <div v-else class="text-center py-12 text-text-secondary">
            <div class="text-4xl mb-4">🏷️</div>
            <p>{{ t('modal.tag.noTags') }}</p>
            <button
              class="mt-4 px-4 py-2 text-sm font-medium text-white bg-accent rounded-md hover:bg-accent/90 transition-colors"
              @click="openAddForm"
            >
              <PhPlus :size="20" class="inline mr-1" />
              {{ t('modal.tag.createTag') }}
            </button>
          </div>
        </div>

        <!-- Bottom Action Bar -->
        <div
          class="p-3 sm:p-5 border-t border-border bg-bg-secondary text-right shrink-0 rounded-none sm:rounded-b-2xl"
        >
          <button class="btn-primary text-sm sm:text-base" @click="openAddForm">
            <PhPlus :size="20" class="inline mr-1" />
            {{ t('modal.tag.addNew') }}
          </button>
        </div>

        <!-- Tag Form Modal -->
        <TagFormModal
          v-if="showAddForm"
          :editing-tag="editingTag"
          @close="closeForm"
          @save="handleSaveTag"
        />
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
@reference "../../../../style.css";

.btn-primary {
  @apply bg-accent text-white border-none px-4 sm:px-5 py-2 sm:py-2.5 rounded-lg cursor-pointer font-semibold hover:bg-accent-hover transition-colors disabled:opacity-70;
}

.feed-count-badge {
  @apply text-[10px] font-medium rounded-full min-w-[16px] h-[16px] px-1 flex items-center justify-center;
  background-color: rgba(255, 255, 255, 0.25);
  color: #ffffff;
}

.dark-mode .feed-count-badge {
  background-color: rgba(0, 0, 0, 0.25) !important;
  color: #ffffff !important;
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
