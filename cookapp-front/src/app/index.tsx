import RecipeCard from '@/components/RecipeCard';
import RecipeModal from '@/components/RecipeModal';
import { API_URL } from '@/constants/api';
import { Colors } from '@/constants/colors';
import { useBreakpoint } from '@/hooks/useBreakpoint';
import { useEffect, useState } from 'react';
import { FlatList, Text, View } from 'react-native';

type Recipe = {
  id: number;
  user_id: number;
  title: string;
  image_path: string;
  like_count: number;
  favorite_count: number;
};

export default function Home() {
  const [recipes, setRecipes] = useState<Recipe[]>([]);
  const [selectedRecipeId, setSelectedRecipeId] = useState<number | null>(null);
  const [error, setError] = useState('');
  const { isDesktop, listColumnWidth } = useBreakpoint();

  const loadRecipes = () => {
    fetch(`${API_URL}/recipes`)
      .then((res) => res.json())
      .then(setRecipes)
      .catch(() => setError('Impossible de charger les recettes'));
  };

  useEffect(() => {
    loadRecipes();
  }, []);

  const handleCloseModal = () => {
    setSelectedRecipeId(null);
    loadRecipes();
  };

  if (error) return <Text>{error}</Text>;

  return (
    <View style={{ flex: 1, flexDirection: 'row', backgroundColor: Colors.background }}>

      <View style={{ width: listColumnWidth, flex: isDesktop ? undefined : 1 }}>
        <FlatList
          data={recipes}
          keyExtractor={(item) => item.id.toString()}
          renderItem={({ item }) => (
            <RecipeCard
              id={item.id}
              title={item.title}
              image_path={item.image_path}
              likeCount={item.like_count}
              favoriteCount={item.favorite_count}
              onPress={() => setSelectedRecipeId(item.id)}
            />
          )}
        />
      </View>

      {isDesktop && (
        <View style={{
          flex: 1,
          borderLeftWidth: 0.5,
          borderColor: Colors.border,
          backgroundColor: Colors.background,
          justifyContent: selectedRecipeId ? 'flex-start' : 'center',
          alignItems: selectedRecipeId ? 'stretch' : 'center',
        }}>
          {selectedRecipeId ? (
            <RecipeModal
              recipeId={selectedRecipeId}
              embedded
              onClose={handleCloseModal}
            />
          ) : (
            <Text style={{ color: Colors.textMedium, fontSize: 15 }}>
              Sélectionne une recette 👈
            </Text>
          )}
        </View>
      )}

      {!isDesktop && (
        <RecipeModal
          recipeId={selectedRecipeId}
          onClose={handleCloseModal}
        />
      )}

    </View>
  );
}