/**
 * API Client - Base HTTP client for Scenee API
 */

import { Platform } from 'react-native';

import type {
  User,
  UpdateUserRequest,
  UpdateUserResponse,
  RegisterRequest,
  LoginRequest,
  LoginResponse,
  LogoutResponse,
  Movie,
  SearchMoviesParams,
  SearchMoviesResponse,
  Watchlist,
  WatchlistWithItems,
  CreateWatchlistRequest,
  UpdateWatchlistRequest,
  AddWatchlistItemRequest,
  WatchlistItem,
  ListWatchlistsParams,
  TrendingParams,
  Review,
  CreateReviewRequest,
  UpdateReviewRequest,
  FollowResponse,
  NotificationsResponse,
  NotificationsParams,
  AskAIRequest,
  AskAIResponse,
  DiscoverParams,
  FeedResponse,
  FeedParams,
  StatsResponse,
  DeleteUserResponse,
  ApiError,
} from './types';

// ============================================================================
// Configuration
// ============================================================================

// Use 10.0.2.2 for Android emulator (maps to host localhost)
// Use localhost for iOS simulator
// For physical devices, replace with your machine's local IP address
const getApiBaseUrl = () => {
  if (Platform.OS === 'android') {
    return 'http://10.0.2.2:8080';
  }
  // iOS simulator or web
  return 'http://localhost:8080';
};

const API_BASE_URL = getApiBaseUrl();
const API_VERSION = '/v1';

// ============================================================================
// Token Management
// ============================================================================

let authToken: string | null = null;

export const setAuthToken = (token: string | null) => {
  authToken = token;
};

export const getAuthToken = () => authToken;

export const clearAuthToken = () => {
  authToken = null;
};

// ============================================================================
// Base HTTP Client
// ============================================================================

type QueryParams = Record<string, string | number | boolean | undefined>;

// Helper to convert typed params to QueryParams
const toQueryParams = <T extends object>(params?: T): QueryParams | undefined => {
  if (!params) return undefined;
  return params as unknown as QueryParams;
};

interface RequestOptions extends Omit<RequestInit, 'body'> {
  body?: unknown;
  params?: QueryParams;
}

class ApiClient {
  private baseUrl: string;
  private apiVersion: string;

  constructor(baseUrl: string, apiVersion: string = '') {
    this.baseUrl = baseUrl;
    this.apiVersion = apiVersion;
  }

  private buildUrl(path: string, params?: Record<string, string | number | boolean | undefined>): string {
    const fullPath = this.apiVersion + path;
    const url = new URL(fullPath, this.baseUrl);
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined) {
          url.searchParams.append(key, String(value));
        }
      });
    }
    return url.toString();
  }

  private async request<T>(path: string, options: RequestOptions = {}): Promise<T> {
    const { body, params, headers: customHeaders, ...rest } = options;

    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...customHeaders,
    };

    if (authToken) {
      (headers as Record<string, string>)['Authorization'] = `Bearer ${authToken}`;
    }

    const config: RequestInit = {
      ...rest,
      headers,
      credentials: 'include', // for cookie-based auth
    };

    if (body) {
      config.body = JSON.stringify(body);
    }

    const url = this.buildUrl(path, params);
    const response = await fetch(url, config);

    if (!response.ok) {
      let errorData: ApiError;
      try {
        errorData = await response.json();
      } catch {
        errorData = { error: `HTTP error ${response.status}` };
      }
      throw new ApiRequestError(errorData.error || 'Unknown error', response.status);
    }

    // Handle 204 No Content
    if (response.status === 204) {
      return {} as T;
    }

    return response.json();
  }

  get<T>(path: string, options?: RequestOptions): Promise<T> {
    return this.request<T>(path, { ...options, method: 'GET' });
  }

  post<T>(path: string, options?: RequestOptions): Promise<T> {
    console.log(API_BASE_URL,path,"pp");
    
    return this.request<T>(path, { ...options, method: 'POST' });
  }

  patch<T>(path: string, options?: RequestOptions): Promise<T> {
    return this.request<T>(path, { ...options, method: 'PATCH' });
  }

  put<T>(path: string, options?: RequestOptions): Promise<T> {
    return this.request<T>(path, { ...options, method: 'PUT' });
  }

  delete<T>(path: string, options?: RequestOptions): Promise<T> {
    return this.request<T>(path, { ...options, method: 'DELETE' });
  }
}

// ============================================================================
// Error Handling
// ============================================================================

export class ApiRequestError extends Error {
  status: number;

  constructor(message: string, status: number) {
    super(message);
    this.name = 'ApiRequestError';
    this.status = status;
  }
}

// ============================================================================
// API Client Instance
// ============================================================================

const client = new ApiClient(API_BASE_URL, API_VERSION);

// ============================================================================
// Auth API
// ============================================================================

