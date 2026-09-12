import { API_URL } from '@/constants/api';
import { Colors } from '@/constants/colors';
import { useAuth } from '@/context/AuthContext';
import { useEffect, useState } from 'react';
import { Image, Modal, ScrollView, Text, TouchableOpacity, View } from 'react-native';

type Ingredient = {
  id: number;
  name: string;
  quantity: string;
  unit: string;
};

type Step = {
  id: number;
  order: number;
  description: string;
};

type Recipe = {
  id: number;
  title: string;
  image_path: string;
  ingredients: Ingredient[];
  steps: Step[];
  like_count: number;
  favorite_count: number;
};

type Props = {
  recipeId: number | null;
  onClose: () => void;
  embedded?: boolean;
};

export default function RecipeModal({ recipeId, onClose, embedded = false }: Props) {
  const { token, fetchWithAuth } = useAuth();
  const [recipe, setRecipe] = useState<Recipe | null>(null);
  const [error, setError] = useState('');
  const [isLiked, setIsLiked] = useState(false);
  const [isFavorited, setIsFavorited] = useState(false);
  const [likeCount, setLikeCount] = useState(0);
  const [favoriteCount, setFavoriteCount] = useState(0);

  useEffect(() => {
    if (!recipeId) return;
    setRecipe(null);
    setError('');

    fetch(`${API_URL}/recipes/${recipeId}`)
      .then((res) => res.json())
      .then((data) => {
        setRecipe(data);
        setLikeCount(data.like_count);
        setFavoriteCount(data.favorite_count);
      })
      .catch(() => setError('Impossible de charger la recette'));

    if (token) {
      fetchWithAuth(`${API_URL}/recipes/${recipeId}/isliked`)
        .then((res) => res.json())
        .then((data) => setIsLiked(data.is_liked));

      fetchWithAuth(`${API_URL}/recipes/${recipeId}/isfavorited`)
        .then((res) => res.json())
        .then((data) => setIsFavorited(data.is_favorited));
    }
  }, [recipeId, token]);

  const toggleLike = async () => {
    if (!token) return;
    const method = isLiked ? 'DELETE' : 'POST';
    await fetchWithAuth(`${API_URL}/recipes/${recipeId}/like`, { method });
    setIsLiked(!isLiked);
    setLikeCount(isLiked ? likeCount - 1 : likeCount + 1);
  };

  const toggleFavorite = async () => {
    if (!token) return;
    const method = isFavorited ? 'DELETE' : 'POST';
    await fetchWithAuth(`${API_URL}/recipes/${recipeId}/favorite`, { method });
    setIsFavorited(!isFavorited);
    setFavoriteCount(isFavorited ? favoriteCount - 1 : favoriteCount + 1);
  };

  const content = (
    <ScrollView style={{ flex: 1 }} contentContainerStyle={{ padding: 16 }}>
      <TouchableOpacity onPress={onClose} style={{ marginBottom: 12 }}>
        <Text style={{ color: Colors.primaryDark, fontSize: 14 }}>← Fermer</Text>
      </TouchableOpacity>

      {error ? (
        <Text style={{ color: 'red' }}>{error}</Text>
      ) : !recipe ? (
        <Text style={{ color: Colors.textMedium }}>Chargement...</Text>
      ) : (
        <View>
          <Text style={{ fontSize: 20, fontWeight: '600', color: Colors.textDark, marginBottom: 12 }}>
            {recipe.title}
          </Text>

          <View style={{ alignItems: 'flex-start', marginBottom: 16 }}>
            <Image
              source={
                recipe.image_path
                  ? { uri: recipe.image_path }
                  : { uri: `${API_URL}/uploads/defaults/default-recipe.svg` }
              }
              style={{ width: 310, height: 240, borderRadius: 10 }}
              resizeMode="cover"
            />
          </View>

          <View style={{ flexDirection: 'row', gap: 16, marginBottom: 20 }}>
            <TouchableOpacity onPress={toggleLike} style={{ flexDirection: 'row', alignItems: 'center', gap: 4 }}>
              <Text style={{ fontSize: 18 }}>{isLiked ? '❤️' : '🤍'}</Text>
              <Text style={{ color: Colors.textMedium }}>{likeCount}</Text>
            </TouchableOpacity>
            <TouchableOpacity onPress={toggleFavorite} style={{ flexDirection: 'row', alignItems: 'center', gap: 4 }}>
              <Text style={{ fontSize: 18 }}>{isFavorited ? '⭐' : '☆'}</Text>
              <Text style={{ color: Colors.textMedium }}>{favoriteCount}</Text>
            </TouchableOpacity>
          </View>

          <Text style={{ fontSize: 16, fontWeight: '600', color: Colors.textDark, marginBottom: 8 }}>
            Ingrédients
          </Text>
          {recipe.ingredients?.map((ing) => (
            <Text key={ing.id} style={{ color: Colors.textMedium, marginBottom: 4 }}>
              • {ing.name} — {ing.quantity} {ing.unit}
            </Text>
          ))}

          <Text style={{ fontSize: 16, fontWeight: '600', color: Colors.textDark, marginTop: 16, marginBottom: 8 }}>
            Étapes
          </Text>
          {recipe.steps?.map((step) => (
            <View key={step.id} style={{ marginBottom: 10 }}>
              <Text style={{ color: Colors.textDark, fontWeight: '500' }}>Étape {step.order}</Text>
              <Text style={{ color: Colors.textMedium }}>{step.description}</Text>
            </View>
          ))}
        </View>
      )}
    </ScrollView>
  );

  // Mode embedded : panneau inline sur desktop
  if (embedded) {
    return (
      <View style={{ flex: 1, backgroundColor: Colors.background }}>
        {content}
      </View>
    );
  }

  // Mode modal : plein écran sur mobile/tablet
  return (
    <Modal
      visible={recipeId !== null}
      animationType="slide"
      onRequestClose={onClose}
    >
      <View style={{ flex: 1, backgroundColor: Colors.background }}>
        {content}
      </View>
    </Modal>
  );
}