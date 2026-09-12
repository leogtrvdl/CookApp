import { API_URL } from '@/constants/api';
import { Colors } from '@/constants/colors';
import { useAuth } from '@/context/AuthContext';
import { useEffect, useState } from 'react';
import { FlatList, Modal, Text, TextInput, TouchableOpacity, View } from 'react-native';

type User = {
  id: number;
  username: string;
  avatar_path: string;
};

type Friendship = {
  id: number;
  requester_id: number;
  receiver_id: number;
  status: string;
  friend_username: string;
};

type Props = {
  visible: boolean;
  onClose: () => void;
};

export default function FriendModal({ visible, onClose }: Props) {
  const { token, userID, fetchWithAuth } = useAuth();
  const [view, setView] = useState<'friends' | 'requests'>('friends');
  const [friends, setFriends] = useState<Friendship[]>([]);
  const [requests, setRequests] = useState<Friendship[]>([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [searchResults, setSearchResults] = useState<User[]>([]);

  const loadFriends = () => {
    if (!token) return;
    fetchWithAuth(`${API_URL}/friends`)
      .then((res) => res.json())
      .then((data) => setFriends(data ?? []))
      .catch(() => {});
  };

  const loadRequests = () => {
    if (!token) return;
    fetchWithAuth(`${API_URL}/friends/requests`)
      .then((res) => res.json())
      .then((data) => setRequests(data ?? []))
      .catch(() => {});
  };

  const handleSearch = (text: string) => {
    setSearchQuery(text);
    if (text.trim() === '') {
      setSearchResults([]);
      return;
    }
    fetch(`${API_URL}/users/search?username=${text}`)
      .then((res) => res.json())
      .then((data) => setSearchResults((data ?? []).filter((u: User) => u.id !== userID)))
      .catch(() => {});
  };

  const sendFriendRequest = async (receiverId: number) => {
    await fetchWithAuth(`${API_URL}/friends/request`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ receiver_id: receiverId }),
    });
    setSearchResults([]);
    setSearchQuery('');
  };

  const acceptRequest = async (friendshipId: number) => {
    await fetchWithAuth(`${API_URL}/friends/${friendshipId}/accept`, { method: 'PUT' });
    loadFriends();
    loadRequests();
  };

  const declineRequest = async (friendshipId: number) => {
    await fetchWithAuth(`${API_URL}/friends/${friendshipId}/decline`, { method: 'PUT' });
    loadRequests();
  };

  useEffect(() => {
    if (!visible || !token) return;
    loadFriends();
    loadRequests();
  }, [visible, token]);

  return (
    <Modal visible={visible} animationType="slide" onRequestClose={onClose}>
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
          <Text style={{ fontSize: 17, fontWeight: '600', color: Colors.textDark }}>Amis</Text>
          <TouchableOpacity onPress={onClose}>
            <Text style={{ color: Colors.primaryDark, fontSize: 14 }}>Fermer</Text>
          </TouchableOpacity>
        </View>

        <View style={{ padding: 16 }}>
          <TextInput
            placeholder="Rechercher un utilisateur..."
            value={searchQuery}
            onChangeText={handleSearch}
            style={{
              borderWidth: 0.5,
              borderColor: Colors.border,
              borderRadius: 8,
              padding: 10,
              marginBottom: 12,
              backgroundColor: Colors.white,
              color: Colors.textDark,
            }}
            placeholderTextColor={Colors.textMedium}
          />

          {searchResults.length > 0 && (
            <View style={{
              backgroundColor: Colors.white,
              borderRadius: 8,
              borderWidth: 0.5,
              borderColor: Colors.border,
              marginBottom: 12,
            }}>
              {searchResults.map((user) => (
                <View key={user.id} style={{
                  flexDirection: 'row',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                  padding: 12,
                  borderBottomWidth: 0.5,
                  borderColor: Colors.border,
                }}>
                  <Text style={{ color: Colors.textDark, fontSize: 14 }}>{user.username}</Text>
                  <TouchableOpacity
                    onPress={() => sendFriendRequest(user.id)}
                    style={{
                      backgroundColor: Colors.primary,
                      paddingHorizontal: 12,
                      paddingVertical: 6,
                      borderRadius: 20,
                    }}
                  >
                    <Text style={{ color: Colors.white, fontSize: 13 }}>Ajouter</Text>
                  </TouchableOpacity>
                </View>
              ))}
            </View>
          )}

          <View style={{ flexDirection: 'row', gap: 8, marginBottom: 16 }}>
            <TouchableOpacity
              onPress={() => setView('friends')}
              style={{
                flex: 1,
                padding: 10,
                borderRadius: 8,
                alignItems: 'center',
                backgroundColor: view === 'friends' ? Colors.primary : Colors.white,
                borderWidth: 0.5,
                borderColor: view === 'friends' ? Colors.primary : Colors.border,
              }}
            >
              <Text style={{ color: view === 'friends' ? Colors.white : Colors.textDark, fontSize: 13 }}>
                Amis
              </Text>
            </TouchableOpacity>
            <TouchableOpacity
              onPress={() => setView('requests')}
              style={{
                flex: 1,
                padding: 10,
                borderRadius: 8,
                alignItems: 'center',
                backgroundColor: view === 'requests' ? Colors.primary : Colors.white,
                borderWidth: 0.5,
                borderColor: view === 'requests' ? Colors.primary : Colors.border,
              }}
            >
              <Text style={{ color: view === 'requests' ? Colors.white : Colors.textDark, fontSize: 13 }}>
                Demandes{requests.length > 0 ? ` (${requests.length})` : ''}
              </Text>
            </TouchableOpacity>
          </View>
        </View>

        {view === 'friends' && (
          <FlatList
            data={friends}
            keyExtractor={(item) => item.id.toString()}
            contentContainerStyle={{ paddingHorizontal: 16 }}
            renderItem={({ item }) => (
              <View style={{
                backgroundColor: Colors.white,
                borderRadius: 8,
                borderWidth: 0.5,
                borderColor: Colors.border,
                padding: 12,
                marginBottom: 8,
              }}>
                <Text style={{ color: Colors.textDark, fontSize: 14 }}>{item.friend_username}</Text>
              </View>
            )}
            ListEmptyComponent={
              <Text style={{ color: Colors.textMedium, padding: 16 }}>Aucun ami pour l'instant</Text>
            }
          />
        )}

        {view === 'requests' && (
          <FlatList
            data={requests}
            keyExtractor={(item) => item.id.toString()}
            contentContainerStyle={{ paddingHorizontal: 16 }}
            renderItem={({ item }) => (
              <View style={{
                backgroundColor: Colors.white,
                borderRadius: 8,
                borderWidth: 0.5,
                borderColor: Colors.border,
                padding: 12,
                marginBottom: 8,
                flexDirection: 'row',
                justifyContent: 'space-between',
                alignItems: 'center',
              }}>
                <Text style={{ color: Colors.textDark, fontSize: 14 }}>{item.friend_username}</Text>
                <View style={{ flexDirection: 'row', gap: 8 }}>
                  <TouchableOpacity
                    onPress={() => acceptRequest(item.id)}
                    style={{
                      backgroundColor: Colors.primary,
                      paddingHorizontal: 12,
                      paddingVertical: 6,
                      borderRadius: 20,
                    }}
                  >
                    <Text style={{ color: Colors.white, fontSize: 13 }}>Accepter</Text>
                  </TouchableOpacity>
                  <TouchableOpacity
                    onPress={() => declineRequest(item.id)}
                    style={{
                      borderWidth: 0.5,
                      borderColor: Colors.border,
                      paddingHorizontal: 12,
                      paddingVertical: 6,
                      borderRadius: 20,
                    }}
                  >
                    <Text style={{ color: Colors.textMedium, fontSize: 13 }}>Refuser</Text>
                  </TouchableOpacity>
                </View>
              </View>
            )}
            ListEmptyComponent={
              <Text style={{ color: Colors.textMedium, padding: 16 }}>Aucune demande en attente</Text>
            }
          />
        )}

      </View>
    </Modal>
  );
}