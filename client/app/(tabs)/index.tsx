import { StyleSheet, ScrollView, View, Image, Pressable } from 'react-native';
import { Text } from '../../components/Themed';
import { MaterialIcons } from '@expo/vector-icons';
import { SafeAreaView } from 'react-native-safe-area-context';
import { router } from 'expo-router';

const PRIMARY_COLOR = '#A855F7';
const BACKGROUND_DARK = '#0A0A0A';
const SURFACE_DARK = '#161616';
const TEXT_PRIMARY = '#EAEAEA';
const TEXT_SECONDARY = '#A9A9A9';
const ACCENT_PINK = '#f20da6';

// Mock data for watchlists
const trendingWatchlists = [
  {
    id: '1',
    title: 'Sci-Fi Classics',
    description: 'A curated list of mind-bending science fiction movies that defined the genre.',
    image: 'https://lh3.googleusercontent.com/aida-public/AB6AXuCbQ0XMSn0o3SE8AkGjcyPM7OHCLIyagS4PbYlj1IcCSyJbg-DD3gHM1XBwI_1oOrKoeYJ74cC8zlDAnwQ96MvFDvrTmNP5IHdftdnmG3lRctW0OQvxbsuY4GpNJR5jkLbGGbN_EmjbCW9SQhxc_5C5b4VW8WS0MpIjsey8hvCNpUa7dKLeuWmkEkWiCpMVJRewxVuQ-U7YCCc8INN4i541D6QGriyIijPpFCWups9F3Oddkq6iQVGbHV8LndAvCccJ4q7ztTA0W4Y',
    likes: '2.1k',
    movieCount: 15,
    saved: false,
  },
  {
    id: '2',
    title: 'Film Noir Favorites',
    description: 'Step into the shadows with these classic tales of crime, mystery, and fatalism.',
    image: 'https://lh3.googleusercontent.com/aida-public/AB6AXuA3gvy_ch1zuI231OUoSCID6KV-2GPsApHNPLHjBStMNl3rcFQPHNUfw4XQDBCjx1aIqQ70t5n9Qw3GsrXOU-ZnKIJUJsEicS3-5Sc-aDgo9hPnRkHaajuTw_8-8r0Et0vEfpijZAT9MQmqpwgViD1sOy7aUepzNfIam0Qa9Aqa6VbEGFR4K7KsQG9r00iElTgUDS_oI4CScJKRzER5n6enMJbZQA8zV5huw_MApjbGP9keLU5kSj3iBQ71u6MYfM-WpPPwbstojfE',
    likes: '1.8k',
    movieCount: 20,
    saved: true,
  },
  {
    id: '3',
    title: 'Animated Adventures',
    description: 'From hand-drawn masterpieces to modern CGI wonders.',
    image: 'https://lh3.googleusercontent.com/aida-public/AB6AXuAJibT841S4PCT_KWrUpz6cI30gd76v8dDGD0pQ7j18zoDYm23zfflc5L79JDQstdcPN2ncYskNCu7clInYkxWJRl9-kSdIWtTgfx-c_rQq9qc30P4xOLmGtRndvOAuJoJDOWDgrkI-rF6KJO2AHNC_YMD93ChISqNA8V8YT-Z48xEqQdVgb-0wpOU4slU7MXfbs1ve6VxajjLHWalP90YHpLA6_xnZkcxUDMxmzLp9JfdJrqB5eh12KrORvOPNhiPyLyUx6URYjnk',
    likes: '3.2k',
    movieCount: 25,
    saved: false,
  },
];

