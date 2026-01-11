<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, nextTick, watch } from 'vue';
import { useAppStore } from '@/stores/app';
import { useI18n } from 'vue-i18n';
import type { Article } from '@/types/models';
import { PhImage, PhHeart, PhList, PhFloppyDisk, PhGlobe, PhX } from '@phosphor-icons/vue';
import { openInBrowser } from '@/utils/browser';

const store = useAppStore();
const { t } = useI18n();

interface Props {
  isSidebarOpen?: boolean;
}

defineProps<Props>();

const emit = defineEmits<{
  toggleSidebar: [];
}>();

// Constants
const ITEMS_PER_PAGE = 30;
const SCROLL_THRESHOLD_PX = 500; // Start loading more items when user is 500px from bottom

const articles = ref<Article[]>([]);
const isLoading = ref(false);
const page = ref(1);
const hasMore = ref(true);
const selectedArticle = ref<Article | null>(null);
const showImageViewer = ref(false);
const allImages = ref<string[]>([]);
const currentImageIndex = ref(0);
const currentImageLoading = ref(false);
const columns = ref<Article[][]>([]);
const columnCount = ref(4);
const containerRef = ref<HTMLElement | null>(null);
// eslint-disable-next-line no-undef
let resizeObserver: ResizeObserver | null = null;
const contextMenu = ref<{ show: boolean; x: number; y: number; article: Article | null }>({
  show: false,
  x: 0,
  y: 0,
  article: null,
});
const imageCountCache = ref<Map<number, number>>(new Map());

// Compute which feed ID to fetch (if viewing a specific feed)
const feedId = computed(() => store.currentFeedId);

// Fetch image gallery articles
async function fetchImages(loadMore = false) {
  if (isLoading.value) return;

  isLoading.value = true;
  try {
    let url = `/api/articles/images?page=${page.value}&limit=${ITEMS_PER_PAGE}`;
    if (feedId.value) {
      url += `&feed_id=${feedId.value}`;
    }

    const res = await fetch(url);
    if (res.ok) {
      const data = await res.json();

      // Validate that data is an array
      if (!Array.isArray(data)) {
        console.error('API response is not an array:', data);
        return;
      }

      const newArticles = data;

      if (loadMore) {
        articles.value = [...articles.value, ...newArticles];
      } else {
        articles.value = newArticles;
      }

      hasMore.value = newArticles.length >= ITEMS_PER_PAGE;

      // Preload image counts for new articles
      newArticles.forEach((article: Article) => {
        if (!imageCountCache.value.has(article.id)) {
          fetchImageCount(article.id);
        }
      });
    }
  } catch (e) {
    console.error('Failed to load images:', e);
  } finally {
    isLoading.value = false;
  }
}

// Fetch image count for an article
async function fetchImageCount(articleId: number) {
  try {
    const res = await fetch(`/api/articles/extract-images?id=${articleId}`);
    if (res.ok) {
      const data = await res.json();
      if (data.images && Array.isArray(data.images)) {
        imageCountCache.value.set(articleId, data.images.length);
      }
    }
  } catch (e) {
    console.error('Failed to fetch image count:', e);
  }
}

// Get image count for an article
function getImageCount(article: Article): number {
  return imageCountCache.value.get(article.id) || 1;
}

// Calculate number of columns based on container width dynamically
function calculateColumns() {
  if (!containerRef.value) return;
  const width = containerRef.value.offsetWidth;

  // Target column width: 250px for optimal image viewing
  // Minimum 2 columns, no maximum
  const targetColumnWidth = 250;
  const calculatedColumns = Math.floor(width / targetColumnWidth);

  // Ensure at least 2 columns
  columnCount.value = Math.max(2, calculatedColumns);

  // Rearrange columns after calculating new count
  arrangeColumns();
}

// Arrange articles into columns by time, balancing heights
function arrangeColumns() {
  if (articles.value.length === 0) {
    columns.value = [];
    return;
  }

  // Initialize columns
  const cols: Article[][] = Array.from({ length: columnCount.value }, () => []);
  const colHeights: number[] = Array(columnCount.value).fill(0);

  // Sort articles by published date (newest first)
  const sortedArticles = [...articles.value].sort((a, b) => {
    return new Date(b.published_at).getTime() - new Date(a.published_at).getTime();
  });

  // Place each article in the shortest column
  sortedArticles.forEach((article) => {
    const shortestColIndex = colHeights.indexOf(Math.min(...colHeights));
    cols[shortestColIndex].push(article);
    // Estimate height: 200px for image + 80px for info
    colHeights[shortestColIndex] += 280;
  });

  columns.value = cols;
}

