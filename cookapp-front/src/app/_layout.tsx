import { Tabs } from 'expo-router';
import { AuthProvider } from '@/context/AuthContext';
import { Colors } from '@/constants/colors';

export default function Layout() {
  return (
    <AuthProvider>
      <Tabs
        screenOptions={{
          tabBarStyle: {
            backgroundColor: Colors.background,
            borderTopColor: Colors.border,
            borderTopWidth: 0.5,
          },
          tabBarActiveTintColor: Colors.primaryDark,
          tabBarInactiveTintColor: Colors.textMedium,
          headerStyle: {
            backgroundColor: Colors.background,
          },
          headerTintColor: Colors.textDark,
          headerTitleStyle: {
            fontWeight: '500',
          },
        }}
      >
        <Tabs.Screen name="index" options={{ title: 'Accueil' }} />
        <Tabs.Screen name="search" options={{ title: 'Recherche' }} />
        <Tabs.Screen name="profile" options={{ title: 'Profil' }} />
      </Tabs>
    </AuthProvider>
  );
}