import { API_URL } from '@/constants/api';
import { Colors } from '@/constants/colors';
import { useAuth } from '@/context/AuthContext';
import { useBreakpoint } from '@/hooks/useBreakpoint';
import { useState } from 'react';
import { Text, TextInput, TouchableOpacity, View } from 'react-native';

export default function AuthForm() {
  const { login } = useAuth();
  const { isMobile } = useBreakpoint();
  const [isLogin, setIsLogin] = useState(true);
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = async () => {
    setError('');
    try {
      const url = isLogin ? `${API_URL}/users/login` : `${API_URL}/users/register`;
      const body = isLogin ? { username, password } : { username, email, password };

      const res = await fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      });

      if (!res.ok) {
        setError(isLogin ? 'Identifiants incorrects' : "Nom d'utilisateur ou email déjà utilisé");
        return;
      }

      if (isLogin) {
        const data = await res.json();
        login(data.token);
      } else {
        setIsLogin(true);
      }
    } catch {
      setError('Impossible de contacter le serveur, réessayez plus tard');
    }
  };

  return (
    <View style={{ flex: 1, backgroundColor: Colors.background, justifyContent: 'center', alignItems: 'center' }}>
      <View style={{
        width: isMobile ? '90%' : 400,
        backgroundColor: Colors.white,
        borderRadius: 16,
        padding: 24,
        borderWidth: 0.5,
        borderColor: Colors.border,
      }}>
        <Text style={{ fontSize: 22, fontWeight: '700', color: Colors.textDark, marginBottom: 4 }}>
          {isLogin ? 'Connexion' : 'Inscription'}
        </Text>
        <Text style={{ fontSize: 13, color: Colors.textMedium, marginBottom: 24 }}>
          {isLogin ? 'Content de te revoir !' : 'Crée ton compte pour commencer'}
        </Text>

        <Text style={{ fontSize: 13, color: Colors.textDark, marginBottom: 6 }}>Nom d'utilisateur</Text>
        <TextInput
          placeholder="Ton nom d'utilisateur"
          value={username}
          onChangeText={setUsername}
          style={{
            borderWidth: 0.5,
            borderColor: Colors.border,
            borderRadius: 8,
            padding: 10,
            marginBottom: 16,
            color: Colors.textDark,
            backgroundColor: Colors.background,
          }}
          placeholderTextColor={Colors.textMedium}
        />

        {!isLogin && (
          <>
            <Text style={{ fontSize: 13, color: Colors.textDark, marginBottom: 6 }}>Email</Text>
            <TextInput
              placeholder="Ton email"
              value={email}
              onChangeText={setEmail}
              style={{
                borderWidth: 0.5,
                borderColor: Colors.border,
                borderRadius: 8,
                padding: 10,
                marginBottom: 16,
                color: Colors.textDark,
                backgroundColor: Colors.background,
              }}
              placeholderTextColor={Colors.textMedium}
            />
          </>
        )}

        <Text style={{ fontSize: 13, color: Colors.textDark, marginBottom: 6 }}>Mot de passe</Text>
        <TextInput
          placeholder="Ton mot de passe"
          value={password}
          onChangeText={setPassword}
          secureTextEntry
          style={{
            borderWidth: 0.5,
            borderColor: Colors.border,
            borderRadius: 8,
            padding: 10,
            marginBottom: 16,
            color: Colors.textDark,
            backgroundColor: Colors.background,
          }}
          placeholderTextColor={Colors.textMedium}
        />

        {error ? (
          <Text style={{ color: 'red', fontSize: 13, marginBottom: 12 }}>{error}</Text>
        ) : null}

        <TouchableOpacity
          onPress={handleSubmit}
          style={{
            backgroundColor: Colors.primary,
            padding: 12,
            borderRadius: 8,
            alignItems: 'center',
            marginBottom: 12,
          }}
        >
          <Text style={{ color: Colors.white, fontWeight: '600', fontSize: 15 }}>
            {isLogin ? 'Se connecter' : "S'inscrire"}
          </Text>
        </TouchableOpacity>

        <TouchableOpacity onPress={() => setIsLogin(!isLogin)} style={{ alignItems: 'center' }}>
          <Text style={{ color: Colors.textMedium, fontSize: 13 }}>
            {isLogin ? "Pas de compte ? S'inscrire" : 'Déjà un compte ? Se connecter'}
          </Text>
        </TouchableOpacity>
      </View>
    </View>
  );
}