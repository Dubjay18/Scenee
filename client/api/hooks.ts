/**
 * React Query Hooks for Scenee API
 * Provides type-safe hooks for data fetching and mutations
 */

import {
  useQuery,
  useMutation,
  useQueryClient,
  UseQueryOptions,
  UseMutationOptions,
} from '@tanstack/react-query';

import {
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
  ApiRequestError,
} from './client';

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
} from './types';

// ============================================================================
// Query Keys
// ============================================================================

export const queryKeys = {
  // Auth
  authUser: ['auth', 'user'] as const,

  // User
  me: ['user', 'me'] as const,

  // Movies
  movies: ['movies'] as const,
  movie: (id: number) => ['movies', id] as const,
  movieSearch: (params: SearchMoviesParams) => ['movies', 'search', params] as const,

  // Watchlists
  watchlists: ['watchlists'] as const,
  watchlist: (id: string) => ['watchlists', id] as const,
  watchlistPublic: (slug: string) => ['watchlists', 'public', slug] as const,
  watchlistsByOwner: (params?: ListWatchlistsParams) => ['watchlists', 'list', params] as const,
  trending: (params?: TrendingParams) => ['watchlists', 'trending', params] as const,

  // Reviews
  reviewsByMovie: (movieId: string) => ['reviews', 'movie', movieId] as const,

  // Social
  followers: (userId: string) => ['users', userId, 'followers'] as const,
  following: (userId: string) => ['users', userId, 'following'] as const,

  // Notifications
  notifications: (params?: NotificationsParams) => ['notifications', params] as const,

  // Discover
  discoverTrending: (params?: DiscoverParams) => ['discover', 'trending', params] as const,
  discoverNew: (params?: DiscoverParams) => ['discover', 'new', params] as const,

  // Feed
  feed: (params?: FeedParams) => ['feed', params] as const,

  // Stats
  stats: ['stats'] as const,
};

// ============================================================================
// Auth Hooks
// ============================================================================

export const useRegister = (
  options?: UseMutationOptions<User, ApiRequestError, RegisterRequest>
) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: authApi.register,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.authUser });
    },
    ...options,
  });
};

export const useLogin = (
  options?: UseMutationOptions<LoginResponse, ApiRequestError, LoginRequest>
) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: authApi.login,
    onSuccess: (data) => {
      queryClient.setQueryData(queryKeys.me, data.user);
      queryClient.invalidateQueries({ queryKey: queryKeys.authUser });
    },
    ...options,
  });
};

export const useLogout = (
  options?: UseMutationOptions<LogoutResponse, ApiRequestError, void>
) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: authApi.logout,
    onSuccess: () => {
      queryClient.clear();
    },
    ...options,
  });
};

export const useAuthUser = (
  options?: Omit<UseQueryOptions<User, ApiRequestError>, 'queryKey' | 'queryFn'>
) => {
  return useQuery({
    queryKey: queryKeys.authUser,
    queryFn: authApi.getUser,
    retry: false,
    staleTime: 5 * 60 * 1000, // 5 minutes
    ...options,
  });
};

// ============================================================================
// User Hooks
// ============================================================================

export const useMe = (
  options?: Omit<UseQueryOptions<User, ApiRequestError>, 'queryKey' | 'queryFn'>
) => {
  return useQuery({
    queryKey: queryKeys.me,
    queryFn: userApi.me,
    staleTime: 5 * 60 * 1000, // 5 minutes
    ...options,
  });
};

export const useUpdateMe = (
  options?: UseMutationOptions<UpdateUserResponse, ApiRequestError, UpdateUserRequest>
) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: userApi.updateMe,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.me });
    },
    ...options,
  });
};

// ============================================================================
// Movie Hooks
// ============================================================================

export const useSearchMovies = (
  params: SearchMoviesParams,
  options?: Omit<UseQueryOptions<SearchMoviesResponse, ApiRequestError>, 'queryKey' | 'queryFn'>
) => {
  return useQuery({
    queryKey: queryKeys.movieSearch(params),
    queryFn: () => movieApi.search(params),
    enabled: !!params.q,
    ...options,
  });
};

export const useMovie = (
  id: number,
  options?: Omit<UseQueryOptions<Movie, ApiRequestError>, 'queryKey' | 'queryFn'>
) => {
  return useQuery({
    queryKey: queryKeys.movie(id),
    queryFn: () => movieApi.get(id),
    enabled: id > 0,
    ...options,
  });
};

// ============================================================================
// Watchlist Hooks
// ============================================================================

