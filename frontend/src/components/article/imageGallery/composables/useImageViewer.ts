import { ref, computed } from 'vue';
import type { Article } from '@/types/models';
import type { ImageViewerReturn } from '../types';

const MIN_SCALE = 0.5;
const MAX_SCALE = 5;
const SCALE_STEP = 0.25;

/**
 * Composable for managing image viewer functionality (zoom, pan, navigation)
 * @param allImages - All images from the current article
 * @param currentImageIndex - Index of currently displayed image
 * @param articles - All articles in the gallery
 * @param selectedArticle - Currently selected article
 * @returns Image viewer state and methods
 */
export function useImageViewer(
  allImages: { value: string[] },
  currentImageIndex: { value: number },
  articles: { value: Article[] },
  selectedArticle: { value: Article | null }
): ImageViewerReturn {
  // State
  const scale = ref(1);
  const position = ref<{ x: number; y: number }>({ x: 0, y: 0 });
  const isDragging = ref(false);
  const currentImageLoading = ref(false);
  const dragStart = ref<{ x: number; y: number }>({ x: 0, y: 0 });

  // Computed
  const imageStyle = computed(() => ({
    transform: `translate(${position.value.x}px, ${position.value.y}px) scale(${scale.value})`,
  }));

  // Find current article index in articles array
  const currentArticleIndex = computed(() => {
    if (!selectedArticle.value) return -1;
    return articles.value.findIndex((a) => a.id === selectedArticle.value!.id);
  });

  // Check if we can navigate to previous image/article
  const canNavigatePrevious = computed(() => {
    // Can navigate if not at first image of first article
    if (currentImageIndex.value > 0) return true;
    if (currentArticleIndex.value > 0) return true;
    return false;
  });

  // Check if we can navigate to next image/article
  const canNavigateNext = computed(() => {
    // Can navigate if not at last image of last article
    if (currentImageIndex.value < allImages.value.length - 1) return true;
    if (currentArticleIndex.value >= 0 && currentArticleIndex.value < articles.value.length - 1)
      return true;
    return false;
  });

  // Zoom functions
  function zoomIn(): void {
    if (scale.value < MAX_SCALE) {
      scale.value = Math.min(scale.value + SCALE_STEP, MAX_SCALE);
    }
  }

  function zoomOut(): void {
    if (scale.value > MIN_SCALE) {
      scale.value = Math.max(scale.value - SCALE_STEP, MIN_SCALE);
      // Reset position if zooming out to 1 or less
      if (scale.value <= 1) {
        position.value = { x: 0, y: 0 };
      }
    }
  }

  function resetView(): void {
    scale.value = 1;
    position.value = { x: 0, y: 0 };
  }

  // Drag functions
  function startDrag(e: MouseEvent): void {
    isDragging.value = true;
    dragStart.value = {
      x: e.clientX - position.value.x,
      y: e.clientY - position.value.y,
    };
  }

  function onDrag(e: MouseEvent): void {
    if (isDragging.value) {
      position.value = {
        x: e.clientX - dragStart.value.x,
        y: e.clientY - dragStart.value.y,
      };
    }
  }

  function stopDrag(): void {
    isDragging.value = false;
  }

  // Handle mouse wheel on main image area for navigation
  function handleImageWheel(e: globalThis.WheelEvent): void {
    // Determine direction
    const isNavigatingForward = e.deltaY > 0 || e.deltaX > 0;
    const isNavigatingBackward = e.deltaY < 0 || e.deltaX < 0;

    // Check if navigation is possible
    if (isNavigatingForward && !canNavigateNext.value) return;
    if (isNavigatingBackward && !canNavigatePrevious.value) return;

    // Prevent default scrolling only if we can navigate
    e.preventDefault();

    // Navigation will be handled by parent component
    // We emit the direction via event
    const event = new CustomEvent('image-wheel-navigate', {
      detail: { direction: isNavigatingForward ? 'next' : 'prev' },
    });
    window.dispatchEvent(event);
  }

  // Handle image load
  function handleImageLoad(): void {
    currentImageLoading.value = false;
  }

  // Handle image error
  function handleImageError(): void {
    currentImageLoading.value = false;
  }

  return {
    scale,
    position,
    isDragging,
    currentImageLoading,
    dragStart,
    imageStyle,
    canNavigatePrevious,
    canNavigateNext,
    zoomIn,
    zoomOut,
    resetView,
    startDrag,
    onDrag,
    stopDrag,
    handleImageWheel,
    handleImageLoad,
    handleImageError,
  };
}
