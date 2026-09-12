import { API_URL } from '@/constants/api';
import { Colors } from '@/constants/colors';
import { useAuth } from '@/context/AuthContext';
import { useBreakpoint } from '@/hooks/useBreakpoint';
import { useEffect, useState } from 'react';
import { Alert, Image, Platform, Text, TouchableOpacity, View } from 'react-native';

type Props = {
  id: number;
  title: string;
  image_path: string;
  likeCount: number;
  favoriteCount: number;
  onPress: () => void;
  onEdit?: () => void;
  onDelete?: () => void;
};

export default function RecipeCard({ id, title, image_path, likeCount, favoriteCount, onPress, onEdit, onDelete }: Props) {
  const { token, fetchWithAuth } = useAuth();
  const { cardImageHeight } = useBreakpoint();
  const [isLiked, setIsLiked] = useState(false);
  const [isFavorited, setIsFavorited] = useState(false);
  const [localLikeCount, setLocalLikeCount] = useState(likeCount);
  const [localFavoriteCount, setLocalFavoriteCount] = useState(favoriteCount);

  useEffect(() => {
    if (!token) return;
    fetchWithAuth(`${API_URL}/recipes/${id}/isliked`)
      .then((res) => res.json())
      .then((data) => setIsLiked(data.is_liked));
    fetchWithAuth(`${API_URL}/recipes/${id}/isfavorited`)
      .then((res) => res.json())
      .then((data) => setIsFavorited(data.is_favorited));
  }, [token, id]);

  const toggleLike = async () => {
    if (!token) return;
    const method = isLiked ? 'DELETE' : 'POST';
    await fetchWithAuth(`${API_URL}/recipes/${id}/like`, { method });
    setIsLiked(!isLiked);
    setLocalLikeCount(isLiked ? localLikeCount - 1 : localLikeCount + 1);
  };

  const toggleFavorite = async () => {
    if (!token) return;
    const method = isFavorited ? 'DELETE' : 'POST';
    await fetchWithAuth(`${API_URL}/recipes/${id}/favorite`, { method });
    setIsFavorited(!isFavorited);
    setLocalFavoriteCount(isFavorited ? localFavoriteCount - 1 : localFavoriteCount + 1);
  };

  const handleDelete = () => {
    if (Platform.OS === 'web') {
      if (window.confirm('Es-tu sûr de vouloir supprimer cette recette ?')) {
        onDelete?.();
      }
    } else {
      Alert.alert(
        'Supprimer la recette',
        'Es-tu sûr de vouloir supprimer cette recette ?',
        [
          { text: 'Annuler', style: 'cancel' },
          { text: 'Supprimer', style: 'destructive', onPress: onDelete },
        ]
      );
    }
  };

  return (
    <View style={{
      backgroundColor: Colors.cardBackground,
      borderRadius: 12,
      borderWidth: 0.5,
      borderColor: Colors.border,
      marginHorizontal: 12,
      marginVertical: 6,
      overflow: 'hidden',
    }}>
      <TouchableOpacity onPress={onPress} activeOpacity={0.85}>
        <Image
          source={
            image_path
              ? { uri: image_path }
              : { uri: `${API_URL}/uploads/defaults/default-recipe.svg` }
          }
          style={{ width: '100%', aspectRatio: 13 / 10 }}
          resizeMode="cover"
        />
      </TouchableOpacity>

      <View style={{ padding: 10 }}>
        <TouchableOpacity onPress={onPress} activeOpacity={0.7}>
          <Text style={{ fontSize: 15, fontWeight: '500', color: Colors.textDark }}>{title}</Text>
        </TouchableOpacity>

        <View style={{ flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', marginTop: 8 }}>
          <View style={{ flexDirection: 'row', gap: 12 }}>
            <TouchableOpacity onPress={toggleLike} style={{ flexDirection: 'row', alignItems: 'center', gap: 4 }}>
              <Text style={{ fontSize: 16 }}>{isLiked ? '❤️' : '🤍'}</Text>
              <Text style={{ fontSize: 13, color: Colors.textMedium }}>{localLikeCount}</Text>
            </TouchableOpacity>
            <TouchableOpacity onPress={toggleFavorite} style={{ flexDirection: 'row', alignItems: 'center', gap: 4 }}>
              <Text style={{ fontSize: 16 }}>{isFavorited ? '⭐' : '☆'}</Text>
              <Text style={{ fontSize: 13, color: Colors.textMedium }}>{localFavoriteCount}</Text>
            </TouchableOpacity>
          </View>

          <View style={{ flexDirection: 'row', gap: 8 }}>
            {onEdit && (
              <TouchableOpacity
                onPress={onEdit}
                style={{
                  paddingHorizontal: 10,
                  paddingVertical: 4,
                  borderRadius: 20,
                  borderWidth: 0.5,
                  borderColor: Colors.border,
                }}
              >
                <Text style={{ fontSize: 12, color: Colors.primaryDark }}>Modifier</Text>
              </TouchableOpacity>
            )}
            {onDelete && (
              <TouchableOpacity
                onPress={handleDelete}
                style={{
                  paddingHorizontal: 10,
                  paddingVertical: 4,
                  borderRadius: 20,
                  backgroundColor: Colors.primary,
                }}
              >
                <Text style={{ fontSize: 12, color: Colors.white }}>Supprimer</Text>
              </TouchableOpacity>
            )}
          </View>
        </View>
      </View>
    </View>
  );
}