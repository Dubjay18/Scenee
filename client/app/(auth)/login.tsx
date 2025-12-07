import React, { useState } from 'react';
import {
  StyleSheet,
  View,
  TextInput,
  Pressable,
  KeyboardAvoidingView,
  Platform,
  ScrollView,
  ActivityIndicator,
} from 'react-native';
import { Text } from '../../components/Themed';
import { SafeAreaView } from 'react-native-safe-area-context';
import { MaterialIcons } from '@expo/vector-icons';
import { Link, useRouter } from 'expo-router';
import { useAuth } from '../../providers/AuthProvider';

const BACKGROUND_DARK = '#0a0a0a';
const SURFACE_DARK = '#161616';
const TEXT_PRIMARY = '#EAEAEA';
const TEXT_SECONDARY = '#A9A9A9';
const NEON_PURPLE = '#A855F7';
const ERROR_RED = '#EF4444';
const BORDER_COLOR = '#2a2a2a';

export default function LoginScreen() {
  const { login } = useAuth();
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [isLoggingIn, setIsLoggingIn] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleLogin = async () => {
    // Clear previous errors
    setError(null);

    // Validate
    if (!email.trim()) {
      setError('Email is required');
      return;
    }
    if (!password) {
      setError('Password is required');
      return;
    }

    setIsLoggingIn(true);
    try {
      await login({ email: email.trim(), password });
      router.replace('/(tabs)');
    } catch (err: any) {
      setError(err?.message || 'Login failed. Please check your credentials.');
    } finally {
      setIsLoggingIn(false);
    }
  };

  return (
    <SafeAreaView style={styles.container}>
      <KeyboardAvoidingView
        behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
        style={styles.keyboardView}
      >
        <ScrollView
          contentContainerStyle={styles.scrollContent}
          keyboardShouldPersistTaps="handled"
          showsVerticalScrollIndicator={false}
        >
          {/* Header */}
          <View style={styles.header}>
            <View style={styles.iconContainer}>
              <MaterialIcons name="movie" size={48} color={NEON_PURPLE} />
            </View>
            <Text style={styles.title}>Welcome Back</Text>
            <Text style={styles.subtitle}>Sign in to continue to Scenee</Text>
          </View>

          {/* Form */}
          <View style={styles.form}>
            {/* Error Message */}
            {error && (
              <View style={styles.errorContainer}>
                <MaterialIcons name="error-outline" size={18} color={ERROR_RED} />
                <Text style={styles.errorText}>{error}</Text>
              </View>
            )}

            {/* Email Input */}
            <View style={styles.inputContainer}>
              <MaterialIcons
                name="email"
                size={20}
                color={TEXT_SECONDARY}
                style={styles.inputIcon}
              />
              <TextInput
                style={styles.input}
                placeholder="Email"
                placeholderTextColor={TEXT_SECONDARY}
                value={email}
                onChangeText={setEmail}
                autoCapitalize="none"
                autoComplete="email"
                keyboardType="email-address"
                textContentType="emailAddress"
                editable={!isLoggingIn}
              />
            </View>

            {/* Password Input */}
            <View style={styles.inputContainer}>
              <MaterialIcons
                name="lock"
                size={20}
                color={TEXT_SECONDARY}
                style={styles.inputIcon}
              />
              <TextInput
                style={styles.input}
                placeholder="Password"
                placeholderTextColor={TEXT_SECONDARY}
                value={password}
                onChangeText={setPassword}
                secureTextEntry={!showPassword}
                autoComplete="password"
                textContentType="password"
                editable={!isLoggingIn}
              />
              <Pressable
                onPress={() => setShowPassword(!showPassword)}
                style={styles.passwordToggle}
              >
                <MaterialIcons
                  name={showPassword ? 'visibility-off' : 'visibility'}
                  size={20}
                  color={TEXT_SECONDARY}
                />
              </Pressable>
            </View>

            {/* Login Button */}
            <Pressable
              style={[styles.loginButton, isLoggingIn && styles.loginButtonDisabled]}
              onPress={handleLogin}
              disabled={isLoggingIn}
            >
              {isLoggingIn ? (
                <ActivityIndicator color="#FFFFFF" size="small" />
              ) : (
                <Text style={styles.loginButtonText}>Sign In</Text>
              )}
            </Pressable>

            {/* Forgot Password */}
            <Pressable style={styles.forgotPassword}>
              <Text style={styles.forgotPasswordText}>Forgot Password?</Text>
            </Pressable>
          </View>

          {/* Register Link */}
          <View style={styles.registerContainer}>
            <Text style={styles.registerText}>Don't have an account? </Text>
            <Link href={"/(auth)/register" as any} asChild>
              <Pressable>
                <Text style={styles.registerLink}>Sign Up</Text>
              </Pressable>
            </Link>
          </View>
        </ScrollView>
      </KeyboardAvoidingView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: BACKGROUND_DARK,
  },
  keyboardView: {
    flex: 1,
  },
  scrollContent: {
    flexGrow: 1,
    paddingHorizontal: 24,
    paddingVertical: 40,
    justifyContent: 'center',
  },
  header: {
    alignItems: 'center',
    marginBottom: 40,
  },
  iconContainer: {
    width: 80,
    height: 80,
    borderRadius: 20,
    backgroundColor: 'rgba(168, 85, 247, 0.15)',
    alignItems: 'center',
    justifyContent: 'center',
    marginBottom: 24,
  },
  title: {
    fontSize: 28,
    fontWeight: '700',
    color: TEXT_PRIMARY,
    marginBottom: 8,
  },
  subtitle: {
    fontSize: 16,
    color: TEXT_SECONDARY,
  },
  form: {
    gap: 16,
  },
  errorContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: 'rgba(239, 68, 68, 0.1)',
    paddingHorizontal: 16,
    paddingVertical: 12,
    borderRadius: 12,
    gap: 8,
  },
  errorText: {
    color: ERROR_RED,
    fontSize: 14,
    flex: 1,
  },
  inputContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: SURFACE_DARK,
    borderRadius: 12,
    borderWidth: 1,
    borderColor: BORDER_COLOR,
    paddingHorizontal: 16,
    height: 56,
  },
  inputIcon: {
    marginRight: 12,
  },
  input: {
    flex: 1,
    fontSize: 16,
    color: TEXT_PRIMARY,
  },
  passwordToggle: {
    padding: 4,
  },
  loginButton: {
    backgroundColor: NEON_PURPLE,
    borderRadius: 12,
    height: 56,
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: 8,
  },
  loginButtonDisabled: {
    opacity: 0.7,
  },
  loginButtonText: {
    color: '#FFFFFF',
    fontSize: 16,
    fontWeight: '600',
  },
  forgotPassword: {
    alignItems: 'center',
    paddingVertical: 8,
  },
  forgotPasswordText: {
    color: NEON_PURPLE,
    fontSize: 14,
  },
  registerContainer: {
    flexDirection: 'row',
    justifyContent: 'center',
    marginTop: 40,
  },
  registerText: {
    color: TEXT_SECONDARY,
    fontSize: 14,
  },
  registerLink: {
    color: NEON_PURPLE,
    fontSize: 14,
    fontWeight: '600',
  },
});
