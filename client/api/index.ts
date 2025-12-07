/**
 * API Module - Exports for Scenee API client
 */

// Types
export * from './types';

// Client
export {
  authApi,
  userApi,
  movieApi,
  watchlistApi,
  reviewApi,
  followApi,
  notificationApi,
  aiApi,
  discoverApi,
  feedApi,
  statsApi,
  adminApi,
  apiClient,
  setAuthToken,
  getAuthToken,
  clearAuthToken,
  ApiRequestError,
  API_BASE_URL,
} from './client';

// Hooks
export * from './hooks';
