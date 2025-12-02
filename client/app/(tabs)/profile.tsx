import { StyleSheet, View, ScrollView, Image, Pressable } from 'react-native';
import { Text } from '../../components/Themed';
import { SafeAreaView } from 'react-native-safe-area-context';
import { MaterialIcons } from '@expo/vector-icons';
import { router } from 'expo-router';
import { useState } from 'react';

const BACKGROUND_DARK = '#0a0a0a';
const TEXT_PRIMARY = '#EAEAEA';
const TEXT_SECONDARY = '#A9A9A9';
const NEON_PURPLE = '#A855F7';
const SURFACE_DARK = '#111111';

type ViewMode = 'grid' | 'list';

// Mock user data
const userData = {
  username: '@cinemafanatic',
  bio: 'Lover of 80s sci-fi and Ghibli films.',
  avatar: 'https://lh3.googleusercontent.com/aida-public/AB6AXuD1mrdTLsOLwTZwOXizI7HvjgHL_5oZCBkZ50VTEhZk4Ugugr88kRs8_VzoS1uqcTqwCBB58Emzg7lZsnYQx2nJ613C5YS7PHbj1gVIy9kWagnam_TwoSlEv8q7YBjNsMYyGXBeHij0xohjBIDR6YI7p-ZPPgIgYB5WKi1ESloQQ6DtpB1JVh-e9H-nWu1qOrOQPkOZww3WkbjkCyFbd0HnL36MQIZMKhj3iMG1IRpGdLrcNY4n3SZNgDiSGXfEkVyBTF97pazHjOM',
  followers: '1.2k',
  following: '150',
};

