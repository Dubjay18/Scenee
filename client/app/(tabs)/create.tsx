import {
  StyleSheet,
  View,
  ScrollView,
  TextInput,
  Pressable,
  Image,
} from 'react-native';
import { Text } from '../../components/Themed';
import { MaterialIcons } from '@expo/vector-icons';
import { SafeAreaView } from 'react-native-safe-area-context';
import { router } from 'expo-router';
import { useState } from 'react';
import { LinearGradient } from 'expo-linear-gradient';

const PRIMARY_COLOR = '#A855F7';
const BACKGROUND_DARK = '#0A0A0A';
const TEXT_PRIMARY = '#E0E0E0';
const TEXT_SECONDARY = '#A0A0A0';
const BORDER_COLOR = 'rgba(255, 255, 255, 0.1)';
const INPUT_BG = 'rgba(255, 255, 255, 0.05)';

type PrivacyType = 'public' | 'private';

// Mock data for added movies
const initialMovies = [
  {
    id: '1',
    title: 'Blade Runner 2049',
    poster: 'https://lh3.googleusercontent.com/aida-public/AB6AXuD9Ha9H_6Y9UiI0dI4e1VggZTZ7HNulpEcF8DhxV-FYqUvdgLRgh7vMntHeTiYVhaqWOvgQD7VHt2HIJlh-eL-08Rlo6rAdqyJpRihsbrTIP8ezGnNSvBoMkoVlJFHN4J9Oo5_bhS_TgSNun8xm3j4v80O3Nqo_YQPNbjcwM4EiUsI9UBxfIm9v65duYfGM4ySDZj4gystx8eMV4A4wYCAWRoZW9FwVvcqqxPbn_Agj0uSchsyQOf-Ohrf_RKmBUe05WM2uTsOZT_I',
  },
  {
    id: '2',
    title: 'Dune',
    poster: 'https://lh3.googleusercontent.com/aida-public/AB6AXuBWpcrYNv2M5J7dq6XaMrTGNMwJHi7nyoiyl3GIZ-Woyzl6feaxAXZJS74Pg0v8Xh_5goKYyLatUjt_H9Ei3UT_dIsWqniX_-LhMiODwo_lXiSbUiI7ocWk3J2bEn-yDoPyFzYmLihGrPyoxYTzVAXACpuZtORY5_LxM_rTqtEKrP0d0Wgeos7m2imWWI2JTb91c6hKb9Dqf4iMR4btrfiQjz_w8iQMP6DbP8iKCNJMF689DyUuUanUcnVQaYpietWE2y0KTulz2ws',
  },
];

const MoviePosterCard = ({ 
  movie, 
  onRemove 
}: { 
  movie: typeof initialMovies[0]; 
  onRemove: () => void;
}) => (
  <View style={styles.moviePosterCard}>
    <Image source={{ uri: movie.poster }} style={styles.moviePosterImage} />
    <Pressable style={styles.removeMovieButton} onPress={onRemove}>
      <MaterialIcons name="close" size={14} color="#FFFFFF" />
    </Pressable>
  </View>
);

const EmptyMoviesPlaceholder = () => (
  <View style={styles.emptyMoviesContainer}>
    <MaterialIcons name="movie" size={48} color="rgba(255, 255, 255, 0.2)" />
    <Text style={styles.emptyMoviesText}>Your watchlist is empty... for now!</Text>
  </View>
);