// Mock data for recommended movies
const recommendedMovies = [
  {
    id: '1',
    title: 'Tenet',
    genre: 'Sci-Fi',
    year: 2020,
    image: 'https://lh3.googleusercontent.com/aida-public/AB6AXuAtfTo8yCLOaCpjsVAMOkRkb8fFfNo7aMlNRqy7GNDuzpdwu2Oso-hUw9vKJ5-DnRowdv6haF6IAvHA-Y_z6_0AdvgJGkbODfoqQOQrz4aW-Sg_7hKSGcAg9ncrXBXP9BNhB4wNUDTQvLGpEu4ECdDXf-z1vpBOZ0v6DQnFqAhkU6ojaINWWL0ZYyK4tRWBDQdfFQIrDmimqh6PvQzdyZKgA3-O3G-oAwlZQtDpYSmnYD6cVZM7VR01qCKJtUugoxE_EFmDc0FTPvY',
    badge: 'Watched',
    badgeColor: PRIMARY_COLOR,
  },
  {
    id: '2',
    title: 'Interstellar',
    genre: 'Sci-Fi',
    year: 2014,
    image: 'https://lh3.googleusercontent.com/aida-public/AB6AXuAwBcOJk-Lq60ZzV3uyBYrgHts4OalxTo0GXstcYJ_BTBglbooc3ykePlK2NwJUZvVRe6EjsF3ZrxBTgPmhPWsKuChu-oCwOtAU-7jMvdcUmgVb6KF38XCdC9SOsuFD0kkS1Pjl-uVRBbm-PtRXDu6WTlxzCbkt4cb2zM9Pn038yxwafCL5csF2BfFZifW_fReurQFIgRUo0UIbjeOBnoHWg1RV_0xpxdBMDXSMMpHdrDhcOWzhWLL7s16cib6Z8T2S89i1OwiBe9I',
  },
  {
    id: '3',
    title: 'The Prestige',
    genre: 'Mystery',
    year: 2006,
    image: 'https://lh3.googleusercontent.com/aida-public/AB6AXuCaMQU3M_veZL0rs01T3RYmcxh7ja4az85LjuQ3F3Wd9HNPo2qYT2LSu5tZw0l51X8PSP3fzMfpZnBti6ylS7NIWs8HOue84w_-1oAZffGWS1iIj6PPm5vObt2gB_IuRNnLDGgle3tDtFp4ZeNRgCRUuQLqGMZRRcKSyfZwsr05CRImRNCRLQA2wtbA1gwj5_MamG4FaNjEaQrQZuFXU536Aze-xixX27Mi7LfPx7RYaTa0epijVUpcsCH5mWz-mhajkDv2kuDFVV4',
    badge: 'On Your List',
    badgeColor: ACCENT_PINK,
  },
  {
    id: '4',
    title: 'Shutter Island',
    genre: 'Thriller',
    year: 2010,
    image: 'https://lh3.googleusercontent.com/aida-public/AB6AXuBmRX5PaFQTHeHXCL23Z6JW04QOB1ITXbJD-JBnBkh0beaE3XyPh98Awfy2SCUjdvpST1X_MffpGF7eIrIXsZghloLKPMZi-FJPOuUeS9ol9HiLoGh72HHcebZm_FUAWeHeWLnSR7ANo6hk7Gf8FbAI3OoDziNp34EKXnzTiyhj-akJkfwe4wR7FuU9vbM1ONVJQf0YW6aLh1-_yR-ADfQ-9qgdydKRIFsIQ5UdxiEBSPYq1Ov_KCw7gLgqYgmaQKNxTHM911a31tk',
  },
  {
    id: '5',
    title: 'Arrival',
    genre: 'Sci-Fi',
    year: 2016,
    image: 'https://lh3.googleusercontent.com/aida-public/AB6AXuAoIFENeqbrX2Q7FgEFJ5ezEQCxE4ADGXYoaD4RcKu2g6EXEs_24Z_3yN6F1Z6V3hEy0dLyauhlaNtQnC4Ce8J6GEbWj3SFEJxTnWVsoaoe8gtSi6o7puaNgPt_eOZt9ScrQVW5OY5IaOXGrwqUmVEa3tYx504OuPpBPhDXpU8WW9Gq5J_0mIPB0iCcbu4N1UxI5G-4DGmEXfKSO6t9ky5BGjrFWQHXW1VKoprGFpP1iFaDn75fclsc0lpyczjS6Oj4bwr_7MM7ap0',
  },
];