// Mock watchlists data
const userWatchlists = [
  {
    id: '1',
    title: 'Mind-Bending Thrillers',
    movieCount: 12,
    isPublic: true,
    posters: [
      'https://lh3.googleusercontent.com/aida-public/AB6AXuAL1Y2B4K70URHJ_XqOuTN0VvxuNdCLMo5dvce113BahXRpdpxxwHhyI8mTNIi86RtvDF4P_cKQLS7hbyDJdJCEUdkFkVMJT1jbXLfYfciVR9W4s6tcAM8yawE_54Fs4Cf11MSumMRCSApC_dZ8SDr21VVBiW7Jylz_xKHtboq3ZJNAfUZyYEuimZqQIwcVBhgTPc03iMEPb1DkzGxpKWgSbxWI592nUfhxPgbjgqmDOPwiBzAn54HBkU6fAXYYmSZhNPUWOdAdAAI',
      'https://lh3.googleusercontent.com/aida-public/AB6AXuDXSN7RIFqu_IIDgPcUujWgrb_gw3I-5IEzBDy8cQaDgDf9l_kT2HqAKNYUOoi339pWS2jTb2DZVU8V_VBgSC9h3yUhp8xhwTWAwhxHeaiIZLZWS_BRseGdW-pqirkPm_nFUwhzp6kuMfTU8yTh-SfyulfjK2aKRzDGY1wIj4c2NWUnIZVKlq1MvUXOzG-LthUkp-08GSEpPa1QBublQm_3ityRnQ2GGcOSIEim_FeVVLvcYIR88YAbsUdH--CIXd8s5sXEruxappQ',
    ],
  },
  {
    id: '2',
    title: 'Cozy Sunday Movies',
    movieCount: 8,
    isPublic: false,
    posters: [
      'https://lh3.googleusercontent.com/aida-public/AB6AXuAfsotcLbjhNIOqfUX3SnoATexl-oSNb2gKTyeacEuYLaXBmsMbuTnwyX-Vw09ZR4rQl0b9dh5ZvwAfI5t4vSCY4F7V5xhQPz0qOjAuvYu3SDy-TmdUtsSY8ds7ZRJK85KE1Um_XfjrFkFfsOTCQcDpD1IkAEJMKEaz5pYIKvMU_3JEN9Jr9zIS9f66rL4Dfd0iR2Jm-fwAUJYnitzFW1oujYZwAiNoNc5iiTOhbUgirFeMvUueitF4MGtFkhtfh9ygShpCYvDwbaY',
      'https://lh3.googleusercontent.com/aida-public/AB6AXuDY-FcPLSsd2mrlY1Uk-z-1YbkK2xJc8OgQqWdlmobTOqdrLs54ZLKHaK4bf6amW-d55NG_FLIz_SI_W_TXeO0MWbxPTESAH555Yhwnm3aQr-6FSCjObLPihu4qxSAoIPWvW4MkswJp0BZ94Z94ueTxKzA1x80arN-89U4XxQeYYX-5ymiUyxHx8_ixy4fadp4yI8angNh2otZDiSUzIcV10_GRP3NopBI0ZT9tdJg9KE-KuJ8SrlvGTtVy6D1Ug5AVEFBP-gP25io',
    ],
  },
  {
    id: '3',
    title: 'Best of 2023',
    movieCount: 25,
    isPublic: true,
    posters: [
      'https://lh3.googleusercontent.com/aida-public/AB6AXuA6_VmvSrFIJlMgK1ignAxxcSoRh-YCCvUoOAHzxCEDNO65UcnvhZIwKg5yc_wNzkZqoEwYX1FWeWwq9QrLDoFYarXUGgEwCn78nHWIzQAcVWQk1e-R9QH78ljK3go_tW88sKDDaLF1g9bCc1kkJkdsglxvtw10sJaYeH_KnLofy4768j2TDIcXoMOdEIF7I79ignBs5d6SUGr_VDkOfEIy2ZmvetwnOf4LmvyhCBC17UYKxDlOirZWw4B0I88GDGHG0rngM4e__40',
      'https://lh3.googleusercontent.com/aida-public/AB6AXuC-Fbbe3ck9m4DkaKwXaL_3omjv6Y8oaDoNv1T4XyzesJno6t2tStEhPGOQ2dWT7bt3gbquSDBP_6MvOrfXaHNZxfNLt6O5z_KwGTG3a8pHN93Rt_s1oL2y4jNnZUCDClZJS6Z7Ojdq0T0x-gG0yCowGN3KZqpjuacVAbv25l_ovqH05eBZNcioPeU9IGwRGJTKx4EokZkCmxb-5FiKt87zN6Z7inTAUBO0wDI9q7YR76HjTsDtUm8Rfd3OQG-8aDSq5c7gr91nEkY',
    ],
  },
  {
    id: '4',
    title: 'Essential Horror',
    movieCount: 31,
    isPublic: true,
    posters: [
      'https://lh3.googleusercontent.com/aida-public/AB6AXuDXa1thaVCX2zkiJx-D5wYspwcutjQUA1yKTEU8-4cUXYGUPn-_p1mn-bSuFqfxjhUkhTjdkYa9t8fKAu1f1uZJoArcWgnWckFRVSk8Uf5nEZY5XnrOVkX0HaPW2irzu88BTOLSztcD0FHdyQteDXDhT0keRboY_q0yop1lQVLuKHIY4FrjXTzymFrBs-8S3HJXspbxOgU9F5Ck33WjRydlpjk8MUPTjvXMeDPLFnugEmmp1PfUmDMEqi5wexaadzBXu39KomzMbz8',
      'https://lh3.googleusercontent.com/aida-public/AB6AXuC0WWhCqC9GFDCFJJif14heM9Y9AjcaV8aB7WFA2BmGjSYtsxAp0Zv2hC4wSxUI004xKLT7jb_uqzazeYoW3_1QpcC3pAJtlgNm4wI0ZYaec2dWBQDXWnYCXs4uEUnveknGHsvXHv1mGybBp762HAm_RLUnR3CwWbRcg-aLoNW7EBrCgZTw5VfIap0JKkDHZ7tdjM5b6eW0_ZYTxmeGui8ZV0--VmLR79I83ytfgRGXNFtDanlNKdp_n4y17t04eKk08hGQHAHxRCk',
    ],
  },
];

