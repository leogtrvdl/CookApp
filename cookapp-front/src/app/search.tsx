import FilterModal from '@/components/FilterModal';
import RecipeCard from '@/components/RecipeCard';
import RecipeModal from '@/components/RecipeModal';
import { API_URL } from '@/constants/api';
import { Colors } from '@/constants/colors';
import { useAuth } from '@/context/AuthContext';
import { useBreakpoint } from '@/hooks/useBreakpoint';
import { useEffect, useState } from 'react';
import { FlatList, Text, TextInput, TouchableOpacity, View } from 'react-native';

type Recipe = {
  id: number;
  user_id: number;
  title: string;
  image_path: string;
  like_count: number;
  favorite_count: number;
};

type Filters = {
  ingredient: string;
  sort: string;
  liked: boolean;
  favorited: boolean;
};

export default function Search() {
  const { token } = useAuth();
  const { isDesktop, listColumnWidth } = useBreakpoint();
  const [recipes, setRecipes] = useState<Recipe[]>([]);
  const [title, setTitle] = useState('');
  const [filters, setFilters] = useState<Filters>({
    ingredient: '',
    sort: '',
    liked: false,
    favorited: false,
  });
  const [filterModalVisible, setFilterModalVisible] = useState(false);
  const [selectedRecipeId, setSelectedRecipeId] = useState<number | null>(null);
  const [error, setError] = useState('');

  const search = (currentTitle: string, currentFilters: Filters) => {
    const params = new URLSearchParams();
    if (currentTitle.trim() !== '') params.append('title', currentTitle);
    if (currentFilters.ingredient !== '') params.append('ingredient', currentFilters.ingredient);
    if (currentFilters.sort !== '') params.append('sort', currentFilters.sort);
    if (currentFilters.liked) params.append('liked', 'true');
    if (currentFilters.favorited) params.append('favorited', 'true');

    const url = `${API_URL}/recipes/search?${params.toString()}`;

    fetch(url, token ? { headers: { Authorization: `Bearer ${token}` } } : {})
      .then((res) => res.json())
      .then((data) => setRecipes(data ?? []))
      .catch(() => setError('Impossible de charger les recettes'));
  };

  useEffect(() => {
    search(title, filters);
  }, []);

  const handleTitleChange = (text: string) => {
    setTitle(text);
    search(text, filters);
  };

  const handleApplyFilters = (newFilters: Filters) => {
    setFilters(newFilters);
    search(title, newFilters);
  };

  const activeFilterCount = [
    filters.ingredient !== '',
    filters.sort !== '',
    filters.liked,
    filters.favorited,
  ].filter(Boolean).length;

  if (error) return <Text>{error}</Text>;

  return (
    <View style={{ flex: 1, flexDirection: 'row', backgroundColor: Colors.background }}>

      <View style={{ width: listColumnWidth, flex: isDesktop ? undefined : 1 }}>
        <View style={{ flexDirection: 'row', padding: 12, gap: 8 }}>
          <TextInput
            placeholder="Rechercher une recette..."
            value={title}
            onChangeText={handleTitleChange}
            style={{
              flex: 1,
              borderWidth: 0.5,
              borderColor: Colors.border,
              borderRadius: 8,
              padding: 8,
              backgroundColor: Colors.white,
              color: Colors.textDark,
            }}
          />
          <TouchableOpacity
            onPress={() => setFilterModalVisible(true)}
            style={{
              justifyContent: 'center',
              paddingHorizontal: 12,
              paddingVertical: 8,
              borderRadius: 8,
              borderWidth: 0.5,
              borderColor: Colors.border,
              backgroundColor: Colors.white,
            }}
          >
            <Text style={{ color: Colors.textDark, fontSize: 13 }}>
              Filtres{activeFilterCount > 0 ? ` (${activeFilterCount})` : ''}
            </Text>
          </TouchableOpacity>
        </View>

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
          ListEmptyComponent={
            <Text style={{ color: Colors.textMedium, padding: 16 }}>Aucune recette trouvée</Text>
          }
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
              onClose={() => {
                setSelectedRecipeId(null);
                search(title, filters);
              }}
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
          onClose={() => {
            setSelectedRecipeId(null);
            search(title, filters);
          }}
        />
      )}

      <FilterModal
        visible={filterModalVisible}
        onClose={() => setFilterModalVisible(false)}
        filters={filters}
        onApply={handleApplyFilters}
      />

    </View>
  );
}