const WatchlistCard = ({ item }: { item: typeof trendingWatchlists[0] }) => (
  <Pressable style={styles.watchlistCard} onPress={() => router.push(`/watchlist/${item.id}`)}>
    <Image source={{ uri: item.image }} style={styles.watchlistImage} />
    <View style={styles.watchlistContent}>
      <View style={styles.watchlistInfo}>
        <Text style={styles.watchlistTitle}>{item.title}</Text>
        <Text style={styles.watchlistDescription} numberOfLines={2}>
          {item.description}
        </Text>
        <View style={styles.watchlistStats}>
          <View style={styles.statItem}>
            <MaterialIcons name="favorite" size={16} color={PRIMARY_COLOR} />
            <Text style={styles.statText}>{item.likes}</Text>
          </View>
          <View style={styles.statItem}>
            <MaterialIcons name="movie" size={16} color={TEXT_SECONDARY} />
            <Text style={styles.statText}>{item.movieCount} Movies</Text>
          </View>
        </View>
      </View>
      <View style={styles.watchlistActions}>
        <Pressable style={[styles.saveButton, item.saved && styles.savedButton]}>
          <MaterialIcons
            name={item.saved ? 'bookmark' : 'bookmark-border'}
            size={18}
            color="#FFFFFF"
          />
          <Text style={styles.saveButtonText}>{item.saved ? 'Saved' : 'Save'}</Text>
        </Pressable>
        <Pressable style={styles.shareButton}>
          <MaterialIcons name="share" size={18} color="#FFFFFF" />
        </Pressable>
      </View>
    </View>
  </Pressable>
);

const MovieCard = ({ item }: { item: typeof recommendedMovies[0] }) => (
  <View style={styles.movieCard}>
    <View style={styles.moviePosterContainer}>
      <Image source={{ uri: item.image }} style={styles.moviePoster} />
      {item.badge && (
        <View style={[styles.movieBadge, { backgroundColor: `${item.badgeColor}CC` }]}>
          <Text style={styles.movieBadgeText}>{item.badge}</Text>
        </View>
      )}
    </View>
    <View style={styles.movieInfo}>
      <Text style={styles.movieTitle} numberOfLines={1}>{item.title}</Text>
      <Text style={styles.movieGenre} numberOfLines={1}>{item.genre}, {item.year}</Text>
    </View>
  </View>
);

