import AuthForm from '@/components/AuthForm';
import CreateRecipeForm from '@/components/CreateRecipeForm';
import EditProfileForm from '@/components/EditProfileForm';
import EditRecipeForm from '@/components/EditRecipeForm';
import FriendModal from '@/components/FriendModal';
import RecipeCard from '@/components/RecipeCard';
import RecipeModal from '@/components/RecipeModal';
import { API_URL } from '@/constants/api';
import { Colors } from '@/constants/colors';
import { useAuth } from '@/context/AuthContext';
import { useBreakpoint } from '@/hooks/useBreakpoint';
import { useEffect, useState } from 'react';
import { FlatList, Image, ScrollView, Text, TouchableOpacity, View } from 'react-native';

type User = {
  id: number;
  username: string;
  email: string;
  bio: string;
  avatar_path: string;
};

type Recipe = {
  id: number;
  user_id: number;
  title: string;
  image_path: string;
  like_count: number;
  favorite_count: number;
};

type FullRecipe = {
  id: number;
  title: string;
  image_path: string;
  ingredients: { id: number; name: string; quantity: string; unit: string }[];
  steps: { id: number; order: number; description: string }[];
};

export default function Profile() {
  const { token, logout, fetchWithAuth } = useAuth();
  const { isDesktop, listColumnWidth } = useBreakpoint();
  const [user, setUser] = useState<User | null>(null);
  const [recipes, setRecipes] = useState<Recipe[]>([]);
  const [selectedRecipeId, setSelectedRecipeId] = useState<number | null>(null);
  const [editingRecipe, setEditingRecipe] = useState<FullRecipe | null>(null);
  const [error, setError] = useState('');
  const [isEditing, setIsEditing] = useState(false);
  const [isCreating, setIsCreating] = useState(false);
  const [friendModalVisible, setFriendModalVisible] = useState(false);

  const loadUser = () => {
    if (!token) return;
    fetchWithAuth(`${API_URL}/users/me`)
      .then((res) => res.json())
      .then(setUser)
      .catch(() => setError('Impossible de charger le profil'));
  };

  const loadRecipes = () => {
    if (!token) return;
    fetchWithAuth(`${API_URL}/users/me/recipes`)
      .then((res) => res.json())
      .then((data) => setRecipes(data ?? []))
      .catch(() => setError('Impossible de charger les recettes'));
  };

  const loadRecipeForEdit = (recipeId: number) => {
    fetch(`${API_URL}/recipes/${recipeId}`)
      .then((res) => res.json())
      .then(setEditingRecipe)
      .catch(() => setError('Impossible de charger la recette'));
  };

  const deleteRecipe = async (recipeId: number) => {
    await fetchWithAuth(`${API_URL}/recipes/${recipeId}`, { method: 'DELETE' });
    loadRecipes();
  };

  useEffect(() => {
    loadUser();
    loadRecipes();
  }, [token]);

  if (!token) return <AuthForm />;
  if (error) return <Text style={{ color: 'red', padding: 16 }}>{error}</Text>;
  if (!user) return <Text style={{ padding: 16, color: Colors.textMedium }}>Chargement...</Text>;

  if (isEditing) {
    return (
      <EditProfileForm
        user={user}
        onClose={() => { setIsEditing(false); loadUser(); }}
      />
    );
  }

  if (isCreating) {
    return (
      <CreateRecipeForm
        onClose={() => { setIsCreating(false); loadRecipes(); }}
      />
    );
  }

  if (editingRecipe) {
    return (
      <EditRecipeForm
        recipe={editingRecipe}
        onClose={() => { setEditingRecipe(null); loadRecipes(); }}
      />
    );
  }

  const profileHeader = (
    <ScrollView style={{ padding: 16 }}>
      <View style={{ alignItems: 'center', marginBottom: 20 }}>
        <Image
          source={{ uri: `${API_URL}/${user.avatar_path}` }}
          style={{ width: 80, height: 80, borderRadius: 40, marginBottom: 12 }}
        />
        <Text style={{ fontSize: 18, fontWeight: '600', color: Colors.textDark }}>{user.username}</Text>
        {user.bio ? (
          <Text style={{ color: Colors.textMedium, marginTop: 4, textAlign: 'center' }}>{user.bio}</Text>
        ) : null}
        <Text style={{ color: Colors.textMedium, fontSize: 12, marginTop: 2 }}>{user.email}</Text>
      </View>

      <View style={{ gap: 8 }}>
        <TouchableOpacity
          onPress={() => setIsEditing(true)}
          style={{
            padding: 10,
            borderRadius: 8,
            backgroundColor: Colors.border,
            alignItems: 'center',
          }}
        >
          <Text style={{ color: Colors.primaryDark, fontSize: 14 }}>Modifier le profil</Text>
        </TouchableOpacity>

        <TouchableOpacity
          onPress={() => setIsCreating(true)}
          style={{
            padding: 10,
            borderRadius: 8,
            backgroundColor: Colors.primary,
            alignItems: 'center',
          }}
        >
          <Text style={{ color: Colors.white, fontSize: 14 }}>Créer une recette</Text>
        </TouchableOpacity>

        <TouchableOpacity
          onPress={() => setFriendModalVisible(true)}
          style={{
            padding: 10,
            borderRadius: 8,
            backgroundColor: Colors.primaryDark,
            alignItems: 'center',
          }}
        >
          <Text style={{ color: Colors.white, fontSize: 14 }}>Amis</Text>
        </TouchableOpacity>

        <TouchableOpacity
          onPress={logout}
          style={{
            padding: 10,
            borderRadius: 8,
            alignItems: 'center',
          }}
        >
          <Text style={{ color: Colors.textMedium, fontSize: 14 }}>Se déconnecter</Text>
        </TouchableOpacity>
      </View>
    </ScrollView>
  );

  return (
    <View style={{ flex: 1, flexDirection: 'row', backgroundColor: Colors.background }}>

      {isDesktop ? (
        <>
          <View style={{ width: 260, borderRightWidth: 0.5, borderColor: Colors.border }}>
            {profileHeader}
          </View>

          <View style={{ width: listColumnWidth, flex: undefined }}>
            <Text style={{ fontSize: 15, fontWeight: '600', color: Colors.textDark, padding: 12 }}>
              Mes recettes
            </Text>
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
                  onEdit={() => loadRecipeForEdit(item.id)}
                  onDelete={() => deleteRecipe(item.id)}
                />
              )}
              ListEmptyComponent={
                <Text style={{ color: Colors.textMedium, padding: 16 }}>Aucune recette pour l'instant</Text>
              }
            />
          </View>

          <View style={{
            flex: 1,
            borderLeftWidth: 0.5,
            borderColor: Colors.border,
            justifyContent: selectedRecipeId ? 'flex-start' : 'center',
            alignItems: selectedRecipeId ? 'stretch' : 'center',
          }}>
            {selectedRecipeId ? (
              <RecipeModal
                recipeId={selectedRecipeId}
                embedded
                onClose={() => { setSelectedRecipeId(null); loadRecipes(); }}
              />
            ) : (
              <Text style={{ color: Colors.textMedium, fontSize: 15 }}>
                Sélectionne une recette 👈
              </Text>
            )}
          </View>
        </>
      ) : (
        <View style={{ flex: 1 }}>
          {profileHeader}
          <Text style={{ fontSize: 15, fontWeight: '600', color: Colors.textDark, paddingHorizontal: 16, paddingBottom: 8 }}>
            Mes recettes
          </Text>
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
                onEdit={() => loadRecipeForEdit(item.id)}
                onDelete={() => deleteRecipe(item.id)}
              />
            )}
            ListEmptyComponent={
              <Text style={{ color: Colors.textMedium, padding: 16 }}>Aucune recette pour l'instant</Text>
            }
          />
        </View>
      )}

      <RecipeModal
        recipeId={!isDesktop ? selectedRecipeId : null}
        onClose={() => { setSelectedRecipeId(null); loadRecipes(); }}
      />

      <FriendModal
        visible={friendModalVisible}
        onClose={() => setFriendModalVisible(false)}
      />

    </View>
  );
}