export const useWatchlists = (
  params?: ListWatchlistsParams,
  options?: Omit<UseQueryOptions<Watchlist[], ApiRequestError>, 'queryKey' | 'queryFn'>
) => {
  return useQuery({
    queryKey: queryKeys.watchlistsByOwner(params),
    queryFn: () => watchlistApi.list(params),
    ...options,
  });
};

export const useWatchlist = (
  id: string,
  options?: Omit<UseQueryOptions<WatchlistWithItems, ApiRequestError>, 'queryKey' | 'queryFn'>
) => {
  return useQuery({
    queryKey: queryKeys.watchlist(id),
    queryFn: () => watchlistApi.get(id),
    enabled: !!id,
    ...options,
  });
};

export const usePublicWatchlist = (
  slug: string,
  options?: Omit<UseQueryOptions<WatchlistWithItems, ApiRequestError>, 'queryKey' | 'queryFn'>
) => {
  return useQuery({
    queryKey: queryKeys.watchlistPublic(slug),
    queryFn: () => watchlistApi.getPublic(slug),
    enabled: !!slug,
    ...options,
  });
};

export const useCreateWatchlist = (
  options?: UseMutationOptions<Watchlist, ApiRequestError, CreateWatchlistRequest>
) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: watchlistApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.watchlists });
    },
    ...options,
  });
};

export const useUpdateWatchlist = (
  options?: UseMutationOptions<Watchlist, ApiRequestError, { id: string; data: UpdateWatchlistRequest }>
) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }) => watchlistApi.update(id, data),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.watchlist(data.id) });
      queryClient.invalidateQueries({ queryKey: queryKeys.watchlists });
    },
    ...options,
  });
};

export const useDeleteWatchlist = (
  options?: UseMutationOptions<void, ApiRequestError, string>
) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: watchlistApi.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.watchlists });
    },
    ...options,
  });
};

export const useAddWatchlistItem = (
  options?: UseMutationOptions<WatchlistItem, ApiRequestError, { watchlistId: string; data: AddWatchlistItemRequest }>
) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ watchlistId, data }) => watchlistApi.addItem(watchlistId, data),
    onSuccess: (_, { watchlistId }) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.watchlist(watchlistId) });
    },
    ...options,
  });
};

export const useRemoveWatchlistItem = (
  options?: UseMutationOptions<void, ApiRequestError, { watchlistId: string; itemId: string }>
) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ watchlistId, itemId }) => watchlistApi.removeItem(watchlistId, itemId),
    onSuccess: (_, { watchlistId }) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.watchlist(watchlistId) });
    },
    ...options,
  });
};

export const useLikeWatchlist = (
  options?: UseMutationOptions<void, ApiRequestError, string>
) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: watchlistApi.like,
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.watchlist(id) });
    },
    ...options,
  });
};

export const useUnlikeWatchlist = (
  options?: UseMutationOptions<void, ApiRequestError, string>
) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: watchlistApi.unlike,
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.watchlist(id) });
    },
    ...options,
  });
};

export const useSaveWatchlist = (
  options?: UseMutationOptions<void, ApiRequestError, string>
) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: watchlistApi.save,
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.watchlist(id) });
    },
    ...options,
  });
};

export const useTrendingWatchlists = (
  params?: TrendingParams,
  options?: Omit<UseQueryOptions<Watchlist[], ApiRequestError>, 'queryKey' | 'queryFn'>
) => {
  return useQuery({
    queryKey: queryKeys.trending(params),
    queryFn: () => watchlistApi.trending(params),
    ...options,
  });
};

// ============================================================================
// Review Hooks
// ============================================================================

export const useMovieReviews = (
  movieId: string,
  options?: Omit<UseQueryOptions<Review[], ApiRequestError>, 'queryKey' | 'queryFn'>
) => {
  return useQuery({
    queryKey: queryKeys.reviewsByMovie(movieId),
    queryFn: () => reviewApi.getByMovie(movieId),
    enabled: !!movieId,
    ...options,
  });
};

export const useCreateReview = (
  options?: UseMutationOptions<Review, ApiRequestError, { movieId: string; data: CreateReviewRequest }>
) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ movieId, data }) => reviewApi.create(movieId, data),
    onSuccess: (_, { movieId }) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.reviewsByMovie(movieId) });
    },
    ...options,
  });
};

export const useUpdateReview = (
  options?: UseMutationOptions<Review, ApiRequestError, { movieId: string; data: UpdateReviewRequest }>
) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ movieId, data }) => reviewApi.update(movieId, data),
    onSuccess: (_, { movieId }) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.reviewsByMovie(movieId) });
    },
    ...options,
  });
};

