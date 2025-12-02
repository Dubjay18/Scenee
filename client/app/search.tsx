import {
  StyleSheet,
  View,
  TextInput,
  Pressable,
  ScrollView,
  Image,
  ImageBackground,
} from 'react-native';
import { Text } from '../components/Themed';
import { MaterialIcons } from '@expo/vector-icons';
import { SafeAreaView } from 'react-native-safe-area-context';
import { router } from 'expo-router';
import { useState } from 'react';

const PRIMARY_COLOR = '#C039FF';
const BACKGROUND_DARK = '#0A0A0A';
const TEXT_PRIMARY = '#FFFFFF';
const TEXT_SECONDARY = '#9CA3AF';

type TabType = 'movies' | 'watchlists' | 'users';

const recentSearches = [
  { id: '1', text: 'Blade Runner 2049' },
  { id: '2', text: 'Dune' },
  { id: '3', text: 'Cyberpunk' },
];

const categories = [
  {
    id: '1',
    name: 'Sci-Fi',
    image: 'https://lh3.googleusercontent.com/aida-public/AB6AXuDlNkRSiR6opvsnag2ndUicnhdrsn3pDpwVSEkvAA0YDRhKrYvYde8YnzjFnTaKcLIWv6m9t1Dox7VqdrIgPSgF4ojIr2GA4MXIVTKri8hf18pwpjD-pLLwyOsZitamifsBLYrcwIOGuZ4AXsNvz3OADSjpXmbqgajELxkQXvULWAwA7nT5aRF7FrWxyScEkT16lA8VGDINxdCGgJmjaJGrLwHKHzTJ8XT32mVj3eM0okjRX6Rf_F8roaF1suOlV8KPBOkzwOg-_3M',
  },
  {
    id: '2',
    name: 'Horror',
    image: 'https://lh3.googleusercontent.com/aida-public/AB6AXuC86Nz_pgZgTwWD1o8J2j6Ablt5jQs5RTHiKWFCGfyzf0rGsSssNgFWJy3qW3OxrPUjuderD9rO63RYmNCs6VdVTiaUp_yO0pCJBOCiG7pnUcx80W1OIQjoH_YriYjLhe6eAnqjhCbsuEvT6btJ7eU31qqokY3Xz7HuOHxmxyvUVQ2WfaxIgW6RC4NwKmqayne-1DVWQt-O92zrrHgPFTJTB3vf0kHWw-hkPaKcREDX-v96o2FKje6DdaI7hcFvmwzVp0caVVkXT2I',
  },
  {
    id: '3',
    name: "80s Classics",
    image: 'https://lh3.googleusercontent.com/aida-public/AB6AXuBUvZVMeGlnPQsak7XMABeqo32f_7212nIRt1C80J8FcspxBS1sHqD_CWPKUmyijc2Ul8wydhX5q-pbfG1iBu5c1Dr1kHpDZq53lRJ9vYJfI0qik9-7i64pSNTZgVWEMTwH6IfZTlaHRyJQ986Ji1wGWfwo0wLwLFHljr0rU3Nqr3nPDCtSehteLqWcM6AUXnAhNgU2Jqv2a8gUutUbv0brHtN538DCjgBBKsbnjm_PMJ1p10_6HwIMhiaJGqFn5UJEiXyRPs78-yc',
  },
  {
    id: '4',
    name: 'Indie Gems',
    image: 'https://lh3.googleusercontent.com/aida-public/AB6AXuD0QHhosGRPE291sWqo8Fg1Uws3GH9dHiZTwmBdL1CpzNuPAuxXXk-YhGdc-kU7dQBh6gssXODaBOqobd3T0MWUVBO2T3xKwNcQzARE-OKTSJo7XDu1sZrxxeB2rWs94DWo0CMqyyQtIuFZSgsmZZhVoZaXhTZxP50TXX9YY-6e-raRNcodHXYmr2KzSQ0SgYygkK26iXOFb5pN4qjmzaNpjGRKkVeHMI5ITBg-lcbke4PjqXe1WRIPhOL2pVF9wzAEiVhk1BIFuxc',
  },
];