// Handle scroll for infinite loading
function handleScroll() {
  if (!containerRef.value) return;

  const scrollTop = containerRef.value.scrollTop;
  const containerHeight = containerRef.value.clientHeight;
  const scrollHeight = containerRef.value.scrollHeight;

  if (
    scrollTop + containerHeight >= scrollHeight - SCROLL_THRESHOLD_PX &&
    !isLoading.value &&
    hasMore.value
  ) {
    // Increment page before fetching
    const nextPage = page.value + 1;
    page.value = nextPage;
    fetchImages(true);
  }
}

// Toggle favorite
async function toggleFavorite(article: Article, event: Event) {
  event.stopPropagation();
  try {
    const res = await fetch(`/api/articles/favorite?id=${article.id}`, {
      method: 'POST',
    });
    if (res.ok) {
      article.is_favorite = !article.is_favorite;
    }
  } catch (e) {
    console.error('Failed to toggle favorite:', e);
  }
}

// Open image viewer
async function openImage(article: Article) {
  selectedArticle.value = article;
  showImageViewer.value = true;
  currentImageLoading.value = true;

  // Fetch all images from the article
  await fetchArticleImages(article);

  // Mark as read
  if (!article.is_read) {
    markAsRead(article);
  }
}

// Fetch all images from article content
async function fetchArticleImages(article: Article) {
  try {
    const res = await fetch(`/api/articles/extract-images?id=${article.id}`);
    if (res.ok) {
      const data = await res.json();
      if (data.images && Array.isArray(data.images) && data.images.length > 0) {
        allImages.value = data.images;
        // Find the index of the article's main image
        currentImageIndex.value = data.images.findIndex((img: string) => img === article.image_url);
        if (currentImageIndex.value < 0) {
          currentImageIndex.value = 0;
        }
      } else {
        // Fallback to just the article's main image
        allImages.value = [article.image_url || ''];
        currentImageIndex.value = 0;
      }
    } else {
      // Fallback on error
      allImages.value = [article.image_url || ''];
      currentImageIndex.value = 0;
    }
  } catch (e) {
    console.error('Failed to fetch article images:', e);
    // Fallback on error
    allImages.value = [article.image_url || ''];
    currentImageIndex.value = 0;
  }
}

// Navigate to previous image
function previousImage() {
  if (currentImageIndex.value > 0) {
    currentImageIndex.value--;
  } else {
    // Wrap to last image
    currentImageIndex.value = allImages.value.length - 1;
  }
  // Reset loading state
  currentImageLoading.value = true;
}

// Navigate to next image
function nextImage() {
  if (currentImageIndex.value < allImages.value.length - 1) {
    currentImageIndex.value++;
  } else {
    // Wrap to first image
    currentImageIndex.value = 0;
  }
  // Reset loading state
  currentImageLoading.value = true;
}

// Handle image load
function handleImageLoad() {
  currentImageLoading.value = false;
}

// Handle image error
function handleImageError() {
  currentImageLoading.value = false;
}

// Mark article as read
async function markAsRead(article: Article) {
  try {
    const res = await fetch(`/api/articles/read?id=${article.id}&read=true`, {
      method: 'POST',
    });
    if (res.ok) {
      article.is_read = true;
    }
  } catch (e) {
    console.error('Failed to mark as read:', e);
  }
}

// Close image viewer
function closeImageViewer() {
  showImageViewer.value = false;
  selectedArticle.value = null;
  allImages.value = [];
  currentImageIndex.value = 0;
}

// Format date
function formatDate(dateString: string): string {
  const date = new Date(dateString);
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  const days = Math.floor(diff / (1000 * 60 * 60 * 24));

  if (days === 0) {
    const hours = Math.floor(diff / (1000 * 60 * 60));
    if (hours === 0) {
      const minutes = Math.floor(diff / (1000 * 60));
      return minutes <= 0 ? t('justNow') : t('minutesAgo', { count: minutes });
    }
    return t('hoursAgo', { count: hours });
  } else if (days < 7) {
    return t('daysAgo', { count: days });
  }
  return date.toLocaleDateString();
}

