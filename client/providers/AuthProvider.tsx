import React, {
  createContext,
  useContext,
  useEffect,
  useState,
  useCallback,
  ReactNode,
} from 'react';
import { useRouter, useSegments } from 'expo-router';
import { authApi, userApi, setAuthToken, clearAuthToken } from '@/api/client';
import {
  storeAccessToken,
  getAccessToken,
  clearAuthStorage,
  storeUser,
  getCachedUser,
} from '@/lib/authStorage';
import type { User, LoginRequest, RegisterRequest } from '@/api/types';

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (credentials: LoginRequest) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => Promise<void>;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const router = useRouter();
  const segments = useSegments();

  // Initialize auth state from storage
  useEffect(() => {
    const initializeAuth = async () => {
      try {
        // First, try to get cached user for instant UI
        const cachedUser = await getCachedUser();
        if (cachedUser) {
          setUser(cachedUser);
        }

        // Then verify token and refresh user data
        const token = await getAccessToken();
        if (token) {
          setAuthToken(token);
          try {
            const freshUser = await userApi.me();
            setUser(freshUser);
            await storeUser(freshUser);
          } catch (error) {
            // Token is invalid, clear everything
            console.log('Token invalid, clearing auth state');
            await clearAuthStorage();
            clearAuthToken();
            setUser(null);
          }
        } else {
          setUser(null);
        }
      } catch (error) {
        console.error('Error initializing auth:', error);
        setUser(null);
      } finally {
        setIsLoading(false);
      }
    };

    initializeAuth();
  }, []);

  // Handle navigation based on auth state
  useEffect(() => {
    if (isLoading) return;

    const firstSegment = segments[0] as string | undefined;
    const inAuthGroup = firstSegment === '(auth)';
    const isOnLanding = firstSegment === undefined;

    if (!user && !inAuthGroup && !isOnLanding) {
      // Redirect to login if not authenticated and not in auth group or landing
      router.replace('/(auth)/login');
    } else if (user && (inAuthGroup || isOnLanding)) {
      // Redirect to home if authenticated and in auth group or landing
      router.replace('/(tabs)');
    }
  }, [user, segments, isLoading, router]);

  const login = useCallback(async (credentials: LoginRequest) => {
    try {
      const response = await authApi.login(credentials);

      // Store token
      await storeAccessToken(response.token);
      setAuthToken(response.token);

      // Fetch and store user
      const userData = await userApi.me();
      await storeUser(userData);
      setUser(userData);
    } catch (error) {
      console.error('Login error:', error);
      throw error;
    }
  }, []);

  const register = useCallback(async (data: RegisterRequest) => {
    try {
      // Register the user
      await authApi.register(data);

      // After successful registration, log them in automatically
      const loginResponse = await authApi.login({
        email: data.email,
        password: data.password,
      });

      // Store token
      await storeAccessToken(loginResponse.token);
      setAuthToken(loginResponse.token);

      // Fetch and store user
      const userData = await userApi.me();
      await storeUser(userData);
      setUser(userData);
    } catch (error) {
      console.error('Register error:', error);
      throw error;
    }
  }, []);

  const logout = useCallback(async () => {
    try {
      // Call logout endpoint (optional, may fail if token already invalid)
      try {
        await authApi.logout();
      } catch (e) {
        // Ignore logout API errors
      }
    } finally {
      // Always clear local state
      await clearAuthStorage();
      clearAuthToken();
      setUser(null);
    }
  }, []);

  const refreshUser = useCallback(async () => {
    try {
      const userData = await userApi.me();
      await storeUser(userData);
      setUser(userData);
    } catch (error) {
      console.error('Error refreshing user:', error);
      throw error;
    }
  }, []);

  return (
    <AuthContext.Provider
      value={{
        user,
        isLoading,
        isAuthenticated: !!user,
        login,
        register,
        logout,
        refreshUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