const WatchlistCard = ({ watchlist }: { watchlist: typeof userWatchlists[0] }) => (
  <Pressable 
    style={styles.watchlistCard}
    onPress={() => router.push(`/watchlist/${watchlist.id}`)}
  >
    {/* Stacked Posters */}
    <View style={styles.postersContainer}>
      {watchlist.posters[1] && (
        <Image
          source={{ uri: watchlist.posters[1] }}
          style={styles.posterBack}
        />
      )}
      <Image
        source={{ uri: watchlist.posters[0] }}
        style={styles.posterFront}
      />
    </View>
    
    {/* Card Info */}
    <View style={styles.watchlistCardInfo}>
      <Text style={styles.watchlistCardTitle} numberOfLines={2}>
        {watchlist.title}
      </Text>
      <View style={styles.watchlistMeta}>
        <MaterialIcons
          name={watchlist.isPublic ? 'public' : 'lock'}
          size={14}
          color={TEXT_SECONDARY}
        />
        <Text style={styles.watchlistMetaText}>{watchlist.movieCount} Movies</Text>
      </View>
    </View>
  </Pressable>
);

export default function ProfileScreen() {
  const [viewMode, setViewMode] = useState<ViewMode>('grid');

  return (
    <View style={styles.container}>
      <SafeAreaView edges={['top']} style={styles.safeArea}>
        {/* Header */}
        <View style={styles.header}>
          <View style={styles.headerSpacer} />
          <Text style={styles.headerTitle}>Profile</Text>
          <Pressable style={styles.settingsButton}>
            <MaterialIcons name="settings" size={28} color={TEXT_PRIMARY} />
          </Pressable>
        </View>

        <ScrollView
          style={styles.scrollView}
          contentContainerStyle={styles.scrollContent}
          showsVerticalScrollIndicator={false}
        >
          {/* Profile Section */}
          <View style={styles.profileSection}>
            {/* Avatar */}
            <View style={styles.avatarContainer}>
              <Image source={{ uri: userData.avatar }} style={styles.avatar} />
            </View>

            {/* User Info */}
            <View style={styles.userInfo}>
              <Text style={styles.username}>{userData.username}</Text>
              <Text style={styles.bio}>{userData.bio}</Text>
            </View>
          </View>

          {/* Stats Section */}
          <View style={styles.statsSection}>
            <Pressable style={styles.statItem}>
              <Text style={styles.statNumber}>{userData.followers}</Text>
              <Text style={styles.statLabel}>Followers</Text>
            </Pressable>
            <Pressable style={styles.statItem}>
              <Text style={styles.statNumber}>{userData.following}</Text>
              <Text style={styles.statLabel}>Following</Text>
            </Pressable>
          </View>

          {/* Edit Profile Button */}
          <View style={styles.editButtonContainer}>
            <Pressable style={styles.editButton}>
              <Text style={styles.editButtonText}>Edit Profile</Text>
            </Pressable>
          </View>

          {/* My Watchlists Section */}
          <View style={styles.watchlistsHeader}>
            <Text style={styles.watchlistsTitle}>My Watchlists</Text>
            <View style={styles.viewModeButtons}>
              <Pressable
                style={[styles.viewModeButton, viewMode === 'grid' && styles.viewModeButtonActive]}
                onPress={() => setViewMode('grid')}
              >
                <MaterialIcons
                  name="grid-view"
                  size={24}
                  color={viewMode === 'grid' ? NEON_PURPLE : TEXT_SECONDARY}
                />
              </Pressable>
              <Pressable
                style={[styles.viewModeButton, viewMode === 'list' && styles.viewModeButtonActive]}
                onPress={() => setViewMode('list')}
              >
                <MaterialIcons
                  name="list"
                  size={24}
                  color={viewMode === 'list' ? NEON_PURPLE : TEXT_SECONDARY}
                />
              </Pressable>
            </View>
          </View>

          {/* Watchlists Grid */}
          <View style={styles.watchlistsGrid}>
            {userWatchlists.map((watchlist) => (
              <WatchlistCard key={watchlist.id} watchlist={watchlist} />
            ))}
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
    paddingVertical: 8,
  },
  headerSpacer: {
    width: 48,
  },
  headerTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: TEXT_PRIMARY,
    flex: 1,
    textAlign: 'center',
  },
  settingsButton: {
    width: 48,
    height: 48,
    alignItems: 'flex-end',
    justifyContent: 'center',
  },
  scrollView: {
    flex: 1,
  },
  scrollContent: {
    paddingBottom: 100,
  },
  profileSection: {
    alignItems: 'center',
    paddingHorizontal: 16,
    paddingVertical: 16,
  },
  avatarContainer: {
    marginBottom: 16,
  },
  avatar: {
    width: 128,
    height: 128,
    borderRadius: 64,
    borderWidth: 2,
    borderColor: 'rgba(168, 85, 247, 0.8)',
  },
  userInfo: {
    alignItems: 'center',
    gap: 4,
  },
  username: {
    fontSize: 22,
    fontWeight: 'bold',
    color: TEXT_PRIMARY,
    textAlign: 'center',
  },
  bio: {
    fontSize: 16,
    color: TEXT_SECONDARY,
    textAlign: 'center',
    marginTop: 4,
  },
  statsSection: {
    flexDirection: 'row',
    justifyContent: 'center',
    gap: 32,
    paddingHorizontal: 16,
    paddingVertical: 12,
  },
  statItem: {
    alignItems: 'center',
    gap: 4,
  },
  statNumber: {
    fontSize: 20,
    fontWeight: 'bold',
    color: TEXT_PRIMARY,
  },
  statLabel: {
    fontSize: 14,
    color: TEXT_SECONDARY,
  },
  editButtonContainer: {
    alignItems: 'center',
    paddingHorizontal: 16,
    paddingTop: 8,
    paddingBottom: 16,
  },
  editButton: {
    paddingHorizontal: 24,
    paddingVertical: 10,
    borderRadius: 9999,
    borderWidth: 1,
    borderColor: 'rgba(168, 85, 247, 0.5)',
  },
  editButtonText: {
    fontSize: 14,
    fontWeight: 'bold',
    color: TEXT_PRIMARY,
  },
  watchlistsHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: 16,
    paddingVertical: 12,
  },
  watchlistsTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    color: TEXT_PRIMARY,
  },
  viewModeButtons: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  viewModeButton: {
    padding: 8,
    borderRadius: 9999,
  },
  viewModeButtonActive: {
    backgroundColor: 'rgba(255, 255, 255, 0.1)',
  },
  watchlistsGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    paddingHorizontal: 12,
    gap: 12,
  },
  watchlistCard: {
    width: '47%',
    backgroundColor: SURFACE_DARK,
    borderRadius: 12,
    overflow: 'hidden',
    borderWidth: 1,
    borderColor: 'rgba(255, 255, 255, 0.1)',
  },
  postersContainer: {
    aspectRatio: 3 / 2,
    position: 'relative',
    overflow: 'hidden',
  },
  posterFront: {
    position: 'absolute',
    width: '66%',
    height: '100%',
    borderRadius: 8,
    zIndex: 2,
  },
  posterBack: {
    position: 'absolute',
    width: '66%',
    height: '100%',
    right: '-10%',
    top: '5%',
    borderRadius: 8,
    zIndex: 1,
  },
  watchlistCardInfo: {
    padding: 12,
    paddingTop: 0,
  },
  watchlistCardTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    color: TEXT_PRIMARY,
    marginBottom: 4,
  },
  watchlistMeta: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 6,
  },
  watchlistMetaText: {
    fontSize: 12,
    color: TEXT_SECONDARY,
  },
});
