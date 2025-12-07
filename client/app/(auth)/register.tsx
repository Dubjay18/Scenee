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
import { Link } from 'expo-router';
import { useAuth } from '../../providers/AuthProvider';

const BACKGROUND_DARK = '#0a0a0a';
const SURFACE_DARK = '#161616';
const TEXT_PRIMARY = '#EAEAEA';
const TEXT_SECONDARY = '#A9A9A9';
const NEON_PURPLE = '#A855F7';
const ERROR_RED = '#EF4444';
const BORDER_COLOR = '#2a2a2a';

export default function RegisterScreen() {
  const { register } = useAuth();

  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [isRegistering, setIsRegistering] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleRegister = async () => {
    // Clear previous errors
    setError(null);

    // Validate
    if (!username.trim()) {
      setError('Username is required');
      return;
    }
    if (username.trim().length < 3) {
      setError('Username must be at least 3 characters');
      return;
    }
    if (!email.trim()) {
      setError('Email is required');
      return;
    }
    if (!/\S+@\S+\.\S+/.test(email)) {
      setError('Please enter a valid email');
      return;
    }
    if (!password) {
      setError('Password is required');
      return;
    }
    if (password.length < 6) {
      setError('Password must be at least 6 characters');
      return;
    }
    if (password !== confirmPassword) {
      setError('Passwords do not match');
      return;
    }

    setIsRegistering(true);
    try {
      await register({
        username: username.trim(),
        email: email.trim(),
        password,
      });
      // Navigation is handled by AuthProvider
    } catch (err: any) {
      setError(err?.message || 'Registration failed. Please try again.');
    } finally {
      setIsRegistering(false);
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
            <Text style={styles.title}>Create Account</Text>
            <Text style={styles.subtitle}>Join the movie community</Text>
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

            {/* Username Input */}
            <View style={styles.inputContainer}>
              <MaterialIcons
                name="person"
                size={20}
                color={TEXT_SECONDARY}
                style={styles.inputIcon}
              />
              <TextInput
                style={styles.input}
                placeholder="Username"
                placeholderTextColor={TEXT_SECONDARY}
                value={username}
                onChangeText={setUsername}
                autoCapitalize="none"
                autoComplete="username"
                textContentType="username"
                editable={!isRegistering}
              />
            </View>

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
                editable={!isRegistering}
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
                autoComplete="new-password"
                textContentType="newPassword"
                editable={!isRegistering}
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

            {/* Confirm Password Input */}
            <View style={styles.inputContainer}>
              <MaterialIcons
                name="lock-outline"
                size={20}
                color={TEXT_SECONDARY}
                style={styles.inputIcon}
              />
              <TextInput
                style={styles.input}
                placeholder="Confirm Password"
                placeholderTextColor={TEXT_SECONDARY}
                value={confirmPassword}
                onChangeText={setConfirmPassword}
                secureTextEntry={!showConfirmPassword}
                autoComplete="new-password"
                textContentType="newPassword"
                editable={!isRegistering}
              />
              <Pressable
                onPress={() => setShowConfirmPassword(!showConfirmPassword)}
                style={styles.passwordToggle}
              >
                <MaterialIcons
                  name={showConfirmPassword ? 'visibility-off' : 'visibility'}
                  size={20}
                  color={TEXT_SECONDARY}
                />
              </Pressable>
            </View>

            {/* Register Button */}
            <Pressable
              style={[styles.registerButton, isRegistering && styles.registerButtonDisabled]}
              onPress={handleRegister}
              disabled={isRegistering}
            >
              {isRegistering ? (
                <ActivityIndicator color="#FFFFFF" size="small" />
              ) : (
                <Text style={styles.registerButtonText}>Create Account</Text>
              )}
            </Pressable>

            {/* Terms */}
            <Text style={styles.termsText}>
              By signing up, you agree to our{' '}
              <Text style={styles.termsLink}>Terms of Service</Text> and{' '}
              <Text style={styles.termsLink}>Privacy Policy</Text>
            </Text>
          </View>

          {/* Login Link */}
          <View style={styles.loginContainer}>
            <Text style={styles.loginText}>Already have an account? </Text>
            <Link href={"/(auth)/login" as any} asChild>
              <Pressable>
                <Text style={styles.loginLink}>Sign In</Text>
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
    marginBottom: 32,
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
  registerButton: {
    backgroundColor: NEON_PURPLE,
    borderRadius: 12,
    height: 56,
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: 8,
  },
  registerButtonDisabled: {
    opacity: 0.7,
  },
  registerButtonText: {
    color: '#FFFFFF',
    fontSize: 16,
    fontWeight: '600',
  },
  termsText: {
    color: TEXT_SECONDARY,
    fontSize: 12,
    textAlign: 'center',
    lineHeight: 18,
  },
  termsLink: {
    color: NEON_PURPLE,
  },
  loginContainer: {
    flexDirection: 'row',
    justifyContent: 'center',
    marginTop: 32,
  },
  loginText: {
    color: TEXT_SECONDARY,
    fontSize: 14,
  },
  loginLink: {
    color: NEON_PURPLE,
    fontSize: 14,
    fontWeight: '600',
  },
});