const RecentSearchChip = ({ text, onRemove }: { text: string; onRemove: () => void }) => (
  <View style={styles.recentChip}>
    <Text style={styles.recentChipText}>{text}</Text>
    <Pressable onPress={onRemove}>
      <MaterialIcons name="close" size={18} color={TEXT_SECONDARY} />
    </Pressable>
  </View>
);

const CategoryCard = ({ name, image }: { name: string; image: string }) => (
  <Pressable style={styles.categoryCard}>
    <ImageBackground source={{ uri: image }} style={styles.categoryImage} imageStyle={styles.categoryImageStyle}>
      <View style={styles.categoryOverlay} />
      <Text style={styles.categoryName}>{name}</Text>
    </ImageBackground>
  </Pressable>
);

const TabButton = ({ 
  label, 
  isActive, 
  onPress 
}: { 
  label: string; 
  isActive: boolean; 
  onPress: () => void;
}) => (
  <Pressable style={[styles.tabButton, isActive && styles.tabButtonActive]} onPress={onPress}>
    <Text style={[styles.tabButtonText, isActive && styles.tabButtonTextActive]}>
      {label}
    </Text>
  </Pressable>
);

const EmptyState = ({ query }: { query: string }) => (
  <View style={styles.emptyState}>
    <MaterialIcons name="movie" size={72} color={PRIMARY_COLOR} />
    <Text style={styles.emptyStateTitle}>Nothing found for '{query}'</Text>
    <Text style={styles.emptyStateSubtitle}>
      Try searching for another movie, user, or watchlist.
    </Text>
  </View>
);