// Handle right-click context menu
function handleContextMenu(event: MouseEvent, article: Article) {
  event.preventDefault();
  event.stopPropagation();
  contextMenu.value = {
    show: true,
    x: event.clientX,
    y: event.clientY,
    article,
  };
}

// Close context menu
function closeContextMenu() {
  contextMenu.value.show = false;
}

// Download image
async function downloadImage(article: Article) {
  try {
    const response = await fetch(article.image_url || '');
    const blob = await response.blob();
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${article.title}.jpg`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    window.URL.revokeObjectURL(url);
  } catch (e) {
    console.error('Failed to download image:', e);
    window.showToast(t('downloadFailed'), 'error');
  }
  closeContextMenu();
}

// Open original article
function openOriginal(article: Article) {
  openInBrowser(article.url);
  closeContextMenu();
}

// Handle keyboard shortcuts
function handleKeyDown(e: KeyboardEvent) {
  // Only handle keyboard when image viewer is open
  if (!showImageViewer.value) return;

  if (e.key === 'Escape') {
    closeImageViewer();
  } else if (e.key === 'ArrowLeft') {
    e.preventDefault();
    previousImage();
  } else if (e.key === 'ArrowRight') {
    e.preventDefault();
    nextImage();
  }
}

// Watch for articles changes and rearrange
watch(articles, () => {
  nextTick(() => {
    arrangeColumns();
  });
});

// Watch for feed ID changes and refetch
watch(feedId, async () => {
  // Close image viewer when switching feeds
  showImageViewer.value = false;
  selectedArticle.value = null;
  allImages.value = [];
  currentImageIndex.value = 0;

  page.value = 1;
  articles.value = [];
  hasMore.value = true;
  await fetchImages();
  // Recalculate columns after fetching new articles
  await nextTick();
  calculateColumns();
});

onMounted(() => {
  fetchImages();
  if (containerRef.value) {
    containerRef.value.addEventListener('scroll', handleScroll);
  }
  window.addEventListener('click', closeContextMenu);
  window.addEventListener('keydown', handleKeyDown);

  // Set up ResizeObserver to watch for container size changes
  if (containerRef.value) {
    // eslint-disable-next-line no-undef
    resizeObserver = new ResizeObserver(() => {
      calculateColumns();
    });
    resizeObserver.observe(containerRef.value);
  }
});

onUnmounted(() => {
  if (containerRef.value) {
    containerRef.value.removeEventListener('scroll', handleScroll);
  }
  if (resizeObserver && containerRef.value) {
    resizeObserver.unobserve(containerRef.value);
    resizeObserver.disconnect();
    resizeObserver = null;
  }
  window.removeEventListener('click', closeContextMenu);
  window.removeEventListener('keydown', handleKeyDown);
});
</script>

<template>
  <div ref="containerRef" class="image-gallery-container">
    <!-- Header -->
    <div class="gallery-header">
      <button class="menu-btn md:hidden" :title="t('toggleSidebar')" @click="emit('toggleSidebar')">
        <PhList :size="24" />
      </button>
      <div class="flex items-center gap-2">
        <PhImage :size="24" class="text-accent" />
        <h1 class="text-xl font-bold text-text-primary">{{ t('imageGallery') }}</h1>
      </div>
    </div>

    <!-- Masonry Grid -->
    <div v-if="articles.length > 0" class="masonry-container">
      <div v-for="(column, colIndex) in columns" :key="colIndex" class="masonry-column">
        <div
          v-for="article in column"
          :key="article.id"
          class="masonry-item"
          @click="openImage(article)"
          @contextmenu="handleContextMenu($event, article)"
        >
          <div class="image-container">
            <img
              :src="article.image_url"
              :alt="article.title"
              class="gallery-image"
              loading="lazy"
            />
            <!-- Image count indicator -->
            <div v-if="getImageCount(article) > 1" class="image-count-badge">
              <PhImage :size="14" />
              <span class="ml-1">{{ getImageCount(article) }}</span>
            </div>
            <div class="image-overlay">
              <button class="favorite-btn" @click="toggleFavorite(article, $event)">
                <PhHeart
                  :size="20"
                  :weight="article.is_favorite ? 'fill' : 'regular'"
                  :class="article.is_favorite ? 'text-red-500' : 'text-white'"
                />
              </button>
            </div>
          </div>
          <div class="image-info">
            <p class="image-title">{{ article.title }}</p>
            <div class="image-meta">
              <span class="feed-name">{{ article.feed_title }}</span>
              <span class="image-date">{{ formatDate(article.published_at) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else-if="!isLoading" class="empty-state">
      <PhImage :size="64" class="text-text-secondary opacity-50" />
      <p class="text-text-secondary">{{ t('noArticles') }}</p>
    </div>

    <!-- Loading Indicator -->
    <div v-if="isLoading" class="loading-indicator">
      <div class="spinner"></div>
    </div>

    <!-- Image Viewer Modal -->
    <div
      v-if="showImageViewer && selectedArticle"
      class="image-viewer-modal"
      data-image-viewer="true"
      @click="closeImageViewer"
    >
      <button class="close-btn" @click="closeImageViewer">
        <PhX :size="20" />
      </button>

      <!-- Image counter (when multiple images) -->
      <div v-if="allImages.length > 1" class="image-counter">
        {{ currentImageIndex + 1 }} / {{ allImages.length }}
      </div>

      <!-- Navigation buttons (when multiple images) -->
      <template v-if="allImages.length > 1">
        <button class="nav-btn nav-btn-prev" @click.stop="previousImage">‹</button>
        <button class="nav-btn nav-btn-next" @click.stop="nextImage">›</button>
      </template>

      <div class="image-viewer-content" @click.stop>
        <!-- Loading placeholder -->
        <div v-if="currentImageLoading" class="image-loading-placeholder">
          <div class="loading-spinner"></div>
        </div>

        <img
          :src="allImages[currentImageIndex] || selectedArticle.image_url"
          :alt="selectedArticle.title"
          class="viewer-image"
          :class="{ loading: currentImageLoading }"
          @load="handleImageLoad"
          @error="handleImageError"
        />
      </div>

      <div class="viewer-info" @click.stop>
        <div class="viewer-header">
          <h2 class="viewer-title">{{ selectedArticle.title }}</h2>
          <a
            :href="selectedArticle.url"
            target="_blank"
            rel="noopener noreferrer"
            class="view-original-btn"
          >
            {{ t('viewOriginal') }}
          </a>
        </div>
        <div class="viewer-meta">
          <span class="feed-name">{{ selectedArticle.feed_title }}</span>
          <span class="image-date">{{ formatDate(selectedArticle.published_at) }}</span>
        </div>
      </div>
    </div>

    <!-- Context Menu -->
    <div
      v-if="contextMenu.show && contextMenu.article"
      class="context-menu"
      :style="{ top: contextMenu.y + 'px', left: contextMenu.x + 'px' }"
      @click.stop
    >
      <button class="context-menu-item" @click="downloadImage(contextMenu.article)">
        <PhFloppyDisk :size="16" />
        <span>{{ t('downloadImage') }}</span>
      </button>
      <button class="context-menu-item" @click="openOriginal(contextMenu.article)">
        <PhGlobe :size="16" />
        <span>{{ t('openInBrowser') }}</span>
      </button>
    </div>
  </div>
</template>

<style scoped>
.image-gallery-container {
  @apply flex flex-col flex-1 h-full overflow-y-auto bg-bg-primary scroll-smooth;
}

.gallery-header {
  @apply sticky top-0 z-10 bg-bg-primary border-b border-border px-4 py-3 flex items-center gap-3;
}

.menu-btn {
  @apply p-2 rounded-lg hover:bg-bg-tertiary text-text-primary transition-colors;
}

.masonry-container {
  @apply p-4 flex gap-4;
}

.masonry-column {
  @apply flex-1 flex flex-col gap-4;
}

.masonry-item {
  @apply cursor-pointer;
}

.image-container {
  @apply relative overflow-hidden rounded-lg bg-bg-secondary;
  @apply transition-transform duration-200 hover:scale-[1.02];
}

.gallery-image {
  @apply w-full h-auto block;
}

.image-count-badge {
  @apply absolute bottom-2 left-2 px-2 py-1 rounded-full;
  @apply bg-black/60 text-white text-xs font-semibold;
  @apply backdrop-blur-sm z-10;
  @apply flex items-center gap-1;
}

.image-overlay {
  @apply absolute inset-0 bg-black/0 hover:bg-black/30 transition-all duration-200;
  @apply flex items-start justify-end p-2;
}

.favorite-btn {
  @apply opacity-0 hover:opacity-100 transition-opacity duration-200;
  @apply bg-black/50 rounded-full p-1.5;
}

.image-container:hover .favorite-btn {
  @apply opacity-100;
}

.image-info {
  @apply p-2;
}

.image-title {
  @apply text-sm font-medium text-text-primary line-clamp-2 mb-1;
}

.image-meta {
  @apply flex items-center justify-between text-xs text-text-secondary;
}

.feed-name {
  @apply truncate flex-1;
}

.image-date {
  @apply ml-2 shrink-0;
}

.empty-state {
  @apply flex flex-col items-center justify-center h-full w-full gap-4;
}

.loading-indicator {
  @apply flex justify-center py-8;
}

.spinner {
  @apply w-8 h-8 border-4 border-accent border-t-transparent rounded-full animate-spin;
}

/* Image Viewer Modal */
.image-viewer-modal {
  @apply fixed inset-0 z-50 bg-black/90 flex flex-col p-4;
}

.close-btn {
  @apply absolute top-4 right-4 w-8 h-8 bg-black/50 hover:bg-black/70;
  @apply rounded-full text-white flex items-center justify-center;
  @apply transition-colors duration-200 z-10 shrink-0;
}

/* Image counter */
.image-counter {
  @apply absolute top-4 left-4 px-2 py-1 rounded;
  @apply text-white text-sm font-medium min-w-[60px] text-center;
  @apply z-10 text-shadow;
}

/* Navigation buttons */
.nav-btn {
  @apply absolute top-1/2 -translate-y-1/2 w-12 h-12 rounded;
  @apply text-white text-4xl;
  @apply flex items-center justify-center transition-all duration-200;
  @apply z-10 text-shadow;
}

.nav-btn-prev {
  @apply left-4;
}

.nav-btn-next {
  @apply right-4;
}

.nav-btn:hover {
  @apply scale-110;
}

.nav-btn:active {
  @apply scale-95;
}

/* Text shadow for better visibility on images */
.text-shadow {
  text-shadow:
    0 1px 3px rgba(0, 0, 0, 0.8),
    0 1px 2px rgba(0, 0, 0, 0.6);
}

.image-viewer-content {
  @apply flex-1 flex items-center justify-center min-h-0 relative;
}

.image-loading-placeholder {
  @apply absolute inset-0 flex items-center justify-center z-10;
}

.loading-spinner {
  @apply w-12 h-12 border-4 border-white/20 border-t-white rounded-full animate-spin;
}

.viewer-image {
  @apply h-full w-full;
  object-fit: contain;
}

.viewer-image.loading {
  @apply opacity-0;
}

.viewer-info {
  @apply bg-bg-primary px-3 py-3 rounded-md shrink-0;
  background-color: var(--bg-primary);
}

.viewer-header {
  @apply flex items-center justify-between gap-4 mb-2;
}

.viewer-title {
  @apply text-base font-bold text-text-primary flex-1 line-clamp-2;
}

.viewer-meta {
  @apply flex items-center gap-4 text-sm text-text-secondary;
}

.view-original-btn {
  @apply px-3 py-1.5 bg-accent text-white rounded-md text-sm whitespace-nowrap;
  @apply hover:bg-accent-hover transition-colors duration-200;
}

/* Context Menu */
.context-menu {
  @apply fixed z-50 bg-bg-primary border border-border rounded-lg shadow-lg py-1 min-w-[180px];
}

.context-menu-item {
  @apply w-full px-4 py-2 flex items-center gap-3 text-sm text-text-primary;
  @apply hover:bg-bg-tertiary transition-colors cursor-pointer;
}

.context-menu-item:active {
  @apply bg-bg-secondary;
}
</style>
