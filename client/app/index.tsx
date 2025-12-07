import { StyleSheet, ImageBackground, Pressable, Dimensions, View, ScrollView } from 'react-native';
import React from 'react';
import { Text } from '../components/Themed';
import { SafeAreaView } from 'react-native-safe-area-context';
import { MaterialIcons } from '@expo/vector-icons';
import { useRouter } from 'expo-router';

const { width } = Dimensions.get('window');
const isSmallScreen = width < 640;

const FeatureCard = ({ icon, title, description }: { icon: keyof typeof MaterialIcons.glyphMap; title: string; description: string }) => (
  <View style={styles.featureCard}>
    <MaterialIcons name={icon} size={28} color="#A855F7" />
    <View style={styles.featureTextContainer}>
      <Text style={styles.featureTitle}>{title}</Text>
      <Text style={styles.featureDescription}>{description}</Text>
    </View>
  </View>
);

const Index = () => {
    const router = useRouter()
  return (
    <ImageBackground
      source={{ uri: 'https://lh3.googleusercontent.com/aida-public/AB6AXuDnjDSobrV453lgix4A6g5jDZx7rtDftgJWDikABXSqYSL6ZFzlxZeY9gew8ZooMqJ2orcJ26ppAJp93Yk_1FMK8TlRlRtDXZjoTksby7QkwzUyCAHxMq4cCUbaVj5KByY0_bc2zoVeUUHhF8N_4GTr8neJCsQxthCfe0SA2Doz8808JM9e8tqXiSAe5m4LY0GmZj_moIkokYpmtqfe0g7JaP1rUqh-241eJ99g6NOqOzJIEzsp3gQXJxjKrbATqBU42W4pGMTiwrE' }}
      style={styles.backgroundImage}
      resizeMode="cover"
    >
      <View style={styles.overlay}>
        <SafeAreaView style={styles.safeArea}>
          <ScrollView 
            contentContainerStyle={styles.scrollContent}
            showsVerticalScrollIndicator={false}
          >
            {/* Header Section */}
            <View style={styles.headerSection}>
              <View style={styles.iconContainer}>
                <MaterialIcons name="movie" size={48} color="#A855F7" />
              </View>
              <Text style={styles.title}>Your Social Movie Universe.</Text>
              <Text style={styles.subtitle}>
                Discover, track, and share movies with friends.
              </Text>
            </View>

            {/* Feature Cards Section */}
            <View style={styles.featuresSection}>
              <View style={styles.featuresGrid}>
                <FeatureCard
                  icon="search"
                  title="Discover"
                  description="Find new and trending films"
                />
                <FeatureCard
                  icon="bookmark-add"
                  title="Track"
                  description="Create personal watchlists"
                />
                <FeatureCard
                  icon="group"
                  title="Share"
                  description="Connect with friends"
                />
              </View>
            </View>

            {/* CTA Section */}
            <View style={styles.ctaSection}>
              <Pressable style={styles.createAccountButton}
                onPress={()=>{
                    router.push("/(auth)/register" as any)
                }}
              >
                <Text style={styles.createAccountText}>Create an Account</Text>
              </Pressable>
              <View style={styles.loginContainer}>
                <Text style={styles.loginText}>Already have an account? </Text>
                <Pressable onPress={() => router.push("/(auth)/login" as any)}>
                  <Text style={styles.loginLink}>Log In</Text>
                </Pressable>
              </View>
            </View>
          </ScrollView>
        </SafeAreaView>
      </View>
    </ImageBackground>
  );
};

export default Index;

const styles = StyleSheet.create({
  backgroundImage: {
    flex: 1,
    width: '100%',
    height: '100%',
  },
  overlay: {
    flex: 1,
    backgroundColor: 'rgba(15, 15, 15, 0.8)',
  },
  safeArea: {
    flex: 1,
  },
  scrollContent: {
    flexGrow: 1,
    paddingHorizontal: 16,
    paddingTop: 40,
    paddingBottom: 24,
    justifyContent: 'space-between',
  },
  headerSection: {
    alignItems: 'center',
  },
  iconContainer: {
    paddingBottom: 32,
  },
  title: {
    fontSize: isSmallScreen ? 32 : 40,
    fontWeight: 'bold',
    color: '#F5F5F5',
    textAlign: 'center',
    letterSpacing: -0.5,
  },
  subtitle: {
    fontSize: 16,
    color: 'rgba(255, 255, 255, 0.8)',
    textAlign: 'center',
    marginTop: 8,
    maxWidth: 300,
  },
  featuresSection: {
    width: '100%',
    maxWidth: 448,
    alignSelf: 'center',
    paddingVertical: 24,
  },
  featuresGrid: {
    flexDirection: isSmallScreen ? 'column' : 'row',
    gap: 12,
  },
  featureCard: {
    flex: 1,
    flexDirection: 'column',
    alignItems: 'center',
    gap: 12,
    borderRadius: 16,
    borderWidth: 1,
    borderColor: 'rgba(168, 85, 247, 0.3)',
    backgroundColor: 'rgba(255, 255, 255, 0.08)',
    padding: 16,
    minHeight: 120,
  },
  featureTextContainer: {
    alignItems: 'center',
    gap: 4,
  },
  featureTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    color: '#F5F5F5',
    textAlign: 'center',
  },
  featureDescription: {
    fontSize: 14,
    color: 'rgba(255, 255, 255, 0.6)',
    textAlign: 'center',
  },
  ctaSection: {
    alignItems: 'center',
    gap: 16,
    paddingTop: 24,
    paddingBottom: 16,
  },
  createAccountButton: {
    width: '100%',
    maxWidth: 384,
    height: 48,
    backgroundColor: '#A855F7',
    borderRadius: 9999,
    alignItems: 'center',
    justifyContent: 'center',
    shadowColor: '#A855F7',
    shadowOffset: { width: 0, height: 0 },
    shadowOpacity: 0.5,
    shadowRadius: 15,
    elevation: 8,
  },
  createAccountText: {
    fontSize: 16,
    fontWeight: 'bold',
    color: '#F5F5F5',
    letterSpacing: 0.3,
  },
  loginContainer: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  loginText: {
    fontSize: 14,
    color: 'rgba(255, 255, 255, 0.6)',
  },
  loginLink: {
    fontSize: 14,
    fontWeight: 'bold',
    color: '#A855F7',
  },
});
