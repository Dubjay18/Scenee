import {
  StyleSheet,
  View,
  ScrollView,
  Image,
  Pressable,
  ImageBackground,
  Dimensions,
} from 'react-native';
import { Text } from '../../components/Themed';
import { MaterialIcons } from '@expo/vector-icons';
import { SafeAreaView } from 'react-native-safe-area-context';
import { router, useLocalSearchParams } from 'expo-router';
import { LinearGradient } from 'expo-linear-gradient';

const PRIMARY_COLOR = '#c084fc';
const BACKGROUND_DARK = '#0A0A0B';
const TEXT_PRIMARY = '#FFFFFF';
const TEXT_SECONDARY = '#9CA3AF';
const SURFACE_DARK = 'rgba(39, 39, 42, 0.4)';

const { width } = Dimensions.get('window');

// Mock data for watchlist
const watchlistData = {
  id: '1',
  title: 'Mind-Bending Sci-Fi Thrillers',
  description: 'A curated list of science fiction films that will challenge your perception of reality, time, and consciousness. Prepare to have your mind blown.',
  coverImage: 'https://lh3.googleusercontent.com/aida-public/AB6AXuACITKE7SDcVfRdF4AROh7JZ_ylT8r-GUZ-MsM36aQYsd-fViYysb7rNcXNjuFairi80JKQtDweWuKp_LVYhHvxAgRZgGFRsZzKx50b-UaRhv7CxpSL1GqOXhPTmfVys_zEJic8wlLKUO_gK2ED3bazkxyemMwT4G8a0TwmOUFJP1K_VpB2iup_lHfCNMQRKYGXBDrAX9ngsup7gqo8Y8mJWfkBVUYwMjO_e0XI7Ts9elRRlUSMXVX3UJhPjn5xxnR4clYtHjZuXzY',
  likes: '1.2k',
  creator: {
    name: 'Alex Rivera',
    avatar: 'https://lh3.googleusercontent.com/aida-public/AB6AXuBrrHpY7ynhAIsLtsBM1AQT-ktaWUFzyUeMqZImI8I11krtrxEU9vdSBUhX-f8_Zf-3pMVJQc4VPfbUSQ28ioJ0r-8S4Ef_NM8pd7BH43srmGAqJyTxXTlaBp1uu4G0ASwRR_9XG2-i2ENG3FwMY-U5VuZza3apl2DEFlc4I40KYlKMeXYbb1S2fmcNYb0deuFLSq-cLQhf_F3iMvHE9rp7RZa_ZGkKNX4am8KT6NHcIWncBgR5x2jgi4epxQCiAsjKjVehlJtQ-Z4',
  },
  movies: [
    {
      id: '1',
      title: 'Inception',
      genres: 'Sci-Fi, Thriller',
      year: 2010,
      poster: 'https://lh3.googleusercontent.com/aida-public/AB6AXuDUO4x1otakKa7yPxbCBAopRxeFSsQJV3sy39i_8_C9v-h4yb93ROpZNXyuB5_TBKdjttgxbsXFjhKQ7RbNvTzxnhGVav23MAus8YOOq5zyThFW2KbbZGcvjGSiILZ4Ck346NCT4ermU670SGlYIwAgYY8PI1hFfmnqZvOZ_BQ89XiUc_v5zlgOd0y1JpBNL-bcf-RfLh0uv6RYfODWuJxKruDv5HSkDdMlZVaqlqmBo36kVS3HNwojDKDJbaZ3Y-n9pOkndwqvoio',
    },
    {
      id: '2',
      title: 'Blade Runner 2049',
      genres: 'Sci-Fi, Neo-noir',
      year: 2017,
      poster: 'https://lh3.googleusercontent.com/aida-public/AB6AXuD-97_sLXltF6Q6MENFyKU_TWFNaHS9GPZSybvvyjaZ588PzuCYUOpEaFVQwpqcT3hPid30S2e-XHxLg84bAyKcID5WU_C59kJLNHmieNVs-dceAWUpL_BclLXqUPlox_oLM0XDywavEcmbRmJS0adm7src-3D-czYyAHq1PpV3zBoFYVRJLp4VWi9WrK0Xe2BAetFYDTj-d5-ocqjJOOihj2WxcL1XrUOokqFnPDFwqnxxNbi8nG3Bs0lwqR6XUhX_dnfcnFmoDfw',
    },
    {
      id: '3',
      title: 'Arrival',
      genres: 'Sci-Fi, Drama',
      year: 2016,
      poster: 'https://lh3.googleusercontent.com/aida-public/AB6AXuC2XNHaJcz_i3W4FQ_fK3gIFry9ooA81vLZlZuvChSSUqYgZqpil38ERcOomSS22K-54FIMgIGIOkcko-SrLk41OeyJRbAwdgdgTCB6HxIn6zJE8c4fzlIdTm6j0edXkb-OP7dOFjLGf0XiDesWXB4y1XNW171f9ZlhmuGWXihdosT_iidjkiT5w4G66lMi-vi3YKKg9xq1sgouadfP3O2sdGmDn3eCQwgD-JCNNCTQexigOsloIlcNZT-1_rB-UIzmXsr0VyYX3tA',
    },
    {
      id: '4',
      title: 'Interstellar',
      genres: 'Sci-Fi, Adventure',
      year: 2014,
      poster: 'https://lh3.googleusercontent.com/aida-public/AB6AXuCrBk8vkreZ0RJW2prJ5Vy4gFy7ZyDkG6lAlDOq5XkvOn-v5jbCl4lb39EyhFQVT5l8ccGVQa-hZehcT35H00aulANN7zyXnImYMCCAJpYqK4LMwWg7KRW4vHVSpAganCcGPNJ3tE-5sgEDshReR0sXS1TSxJFVwKtRvYkujbsQ6eKrMuxL0PHfRPct_oYWB4ep3Rf-W-5U1Yc-49OoEN82NtZ59iWYZylUUlrTbV84rNV0fqbfOFQZ1fmUfO4cT1_7_p29MYNVExk',
    },
  ],
};