export default function SearchScreen() {
  const [searchQuery, setSearchQuery] = useState('');
  const [activeTab, setActiveTab] = useState<TabType>('movies');
  const [searches, setSearches] = useState(recentSearches);

  const handleRemoveSearch = (id: string) => {
    setSearches(searches.filter((s) => s.id !== id));
  };

  const handleClearSearch = () => {
    setSearchQuery('');
  };

  const showEmptyState = searchQuery.length > 0 && searchQuery !== 'Search for a movie, user, or watchlist...';

  return (
    <View style={styles.container}>
      <SafeAreaView edges={['top']} style={styles.safeArea}>
        {/* Header */}
        <View style={styles.header}>
          <Pressable style={styles.backButton} onPress={() => router.back()}>
            <MaterialIcons name="arrow-back" size={28} color={TEXT_PRIMARY} />
          </Pressable>
          <Text style={styles.headerTitle}>Search</Text>
          <View style={styles.headerSpacer} />
        </View>

        {/* Search Input */}
        <View style={styles.searchContainer}>
          <View style={styles.searchInputWrapper}>
            <View style={styles.searchIconContainer}>
              <MaterialIcons name="search" size={24} color={PRIMARY_COLOR} />
            </View>
            <TextInput
              style={styles.searchInput}
              placeholder="Search for a movie, user, or watchlist..."
              placeholderTextColor={TEXT_SECONDARY}
              value={searchQuery}
              onChangeText={setSearchQuery}
              autoFocus
            />
            {searchQuery.length > 0 && (
              <Pressable style={styles.clearButton} onPress={handleClearSearch}>
                <MaterialIcons name="close" size={22} color={TEXT_SECONDARY} />
              </Pressable>
            )}
          </View>
        </View>

        {/* Tabs */}
        <View style={styles.tabsContainer}>
          <View style={styles.tabsWrapper}>
            <TabButton
              label="Movies"
              isActive={activeTab === 'movies'}
              onPress={() => setActiveTab('movies')}
            />
            <TabButton
              label="Watchlists"
              isActive={activeTab === 'watchlists'}
              onPress={() => setActiveTab('watchlists')}
            />
            <TabButton
              label="Users"
              isActive={activeTab === 'users'}
              onPress={() => setActiveTab('users')}
            />
          </View>
        </View>

        {/* Content */}
        <ScrollView
          style={styles.scrollView}
          contentContainerStyle={styles.scrollContent}
          showsVerticalScrollIndicator={false}
        >
          {showEmptyState ? (
            <EmptyState query={searchQuery} />
          ) : (
            <>
              {/* Recent Searches */}
              {searches.length > 0 && (
                <View style={styles.section}>
                  <Text style={styles.sectionTitle}>Recent Searches</Text>
                  <ScrollView
                    horizontal
                    showsHorizontalScrollIndicator={false}
                    contentContainerStyle={styles.recentSearchesContainer}
                  >
                    {searches.map((search) => (
                      <RecentSearchChip
                        key={search.id}
                        text={search.text}
                        onRemove={() => handleRemoveSearch(search.id)}
                      />
                    ))}
                  </ScrollView>
                </View>
              )}

              {/* Popular Categories */}
              <View style={styles.section}>
                <Text style={styles.sectionTitle}>Popular Categories</Text>
                <View style={styles.categoriesGrid}>
                  {categories.map((category) => (
                    <CategoryCard
                      key={category.id}
                      name={category.name}
                      image={category.image}
                    />
                  ))}
                </View>
              </View>
            </>
          )}
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
    paddingVertical: 8,
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
  searchContainer: {
    paddingHorizontal: 16,
    paddingVertical: 12,
  },
  searchInputWrapper: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: 'rgba(255, 255, 255, 0.05)',
    borderRadius: 12,
    height: 56,
  },
  searchIconContainer: {
    paddingLeft: 16,
    paddingRight: 8,
  },
  searchInput: {
    flex: 1,
    height: '100%',
    color: TEXT_PRIMARY,
    fontSize: 16,
  },
  clearButton: {
    padding: 16,
  },
  tabsContainer: {
    paddingHorizontal: 16,
    paddingBottom: 12,
  },
  tabsWrapper: {
    flexDirection: 'row',
    borderBottomWidth: 1,
    borderBottomColor: 'rgba(255, 255, 255, 0.1)',
  },
  tabButton: {
    flex: 1,
    paddingVertical: 16,
    alignItems: 'center',
    borderBottomWidth: 3,
    borderBottomColor: 'transparent',
  },
  tabButtonActive: {
    borderBottomColor: PRIMARY_COLOR,
  },
  tabButtonText: {
    fontSize: 14,
    fontWeight: 'bold',
    color: TEXT_SECONDARY,
  },
  tabButtonTextActive: {
    color: TEXT_PRIMARY,
  },
  scrollView: {
    flex: 1,
  },
  scrollContent: {
    paddingBottom: 32,
  },
  section: {
    paddingTop: 16,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: TEXT_PRIMARY,
    paddingHorizontal: 16,
    paddingBottom: 8,
  },
  recentSearchesContainer: {
    paddingHorizontal: 16,
    paddingTop: 8,
    gap: 8,
  },
  recentChip: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    backgroundColor: 'rgba(255, 255, 255, 0.1)',
    borderRadius: 9999,
    paddingLeft: 16,
    paddingRight: 12,
    height: 32,
    marginRight: 8,
  },
  recentChipText: {
    fontSize: 14,
    fontWeight: '500',
    color: TEXT_PRIMARY,
  },
  categoriesGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    paddingHorizontal: 16,
    paddingTop: 12,
    gap: 16,
  },
  categoryCard: {
    width: '47%',
    height: 96,
    borderRadius: 16,
    overflow: 'hidden',
  },
  categoryImage: {
    flex: 1,
    justifyContent: 'flex-end',
    padding: 16,
  },
  categoryImageStyle: {
    borderRadius: 16,
  },
  categoryOverlay: {
    ...StyleSheet.absoluteFillObject,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
  },
  categoryName: {
    fontSize: 18,
    fontWeight: 'bold',
    color: TEXT_PRIMARY,
    zIndex: 1,
  },
  emptyState: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    paddingHorizontal: 32,
    paddingTop: 64,
  },
  emptyStateTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    color: TEXT_PRIMARY,
    marginTop: 16,
    textAlign: 'center',
  },
  emptyStateSubtitle: {
    fontSize: 16,
    color: TEXT_SECONDARY,
    marginTop: 8,
    textAlign: 'center',
  },
});