export default function HomeScreen() {
  return (
    <View style={styles.container}>
      <SafeAreaView edges={['top']} style={styles.safeArea}>
        {/* Header */}
        <View style={styles.header}>
          <View style={styles.headerLeft}>
            <MaterialIcons name="play-circle-filled" size={32} color={PRIMARY_COLOR} />
            <Text style={styles.headerTitle}>Watchlyst</Text>
          </View>
          <View style={styles.headerRight}>
            <Pressable style={styles.headerButton} onPress={() => router.push('search')}>
              <MaterialIcons name="search" size={26} color={TEXT_PRIMARY} />
            </Pressable>
            <Pressable style={styles.headerButton}>
              <MaterialIcons name="notifications" size={26} color={TEXT_PRIMARY} />
            </Pressable>
          </View>
        </View>

        {/* Content */}
        <ScrollView
          style={styles.scrollView}
          contentContainerStyle={styles.scrollContent}
          showsVerticalScrollIndicator={false}
        >
          {/* Trending Watchlists Section */}
          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Trending Watchlists</Text>
            <ScrollView
              horizontal
              showsHorizontalScrollIndicator={false}
              contentContainerStyle={styles.horizontalScrollContent}
            >
              {trendingWatchlists.map((item) => (
                <WatchlistCard key={item.id} item={item} />
              ))}
            </ScrollView>
          </View>

          {/* Recommended Movies Section */}
          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Because You Watched Inception</Text>
            <ScrollView
              horizontal
              showsHorizontalScrollIndicator={false}
              contentContainerStyle={styles.horizontalScrollContent}
            >
              {recommendedMovies.map((item) => (
                <MovieCard key={item.id} item={item} />
              ))}
            </ScrollView>
          </View>
        </ScrollView>
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
    paddingVertical: 12,
    backgroundColor: 'rgba(10, 10, 10, 0.9)',
  },
  headerLeft: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  headerTitle: {
    fontSize: 22,
    fontWeight: 'bold',
    color: TEXT_PRIMARY,
    letterSpacing: -0.5,
  },
  headerRight: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
  },
  headerButton: {
    width: 40,
    height: 40,
    alignItems: 'center',
    justifyContent: 'center',
  },
  scrollView: {
    flex: 1,
  },
  scrollContent: {
    paddingBottom: 100,
  },
  section: {
    paddingTop: 20,
  },
  sectionTitle: {
    fontSize: 22,
    fontWeight: 'bold',
    color: TEXT_PRIMARY,
    paddingHorizontal: 16,
    paddingBottom: 12,
    letterSpacing: -0.3,
  },
  horizontalScrollContent: {
    paddingHorizontal: 16,
    gap: 16,
  },
  // Watchlist Card Styles
  watchlistCard: {
    width: 288,
    backgroundColor: SURFACE_DARK,
    borderRadius: 16,
    overflow: 'hidden',
  },
  watchlistImage: {
    width: '100%',
    aspectRatio: 16 / 9,
  },
  watchlistContent: {
    padding: 16,
    paddingTop: 0,
    marginTop: -8,
  },
  watchlistInfo: {
    gap: 8,
  },
  watchlistTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    color: TEXT_PRIMARY,
    marginTop: 16,
  },
  watchlistDescription: {
    fontSize: 14,
    color: TEXT_SECONDARY,
    lineHeight: 20,
  },
  watchlistStats: {
    flexDirection: 'row',
    gap: 16,
    marginTop: 8,
  },
  statItem: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 6,
  },
  statText: {
    fontSize: 14,
    color: TEXT_SECONDARY,
  },
  watchlistActions: {
    flexDirection: 'row',
    gap: 8,
    marginTop: 16,
  },
  saveButton: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: 8,
    height: 40,
    backgroundColor: PRIMARY_COLOR,
    borderRadius: 9999,
    shadowColor: PRIMARY_COLOR,
    shadowOffset: { width: 0, height: 0 },
    shadowOpacity: 0.6,
    shadowRadius: 5,
    elevation: 4,
  },
  savedButton: {
    backgroundColor: 'rgba(255, 255, 255, 0.1)',
    shadowOpacity: 0,
  },
  saveButtonText: {
    fontSize: 14,
    fontWeight: 'bold',
    color: '#FFFFFF',
  },
  shareButton: {
    width: 40,
    height: 40,
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: 'rgba(255, 255, 255, 0.1)',
    borderRadius: 9999,
  },
  // Movie Card Styles
  movieCard: {
    width: 160,
    gap: 8,
  },
  moviePosterContainer: {
    position: 'relative',
  },
  moviePoster: {
    width: '100%',
    aspectRatio: 2 / 3,
    borderRadius: 12,
    backgroundColor: SURFACE_DARK,
  },
  movieBadge: {
    position: 'absolute',
    top: 8,
    right: 8,
    paddingHorizontal: 8,
    paddingVertical: 4,
    borderRadius: 9999,
  },
  movieBadgeText: {
    fontSize: 12,
    fontWeight: 'bold',
    color: '#FFFFFF',
  },
  movieInfo: {
    gap: 2,
  },
  movieTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    color: TEXT_PRIMARY,
  },
  movieGenre: {
    fontSize: 14,
    color: TEXT_SECONDARY,
  },
});
