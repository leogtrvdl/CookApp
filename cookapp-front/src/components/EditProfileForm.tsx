import { API_URL } from '@/constants/api';
import { Colors } from '@/constants/colors';
import { useAuth } from '@/context/AuthContext';
import { useBreakpoint } from '@/hooks/useBreakpoint';
import { useState } from 'react';
import { ScrollView, Text, TextInput, TouchableOpacity, View } from 'react-native';

type User = {
  username: string;
  email: string;
  bio: string;
};

type Props = {
  user: User;
  onClose: () => void;
};

export default function EditProfileForm({ user, onClose }: Props) {
  const { fetchWithAuth } = useAuth();
  const { isMobile } = useBreakpoint();
  const [username, setUsername] = useState(user.username);
  const [email, setEmail] = useState(user.email);
  const [bio, setBio] = useState(user.bio);
  const [error, setError] = useState('');

  const handleSubmit = async () => {
    setError('');
    try {
      const res = await fetchWithAuth(`${API_URL}/users/me`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, email, bio }),
      });
      if (!res.ok) { setError('Erreur lors de la mise à jour'); return; }
      onClose();
    } catch {
      setError('Impossible de contacter le serveur');
    }
  };

  const inputStyle = {
    borderWidth: 0.5,
    borderColor: Colors.border,
    borderRadius: 8,
    padding: 10,
    color: Colors.textDark,
    backgroundColor: Colors.background,
    fontSize: 14,
  };

  return (
    <View style={{ flex: 1, backgroundColor: Colors.background }}>
      <View style={{
        backgroundColor: Colors.white,
        borderBottomWidth: 0.5,
        borderColor: Colors.border,
        padding: 16,
        paddingTop: 24,
        flexDirection: 'row',
        alignItems: 'center',
        justifyContent: 'space-between',
      }}>
        <Text style={{ fontSize: 17, fontWeight: '600', color: Colors.textDark }}>Modifier le profil</Text>
        <TouchableOpacity onPress={onClose}>
          <Text style={{ color: Colors.primaryDark, fontSize: 14 }}>Annuler</Text>
        </TouchableOpacity>
      </View>

      <ScrollView contentContainerStyle={{ padding: 16, maxWidth: isMobile ? undefined : 480, alignSelf: 'center', width: '100%' }}>

        <Text style={{ fontSize: 13, color: Colors.textDark, marginBottom: 6 }}>Nom d'utilisateur</Text>
        <TextInput
          placeholder="Ton nom d'utilisateur"
          value={username}
          onChangeText={setUsername}
          style={{ ...inputStyle, marginBottom: 16 }}
          placeholderTextColor={Colors.textMedium}
        />

        <Text style={{ fontSize: 13, color: Colors.textDark, marginBottom: 6 }}>Email</Text>
        <TextInput
          placeholder="Ton email"
          value={email}
          onChangeText={setEmail}
          style={{ ...inputStyle, marginBottom: 16 }}
          placeholderTextColor={Colors.textMedium}
        />

        <Text style={{ fontSize: 13, color: Colors.textDark, marginBottom: 6 }}>Bio</Text>
        <TextInput
          placeholder="Parle-nous de toi..."
          value={bio}
          onChangeText={setBio}
          multiline
          style={{ ...inputStyle, minHeight: 100, marginBottom: 24 }}
          placeholderTextColor={Colors.textMedium}
        />

        {error ? <Text style={{ color: 'red', fontSize: 13, marginBottom: 12 }}>{error}</Text> : null}

        <TouchableOpacity
          onPress={handleSubmit}
          style={{
            backgroundColor: Colors.primary,
            padding: 14,
            borderRadius: 8,
            alignItems: 'center',
          }}
        >
          <Text style={{ color: Colors.white, fontWeight: '600', fontSize: 15 }}>Sauvegarder</Text>
        </TouchableOpacity>

      </ScrollView>
    </View>
  );
}