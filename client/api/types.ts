/**
 * API Types - Request and Response types for Scenee API
 * Based on Go domain types and handler structures
 */

// ============================================================================
// Base Types
// ============================================================================

export interface ApiError {
  error: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  page: number;
  totalPages: number;
  totalCount: number;
}

// ============================================================================
// User Types
// ============================================================================

export interface User {
  id: string;
  created_at: string;
  updated_at: string;
  bio: string;
  email: string;
  username: string;
  avatar_url: string;
}

export interface UpdateUserRequest {
  bio?: string;
  avatar_url?: string;
}

export interface UpdateUserResponse {
  message: string;
}

// ============================================================================
// Auth Types
// ============================================================================

export interface RegisterRequest {
  email: string;
  username: string;
  password: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface LogoutResponse {
  message: string;
}

// ============================================================================
// Movie Types
// ============================================================================

export interface Movie {
  id: string;
  tmdb_id: number;
  title: string;
  year: number;
  release_date?: string;
  overview: string;
  poster_path?: string;
  backdrop_path?: string;
  genres: string[];
  runtime?: number;
  metadata?: Record<string, unknown>;
  created_at: string;
  updated_at: string;
}

export interface SearchMoviesParams {
  q: string;
  page?: number;
}

export interface SearchMoviesResponse {
  movies: Movie[];
  total_count: number;
  total_pages: number;
  page: number;
}

// ============================================================================
// Watchlist Types
// ============================================================================

export type WatchlistVisibility = 'public' | 'private' | 'unlisted';

export interface Watchlist {
  id: string;
  created_at: string;
  updated_at: string;
  owner_id: string;
  owner?: User;
  slug: string;
  title: string;
  description: string;
  cover_url: string;
  like_count: number;
  save_count: number;
  item_count: number;
  saved_by: string[];
  visibility?: WatchlistVisibility;
  tags?: string[];
}

export interface WatchlistItem {
  id: string;
  watchlist_id: string;
  movie_id: string;
  movie?: Movie;
  tmdb_id: number;
  notes: string;
  position: number;
  watched: boolean;
  created_at: string;
  updated_at: string;
}

export interface WatchlistWithItems extends Watchlist {
  items: WatchlistItem[];
}

export interface CreateWatchlistRequest {
  title: string;
  description?: string;
  is_public?: boolean;
}

export interface UpdateWatchlistRequest {
  title?: string;
  description?: string;
  visibility?: WatchlistVisibility;
  tags?: string[];
}

export interface AddWatchlistItemRequest {
  tmdb_id: number;
  notes?: string;
}

export interface ListWatchlistsParams {
  owner?: string;
}

export interface TrendingParams {
  window?: 'week' | 'month';
  limit?: number;
}

// ============================================================================
// Review Types
// ============================================================================

export interface Review {
  id: string;
  user_id: string;
  movie_id: string;
  user?: User;
  rating: number;
  review: string;
  created_at: string;
  updated_at: string;
}

export interface CreateReviewRequest {
  rating: number;
  review?: string;
}

export interface UpdateReviewRequest {
  rating?: number;
  review?: string;
}

// ============================================================================
// Social Types (Follow)
// ============================================================================

export interface FollowResponse {
  message: string;
}

export interface FollowersResponse {
  followers: User[];
}

export interface FollowingResponse {
  following: User[];
}

// ============================================================================
// Notification Types
// ============================================================================

export type NotificationType = 'follow' | 'like' | 'save' | 'comment';

export interface Notification {
  id: string;
  user_id: string;
  type: NotificationType;
  actor_id: string;
  entity_id: string;
  is_read: boolean;
  created_at: string;
}

export interface NotificationsResponse {
  notifications: Notification[];
  count: number;
}

export interface NotificationsParams {
  unread?: boolean;
}

// ============================================================================
// AI Types
// ============================================================================

export interface AskAIRequest {
  query: string;
}

export interface AskAIResponse {
  answer: string;
}

// ============================================================================
// Discover Types
// ============================================================================

export interface DiscoverParams {
  window?: '7d' | '30d' | 'day' | 'week' | 'month';
  page?: number;
  genre?: string;
  region?: string;
}

// ============================================================================
// Feed Types
// ============================================================================

export interface FeedParams {
  page?: number;
  limit?: number;
}

export interface FeedResponse {
  following: Watchlist[];
  results:  Movie[];
  page: number;
  limit: number;
  has_more: boolean;
}

// ============================================================================
// Stats Types
// ============================================================================

export interface StatsResponse {
  total_users: number;
  total_watchlists: number;
  total_reviews: number;
}

// ============================================================================
// Admin Types
// ============================================================================

export interface DeleteUserResponse {
  message: string;
}