export const authApi = {
  register: (data: RegisterRequest) =>
    client.post<User>('/auth/register', { body: data }),

  login: async (data: LoginRequest) => {
    const response = await client.post<LoginResponse>('/auth/login', { body: data });
    if (response.token) {
      setAuthToken(response.token);
    }
    return response;
  },

  logout: async () => {
    const response = await client.post<LogoutResponse>('/auth/logout');
    clearAuthToken();
    return response;
  },

  getUser: () => client.get<User>('/auth/user'),
};

// ============================================================================
// User API
// ============================================================================

export const userApi = {
  me: () => client.get<User>('/me'),

  updateMe: (data: UpdateUserRequest) =>
    client.patch<UpdateUserResponse>('/me', { body: data }),
};

// ============================================================================
// Movie API
// ============================================================================

export const movieApi = {
  search: (params: SearchMoviesParams) =>
    client.get<SearchMoviesResponse>('/search/movies', { params: toQueryParams(params) }),

  get: (id: number) => client.get<Movie>(`/movies/${id}`),
};

// ============================================================================
// Watchlist API
// ============================================================================

export const watchlistApi = {
  list: (params?: ListWatchlistsParams) =>
    client.get<Watchlist[]>('/watchlists', { params: toQueryParams(params) }),

  get: (id: string) => client.get<WatchlistWithItems>(`/watchlists/${id}`),

  getPublic: (slug: string) =>
    client.get<WatchlistWithItems>(`/watchlists/public/${slug}`),

  create: (data: CreateWatchlistRequest) =>
    client.post<Watchlist>('/watchlists', { body: data }),

  update: (id: string, data: UpdateWatchlistRequest) =>
    client.patch<Watchlist>(`/watchlists/${id}`, { body: data }),

  delete: (id: string) => client.delete<void>(`/watchlists/${id}`),

  // Items
  addItem: (watchlistId: string, data: AddWatchlistItemRequest) =>
    client.post<WatchlistItem>(`/watchlists/${watchlistId}/items`, { body: data }),

  removeItem: (watchlistId: string, itemId: string) =>
    client.delete<void>(`/watchlists/${watchlistId}/items/${itemId}`),

  // Likes
  like: (id: string) => client.post<void>(`/watchlists/${id}/like`),

  unlike: (id: string) => client.delete<void>(`/watchlists/${id}/like`),

  // Save
  save: (id: string) => client.post<void>(`/watchlists/${id}/save`),

  // Trending
  trending: (params?: TrendingParams) =>
    client.get<Watchlist[]>('/trending', { params: toQueryParams(params) }),
};

// ============================================================================
// Review API
// ============================================================================

export const reviewApi = {
  getByMovie: (movieId: string) =>
    client.get<Review[]>(`/movies/${movieId}/reviews`),

  create: (movieId: string, data: CreateReviewRequest) =>
    client.post<Review>(`/movies/${movieId}/reviews`, { body: data }),

  update: (movieId: string, data: UpdateReviewRequest) =>
    client.put<Review>(`/movies/${movieId}/reviews`, { body: data }),

  delete: (movieId: string, reviewId: string) =>
    client.delete<void>(`/movies/${movieId}/reviews/${reviewId}`),
};

// ============================================================================
// Follow API
// ============================================================================

export const followApi = {
  follow: (userId: string) =>
    client.post<FollowResponse>(`/users/${userId}/follow`),

  unfollow: (userId: string) =>
    client.delete<FollowResponse>(`/users/${userId}/follow`),

  getFollowers: (userId: string) =>
    client.get<User[]>(`/users/${userId}/followers`),

  getFollowing: (userId: string) =>
    client.get<User[]>(`/users/${userId}/following`),
};

// ============================================================================
// Notification API
// ============================================================================

export const notificationApi = {
  list: (params?: NotificationsParams) =>
    client.get<NotificationsResponse>('/notifications', { params: toQueryParams(params) }),

  markAsRead: (id: string) =>
    client.post<void>(`/notifications/${id}/mark-read`),
};

// ============================================================================
// AI API
// ============================================================================

export const aiApi = {
  ask: (data: AskAIRequest) =>
    client.post<AskAIResponse>('/ai/ask', { body: data }),
};

// ============================================================================
// Discover API
// ============================================================================

export const discoverApi = {
  trending: (params?: DiscoverParams) =>
    client.get<Watchlist[]>('/discover/trending', { params: toQueryParams(params) }),

  new: (params?: DiscoverParams) =>
    client.get<Movie[]>('/discover/new', { params: toQueryParams(params) }),
};

// ============================================================================
// Feed API
// ============================================================================

export const feedApi = {
  get: (params?: FeedParams) =>
    client.get<FeedResponse>('/feed?type=trending', { params: toQueryParams(params) }),
};

// ============================================================================
// Stats API
// ============================================================================

export const statsApi = {
  get: () => client.get<StatsResponse>('/stats'),
};

// ============================================================================
// Admin API
// ============================================================================

export const adminApi = {
  deleteUser: (id: string) =>
    client.delete<DeleteUserResponse>(`/admin/users/${id}`),
};

// ============================================================================
// Exports
// ============================================================================

export {
  API_BASE_URL,
  client as apiClient,
};
