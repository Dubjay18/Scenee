import { StyleSheet, View } from 'react-native';
import { Text } from '../../components/Themed';
import { SafeAreaView } from 'react-native-safe-area-context';
import { MaterialIcons } from '@expo/vector-icons';

const BACKGROUND_DARK = '#0A0A0A';
const TEXT_PRIMARY = '#EAEAEA';
const TEXT_SECONDARY = '#A9A9A9';
const PRIMARY_COLOR = '#A855F7';

export default function ActivityScreen() {
  return (
    <View style={styles.container}>
      <SafeAreaView edges={['top']} style={styles.safeArea}>
        <View style={styles.header}>
          <Text style={styles.headerTitle}>Activity</Text>
        </View>
        <View style={styles.content}>
          <MaterialIcons name="notifications" size={64} color={PRIMARY_COLOR} />
          <Text style={styles.title}>Your Activity</Text>
          <Text style={styles.subtitle}>
            See updates from friends and your watchlist activity.
          </Text>
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
    paddingHorizontal: 16,
    paddingVertical: 12,
    backgroundColor: 'rgba(10, 10, 10, 0.9)',
  },
  headerTitle: {
    fontSize: 22,
    fontWeight: 'bold',
    color: TEXT_PRIMARY,
    letterSpacing: -0.5,
  },
  content: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    paddingHorizontal: 32,
    paddingBottom: 100,
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    color: TEXT_PRIMARY,
    marginTop: 16,
  },
  subtitle: {
    fontSize: 16,
    color: TEXT_SECONDARY,
    textAlign: 'center',
    marginTop: 8,
  },
});
