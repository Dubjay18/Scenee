/**
 * Auth Storage - Secure token storage and management
 */

import AsyncStorage from '@react-native-async-storage/async-storage';
import type { User } from '@/api/types';

// Storage keys
const KEYS = {
  ACCESS_TOKEN: '@scenee/access_token',
  USER: '@scenee/user',
} as const;

// ============================================================================
// Token Management
// ============================================================================

/**
 * Store the access token securely
 */
export async function storeAccessToken(token: string): Promise<void> {
  try {
    await AsyncStorage.setItem(KEYS.ACCESS_TOKEN, token);
  } catch (error) {
    console.error('Failed to store access token:', error);
    throw error;
  }
}

/**
 * Retrieve the stored access token
 */
export async function getAccessToken(): Promise<string | null> {
  try {
    return await AsyncStorage.getItem(KEYS.ACCESS_TOKEN);
  } catch (error) {
    console.error('Failed to get access token:', error);
    return null;
  }
}

/**
 * Remove the stored access token
 */
export async function removeAccessToken(): Promise<void> {
  try {
    await AsyncStorage.removeItem(KEYS.ACCESS_TOKEN);
  } catch (error) {
    console.error('Failed to remove access token:', error);
  }
}

// ============================================================================
// User Cache Management
// ============================================================================

/**
 * Store user data for offline access
 */
export async function storeUser(user: User): Promise<void> {
  try {
    await AsyncStorage.setItem(KEYS.USER, JSON.stringify(user));
  } catch (error) {
    console.error('Failed to store user:', error);
  }
}

/**
 * Retrieve cached user data
 */
export async function getCachedUser(): Promise<User | null> {
  try {
    const data = await AsyncStorage.getItem(KEYS.USER);
    return data ? JSON.parse(data) : null;
  } catch (error) {
    console.error('Failed to get cached user:', error);
    return null;
  }
}

/**
 * Remove cached user data
 */
export async function removeCachedUser(): Promise<void> {
  try {
    await AsyncStorage.removeItem(KEYS.USER);
  } catch (error) {
    console.error('Failed to remove cached user:', error);
  }
}

// ============================================================================
// Clear All Auth Data
// ============================================================================

/**
 * Clear all auth-related stored data
 */
export async function clearAuthStorage(): Promise<void> {
  try {
    await AsyncStorage.multiRemove([KEYS.ACCESS_TOKEN, KEYS.USER]);
  } catch (error) {
    console.error('Failed to clear auth storage:', error);
  }
}

/**
 * Check if user has a stored token (quick auth check)
 */
export async function hasStoredToken(): Promise<boolean> {
  const token = await getAccessToken();
  return !!token;
}