export default function CreateWatchlistScreen() {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [privacy, setPrivacy] = useState<PrivacyType>('public');
  const [movieSearch, setMovieSearch] = useState('');
  const [movies, setMovies] = useState(initialMovies);

  const handleRemoveMovie = (id: string) => {
    setMovies(movies.filter((m) => m.id !== id));
  };

  const handleCreateWatchlist = () => {
    // TODO: Implement create watchlist logic
    router.navigate('/(tabs)');
  };

  return (
    <View style={styles.container}>
      <SafeAreaView edges={['top']} style={styles.safeArea}>
        {/* Header */}
        <View style={styles.header}>
          <Pressable style={styles.backButton} onPress={() => router.navigate('/(tabs)')}>
            <MaterialIcons name="arrow-back" size={28} color="rgba(255, 255, 255, 0.9)" />
          </Pressable>
          <Text style={styles.headerTitle}>Create New Watchlist</Text>
          <View style={styles.headerSpacer} />
        </View>

        {/* Content */}
        <ScrollView
          style={styles.scrollView}
          contentContainerStyle={styles.scrollContent}
          showsVerticalScrollIndicator={false}
        >
          {/* Cover Image Upload */}
          <Pressable style={styles.coverUpload}>
            <MaterialIcons name="add-photo-alternate" size={36} color={PRIMARY_COLOR} />
            <View style={styles.coverUploadTextContainer}>
              <Text style={styles.coverUploadTitle}>Add a Cover Image</Text>
              <Text style={styles.coverUploadSubtitle}>Tap to upload a cover for your list.</Text>
            </View>
          </Pressable>

          {/* Form Fields */}
          <View style={styles.formSection}>
            {/* Title Input */}
            <View style={styles.inputGroup}>
              <Text style={styles.inputLabel}>Watchlist Title</Text>
              <TextInput
                style={styles.textInput}
                placeholder="My Sci-Fi Favorites"
                placeholderTextColor="rgba(160, 160, 160, 0.5)"
                value={title}
                onChangeText={setTitle}
              />
            </View>

            {/* Description Input */}
            <View style={styles.inputGroup}>
              <Text style={styles.inputLabel}>Description</Text>
              <TextInput
                style={[styles.textInput, styles.textArea]}
                placeholder="Add a cool description..."
                placeholderTextColor="rgba(160, 160, 160, 0.5)"
                value={description}
                onChangeText={setDescription}
                multiline
                numberOfLines={4}
                textAlignVertical="top"
              />
            </View>
          </View>

          {/* Privacy Settings */}
          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Privacy Settings</Text>
            <View style={styles.privacyToggle}>
              <Pressable
                style={[styles.privacyButton, privacy === 'public' && styles.privacyButtonActive]}
                onPress={() => setPrivacy('public')}
              >
                <MaterialIcons
                  name="public"
                  size={18}
                  color={privacy === 'public' ? BACKGROUND_DARK : TEXT_SECONDARY}
                />
                <Text
                  style={[
                    styles.privacyButtonText,
                    privacy === 'public' && styles.privacyButtonTextActive,
                  ]}
                >
                  Public
                </Text>
              </Pressable>
              <Pressable
                style={[styles.privacyButton, privacy === 'private' && styles.privacyButtonActive]}
                onPress={() => setPrivacy('private')}
              >
                <MaterialIcons
                  name="lock"
                  size={18}
                  color={privacy === 'private' ? BACKGROUND_DARK : TEXT_SECONDARY}
                />
                <Text
                  style={[
                    styles.privacyButtonText,
                    privacy === 'private' && styles.privacyButtonTextActive,
                  ]}
                >
                  Private
                </Text>
              </Pressable>
            </View>
          </View>

          {/* Add Movies Section */}
          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Add Movies</Text>
            <View style={styles.searchInputContainer}>
              <MaterialIcons
                name="search"
                size={22}
                color="rgba(160, 160, 160, 0.5)"
                style={styles.searchIcon}
              />
              <TextInput
                style={styles.searchInput}
                placeholder="Search for movies to add..."
                placeholderTextColor="rgba(160, 160, 160, 0.5)"
                value={movieSearch}
                onChangeText={setMovieSearch}
              />
            </View>
          </View>

          {/* Movies Grid */}
          <View style={styles.moviesGrid}>
            {movies.map((movie) => (
              <MoviePosterCard
                key={movie.id}
                movie={movie}
                onRemove={() => handleRemoveMovie(movie.id)}
              />
            ))}
            {movies.length === 0 && <EmptyMoviesPlaceholder />}
          </View>
        </ScrollView>

        {/* Create Button */}
        <View style={styles.bottomContainer}>
          <LinearGradient
            colors={['transparent', BACKGROUND_DARK]}
            style={styles.bottomGradient}
          />
          <Pressable style={styles.createButton} onPress={handleCreateWatchlist}>
            <Text style={styles.createButtonText}>Create Watchlist</Text>
          </Pressable>
        </View>
      </SafeAreaView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: BACKGROUND_DARK,
  },
  safeArea: {
    flex: 1,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: 16,
    paddingVertical: 8,
    backgroundColor: 'rgba(10, 10, 10, 0.8)',
  },
  backButton: {
    width: 48,
    height: 48,
    alignItems: 'flex-start',
    justifyContent: 'center',
  },
  headerTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: TEXT_PRIMARY,
    flex: 1,
    textAlign: 'center',
  },
  headerSpacer: {
    width: 48,
  },
  scrollView: {
    flex: 1,
  },
  scrollContent: {
    padding: 16,
    paddingBottom: 100,
  },
  coverUpload: {
    alignItems: 'center',
    justifyContent: 'center',
    gap: 24,
    paddingVertical: 56,
    paddingHorizontal: 24,
    borderWidth: 2,
    borderStyle: 'dashed',
    borderColor: 'rgba(168, 85, 247, 0.3)',
    borderRadius: 12,
    marginBottom: 16,
  },
  coverUploadTextContainer: {
    alignItems: 'center',
    gap: 4,
  },
  coverUploadTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: TEXT_PRIMARY,
    textAlign: 'center',
  },
  coverUploadSubtitle: {
    fontSize: 14,
    color: TEXT_SECONDARY,
    textAlign: 'center',
  },
  formSection: {
    gap: 16,
    marginBottom: 16,
  },
  inputGroup: {
    gap: 8,
  },
  inputLabel: {
    fontSize: 16,
    fontWeight: '500',
    color: TEXT_PRIMARY,
  },
  textInput: {
    height: 56,
    backgroundColor: INPUT_BG,
    borderWidth: 1,
    borderColor: BORDER_COLOR,
    borderRadius: 8,
    paddingHorizontal: 16,
    fontSize: 16,
    color: '#FFFFFF',
  },
  textArea: {
    height: 144,
    paddingTop: 16,
    paddingBottom: 16,
  },
  section: {
    marginBottom: 16,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: TEXT_PRIMARY,
    marginBottom: 8,
    paddingHorizontal: 4,
  },
  privacyToggle: {
    flexDirection: 'row',
    backgroundColor: INPUT_BG,
    borderRadius: 12,
    padding: 4,
  },
  privacyButton: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: 8,
    paddingVertical: 12,
    borderRadius: 8,
  },
  privacyButtonActive: {
    backgroundColor: PRIMARY_COLOR,
  },
  privacyButtonText: {
    fontSize: 14,
    fontWeight: 'bold',
    color: TEXT_SECONDARY,
  },
  privacyButtonTextActive: {
    color: BACKGROUND_DARK,
  },
  searchInputContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    height: 56,
    backgroundColor: INPUT_BG,
    borderWidth: 1,
    borderColor: BORDER_COLOR,
    borderRadius: 8,
  },
  searchIcon: {
    marginLeft: 16,
  },
  searchInput: {
    flex: 1,
    height: '100%',
    paddingHorizontal: 12,
    fontSize: 16,
    color: '#FFFFFF',
  },
  moviesGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 12,
  },
  moviePosterCard: {
    width: '30%',
    aspectRatio: 2 / 3,
    borderRadius: 8,
    overflow: 'hidden',
    position: 'relative',
  },
  moviePosterImage: {
    width: '100%',
    height: '100%',
  },
  removeMovieButton: {
    position: 'absolute',
    top: 4,
    right: 4,
    width: 24,
    height: 24,
    borderRadius: 12,
    backgroundColor: 'rgba(0, 0, 0, 0.6)',
    alignItems: 'center',
    justifyContent: 'center',
  },
  emptyMoviesContainer: {
    width: '100%',
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: 40,
    borderWidth: 2,
    borderStyle: 'dashed',
    borderColor: BORDER_COLOR,
    borderRadius: 12,
    gap: 16,
  },
  emptyMoviesText: {
    fontSize: 14,
    fontWeight: '500',
    color: TEXT_SECONDARY,
    textAlign: 'center',
  },
  bottomContainer: {
    position: 'absolute',
    bottom: 0,
    left: 0,
    right: 0,
    paddingHorizontal: 16,
    paddingBottom: 24,
    paddingTop: 16,
  },
  bottomGradient: {
    ...StyleSheet.absoluteFillObject,
  },
  createButton: {
    height: 56,
    backgroundColor: PRIMARY_COLOR,
    borderRadius: 9999,
    alignItems: 'center',
    justifyContent: 'center',
    shadowColor: PRIMARY_COLOR,
    shadowOffset: { width: 0, height: 0 },
    shadowOpacity: 0.5,
    shadowRadius: 20,
    elevation: 8,
  },
  createButtonText: {
    fontSize: 16,
    fontWeight: 'bold',
    color: '#FFFFFF',
  },
});
