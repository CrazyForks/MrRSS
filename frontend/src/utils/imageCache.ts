/**
 * Global image cache to persist loaded images across component refreshes
 * This prevents images from disappearing when the article list refreshes
 */

interface ImageCacheEntry {
  url: string;
  timestamp: number;
}

interface LoadState {
  status: 'loading' | 'loaded' | 'error';
  retryCount: number;
}

class ImageCacheManager {
  private cache = new Map<string, ImageCacheEntry>();
  private loadStates = new Map<string, LoadState>();
  private readonly MAX_RETRIES = 3;
  private readonly CACHE_TTL = 24 * 60 * 60 * 1000; // 24 hours

  /**
   * Get the current URL for an image, returning cached version if available
   */
  getImageUrl(originalUrl: string): string {
    const state = this.loadStates.get(originalUrl);

    // If current load failed but we have a cache, return the cached URL
    if (state?.status === 'error' && this.cache.has(originalUrl)) {
      return this.cache.get(originalUrl)!.url;
    }

    return originalUrl;
  }

  /**
   * Mark an image as successfully loaded and cache it
   */
  markAsLoaded(url: string): void {
    // Cache the successful URL
    this.cache.set(url, {
      url,
      timestamp: Date.now(),
    });

    // Update load state
    this.loadStates.set(url, {
      status: 'loaded',
      retryCount: 0,
    });

    // Clean up old entries
    this.cleanUp();
  }

  /**
   * Handle image load error with retry logic
   */
  handleLoadError(url: string): { shouldRetry: boolean; delay?: number } {
    const state = this.loadStates.get(url) || { status: 'loading', retryCount: 0 };

    // If we have a cached version, restore from cache
    if (this.cache.has(url)) {
      this.loadStates.set(url, {
        status: 'loaded',
        retryCount: 0,
      });
      return { shouldRetry: false };
    }

    // Check if we should retry
    if (state.retryCount < this.MAX_RETRIES) {
      const newRetryCount = state.retryCount + 1;
      this.loadStates.set(url, {
        status: 'loading',
        retryCount: newRetryCount,
      });

      // Exponential backoff: 1s, 2s, 3s
      return {
        shouldRetry: true,
        delay: 1000 * newRetryCount,
      };
    }

    // Final failure - mark as error
    this.loadStates.set(url, {
      status: 'error',
      retryCount: state.retryCount,
    });

    return { shouldRetry: false };
  }

  /**
   * Get the current load state of an image
   */
  getLoadState(url: string): LoadState | undefined {
    return this.loadStates.get(url);
  }

  /**
   * Check if an image has been cached
   */
  hasCached(url: string): boolean {
    return this.cache.has(url);
  }

  /**
   * Reset the retry counter for a URL (used before retrying)
   */
  resetRetry(url: string): void {
    const state = this.loadStates.get(url);
    if (state) {
      this.loadStates.set(url, {
        ...state,
        retryCount: 0,
      });
    }
  }

  /**
   * Clean up old cache entries
   */
  private cleanUp(): void {
    const now = Date.now();
    for (const [url, entry] of this.cache.entries()) {
      if (now - entry.timestamp > this.CACHE_TTL) {
        this.cache.delete(url);
        this.loadStates.delete(url);
      }
    }
  }

  /**
   * Clear all cache (useful for testing or manual refresh)
   */
  clearAll(): void {
    this.cache.clear();
    this.loadStates.clear();
  }
}

// Export a singleton instance
export const imageCache = new ImageCacheManager();