export const useDeleteReview = (
  options?: UseMutationOptions<void, ApiRequestError, { movieId: string; reviewId: string }>
) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ movieId, reviewId }) => reviewApi.delete(movieId, reviewId),
    onSuccess: (_, { movieId }) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.reviewsByMovie(movieId) });
    },
    ...options,
  });
};

// ============================================================================
// Social Hooks (Follow)
// ============================================================================

export const useFollow = (
  options?: UseMutationOptions<FollowResponse, ApiRequestError, string>
) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: followApi.follow,
    onSuccess: (_, userId) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.followers(userId) });
      queryClient.invalidateQueries({ queryKey: queryKeys.following(userId) });
    },
    ...options,
  });
};

export const useUnfollow = (
  options?: UseMutationOptions<FollowResponse, ApiRequestError, string>
) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: followApi.unfollow,
    onSuccess: (_, userId) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.followers(userId) });
      queryClient.invalidateQueries({ queryKey: queryKeys.following(userId) });
    },
    ...options,
  });
};

export const useFollowers = (
  userId: string,
  options?: Omit<UseQueryOptions<User[], ApiRequestError>, 'queryKey' | 'queryFn'>
) => {
  return useQuery({
    queryKey: queryKeys.followers(userId),
    queryFn: () => followApi.getFollowers(userId),
    enabled: !!userId,
    ...options,
  });
};

export const useFollowing = (
  userId: string,
  options?: Omit<UseQueryOptions<User[], ApiRequestError>, 'queryKey' | 'queryFn'>
) => {
  return useQuery({
    queryKey: queryKeys.following(userId),
    queryFn: () => followApi.getFollowing(userId),
    enabled: !!userId,
    ...options,
  });
};

// ============================================================================
// Notification Hooks
// ============================================================================

export const useNotifications = (
  params?: NotificationsParams,
  options?: Omit<UseQueryOptions<NotificationsResponse, ApiRequestError>, 'queryKey' | 'queryFn'>
) => {
  return useQuery({
    queryKey: queryKeys.notifications(params),
    queryFn: () => notificationApi.list(params),
    ...options,
  });
};

export const useMarkNotificationAsRead = (
  options?: UseMutationOptions<void, ApiRequestError, string>
) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: notificationApi.markAsRead,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] });
    },
    ...options,
  });
};

// ============================================================================
// AI Hooks
// ============================================================================

export const useAskAI = (
  options?: UseMutationOptions<AskAIResponse, ApiRequestError, AskAIRequest>
) => {
  return useMutation({
    mutationFn: aiApi.ask,
    ...options,
  });
};

// ============================================================================
// Discover Hooks
// ============================================================================

export const useDiscoverTrending = (
  params?: DiscoverParams,
  options?: Omit<UseQueryOptions<Watchlist[], ApiRequestError>, 'queryKey' | 'queryFn'>
) => {
  return useQuery({
    queryKey: queryKeys.discoverTrending(params),
    queryFn: () => discoverApi.trending(params),
    ...options,
  });
};

export const useDiscoverNew = (
  params?: DiscoverParams,
  options?: Omit<UseQueryOptions<Movie[], ApiRequestError>, 'queryKey' | 'queryFn'>
) => {
  return useQuery({
    queryKey: queryKeys.discoverNew(params),
    queryFn: () => discoverApi.new(params),
    ...options,
  });
};

// ============================================================================
// Feed Hooks
// ============================================================================

export const useFeed = (
  params?: FeedParams,
  options?: Omit<UseQueryOptions<FeedResponse, ApiRequestError>, 'queryKey' | 'queryFn'>
) => {
  return useQuery({
    queryKey: queryKeys.feed(params),
    queryFn: () => feedApi.get(params),
    ...options,
  });
};

// ============================================================================
// Stats Hooks
// ============================================================================

export const useStats = (
  options?: Omit<UseQueryOptions<StatsResponse, ApiRequestError>, 'queryKey' | 'queryFn'>
) => {
  return useQuery({
    queryKey: queryKeys.stats,
    queryFn: statsApi.get,
    ...options,
  });
};

// ============================================================================
// Admin Hooks
// ============================================================================

export const useDeleteUser = (
  options?: UseMutationOptions<DeleteUserResponse, ApiRequestError, string>
) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: adminApi.deleteUser,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.stats });
    },
    ...options,
  });
};