const MovieListItem = ({ movie }: { movie: typeof watchlistData.movies[0] }) => (
  <Pressable style={styles.movieItem}>
    <Image source={{ uri: movie.poster }} style={styles.moviePoster} />
    <View style={styles.movieInfo}>
      <Text style={styles.movieTitle}>{movie.title}</Text>
      <Text style={styles.movieGenres}>{movie.genres}</Text>
      <Text style={styles.movieYear}>{movie.year}</Text>
    </View>
  </Pressable>
);

export default function WatchlistDetailsScreen() {
  const params = useLocalSearchParams();
  
  return (
    <View style={styles.container}>
      <ScrollView 
        style={styles.scrollView} 
        showsVerticalScrollIndicator={false}
        bounces={false}
      >
        {/* Hero Section */}
        <View style={styles.heroSection}>
          <ImageBackground
            source={{ uri: watchlistData.coverImage }}
            style={styles.heroImage}
            resizeMode="cover"
          >
            <LinearGradient
              colors={['transparent', 'rgba(10, 10, 11, 0.7)', BACKGROUND_DARK]}
              style={styles.heroGradient}
            />
            
            {/* Back Button */}
            <SafeAreaView edges={['top']} style={styles.headerContainer}>
              <Pressable style={styles.backButton} onPress={() => router.back()}>
                <MaterialIcons name="arrow-back-ios" size={22} color={TEXT_PRIMARY} />
              </Pressable>
            </SafeAreaView>

            {/* Creator & Stats */}
            <View style={styles.heroContent}>
              {/* Creator Info */}
              <View style={styles.creatorContainer}>
                <Image source={{ uri: watchlistData.creator.avatar }} style={styles.creatorAvatar} />
                <Text style={styles.creatorName}>{watchlistData.creator.name}</Text>
              </View>

              {/* Stats & Actions */}
              <View style={styles.statsContainer}>
                <View style={styles.likesContainer}>
                  <MaterialIcons name="favorite" size={24} color={PRIMARY_COLOR} />
                  <Text style={styles.likesText}>{watchlistData.likes}</Text>
                </View>
                <View style={styles.actionsContainer}>
                  <Pressable style={styles.actionButton}>
                    <MaterialIcons name="bookmark-add" size={24} color={TEXT_PRIMARY} />
                  </Pressable>
                  <Pressable style={styles.actionButton}>
                    <MaterialIcons name="ios-share" size={24} color={TEXT_PRIMARY} />
                  </Pressable>
                </View>
              </View>
            </View>
          </ImageBackground>
        </View>

        {/* Watchlist Info */}
        <View style={styles.infoSection}>
          <Text style={styles.watchlistTitle}>{watchlistData.title}</Text>
          <Text style={styles.watchlistDescription}>{watchlistData.description}</Text>
        </View>

        {/* Movies List */}
        <View style={styles.moviesSection}>
          {watchlistData.movies.map((movie) => (
            <MovieListItem key={movie.id} movie={movie} />
          ))}
        </View>
      </ScrollView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: BACKGROUND_DARK,
  },
  scrollView: {
    flex: 1,
  },
  heroSection: {
    width: '100%',
    height: 400,
  },
  heroImage: {
    flex: 1,
    justifyContent: 'space-between',
  },
  heroGradient: {
    ...StyleSheet.absoluteFillObject,
  },
  headerContainer: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    zIndex: 10,
    paddingHorizontal: 16,
    paddingTop: 8,
  },
  backButton: {
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: 'rgba(0, 0, 0, 0.3)',
    alignItems: 'center',
    justifyContent: 'center',
  },
  heroContent: {
    position: 'absolute',
    bottom: 0,
    left: 0,
    right: 0,
    padding: 16,
    gap: 12,
  },
  creatorContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
    backgroundColor: 'rgba(0, 0, 0, 0.4)',
    padding: 12,
    borderRadius: 12,
    backdropFilter: 'blur(10px)',
  },
  creatorAvatar: {
    width: 40,
    height: 40,
    borderRadius: 20,
  },
  creatorName: {
    flex: 1,
    fontSize: 16,
    fontWeight: '500',
    color: TEXT_PRIMARY,
  },
  statsContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    backgroundColor: 'rgba(0, 0, 0, 0.4)',
    padding: 4,
    borderRadius: 12,
  },
  likesContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    paddingHorizontal: 12,
    paddingVertical: 8,
  },
  likesText: {
    fontSize: 14,
    fontWeight: 'bold',
    color: TEXT_PRIMARY,
  },
  actionsContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 16,
    paddingHorizontal: 12,
    paddingVertical: 8,
  },
  actionButton: {
    padding: 4,
  },
  infoSection: {
    paddingHorizontal: 16,
    paddingTop: 24,
    paddingBottom: 16,
  },
  watchlistTitle: {
    fontSize: 32,
    fontWeight: 'bold',
    color: TEXT_PRIMARY,
    letterSpacing: -0.5,
    marginBottom: 8,
  },
  watchlistDescription: {
    fontSize: 16,
    color: TEXT_SECONDARY,
    lineHeight: 24,
  },
  moviesSection: {
    paddingHorizontal: 16,
    paddingBottom: 32,
    gap: 16,
  },
  movieItem: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 16,
    backgroundColor: SURFACE_DARK,
    padding: 12,
    borderRadius: 12,
  },
  moviePoster: {
    width: 75,
    height: 112,
    borderRadius: 8,
    backgroundColor: '#27272a',
  },
  movieInfo: {
    flex: 1,
    justifyContent: 'center',
  },
  movieTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: TEXT_PRIMARY,
  },
  movieGenres: {
    fontSize: 14,
    color: TEXT_SECONDARY,
    marginTop: 4,
  },
  movieYear: {
    fontSize: 14,
    color: '#6B7280',
    marginTop: 4,
  },
});
