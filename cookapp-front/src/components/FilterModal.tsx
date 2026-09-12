import { API_URL } from '@/constants/api';
import { Colors } from '@/constants/colors';
import { useAuth } from '@/context/AuthContext';
import { useEffect, useState } from 'react';
import { Modal, ScrollView, Text, TextInput, TouchableOpacity, View } from 'react-native';

type Ingredient = {
  id: number;
  name: string;
};

type Filters = {
  ingredient: string;
  sort: string;
  liked: boolean;
  favorited: boolean;
};

type Props = {
  visible: boolean;
  onClose: () => void;
  filters: Filters;
  onApply: (filters: Filters) => void;
};

export default function FilterModal({ visible, onClose, filters, onApply }: Props) {
  const { token } = useAuth();
  const [ingredients, setIngredients] = useState<Ingredient[]>([]);
  const [ingredientSearch, setIngredientSearch] = useState('');
  const [localFilters, setLocalFilters] = useState<Filters>(filters);

  useEffect(() => {
    if (!visible) return;
    setLocalFilters(filters);
    fetch(`${API_URL}/ingredients`)
      .then((res) => res.json())
      .then((data) => setIngredients(data ?? []))
      .catch(() => {});
  }, [visible]);

  const toggleSort = (sort: string) => {
    setLocalFilters((prev) => ({ ...prev, sort: prev.sort === sort ? '' : sort }));
  };

  const reset = () => {
    setLocalFilters({ ingredient: '', sort: '', liked: false, favorited: false });
  };

  const filteredIngredients = ingredients.filter((i) =>
    i.name.toLowerCase().includes(ingredientSearch.toLowerCase())
  );

  const chipStyle = (active: boolean) => ({
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: 20,
    borderWidth: 0.5,
    borderColor: active ? Colors.primary : Colors.border,
    backgroundColor: active ? Colors.primary : Colors.white,
  });

  const chipText = (active: boolean) => ({
    fontSize: 13,
    color: active ? Colors.white : Colors.textDark,
  });

  return (
    <Modal visible={visible} transparent animationType="slide" onRequestClose={onClose}>
      <TouchableOpacity
        style={{ flex: 1, backgroundColor: 'rgba(0,0,0,0.4)' }}
        onPress={onClose}
        activeOpacity={1}
      />
      <View style={{
        backgroundColor: Colors.white,
        padding: 20,
        borderTopLeftRadius: 20,
        borderTopRightRadius: 20,
        borderTopWidth: 0.5,
        borderColor: Colors.border,
      }}>
        <View style={{ flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
          <Text style={{ fontSize: 16, fontWeight: '600', color: Colors.textDark }}>Filtres</Text>
          <TouchableOpacity onPress={reset}>
            <Text style={{ fontSize: 13, color: Colors.textMedium }}>Réinitialiser</Text>
          </TouchableOpacity>
        </View>

        <Text style={{ fontSize: 13, color: Colors.textDark, marginBottom: 8 }}>Trier par</Text>
        <View style={{ flexDirection: 'row', gap: 8, marginBottom: 16 }}>
          {[
            { label: 'Plus likés', value: 'most_liked' },
            { label: 'Plus favoris', value: 'most_favorited' },
            { label: 'Plus récents', value: 'newest' },
          ].map((option) => (
            <TouchableOpacity
              key={option.value}
              onPress={() => toggleSort(option.value)}
              style={chipStyle(localFilters.sort === option.value)}
            >
              <Text style={chipText(localFilters.sort === option.value)}>{option.label}</Text>
            </TouchableOpacity>
          ))}
        </View>

        {token && (
          <>
            <Text style={{ fontSize: 13, color: Colors.textDark, marginBottom: 8 }}>Mes recettes</Text>
            <View style={{ flexDirection: 'row', gap: 8, marginBottom: 16 }}>
              <TouchableOpacity
                onPress={() => setLocalFilters((prev) => ({ ...prev, liked: !prev.liked }))}
                style={chipStyle(localFilters.liked)}
              >
                <Text style={chipText(localFilters.liked)}>❤️ Likées</Text>
              </TouchableOpacity>
              <TouchableOpacity
                onPress={() => setLocalFilters((prev) => ({ ...prev, favorited: !prev.favorited }))}
                style={chipStyle(localFilters.favorited)}
              >
                <Text style={chipText(localFilters.favorited)}>⭐ Favorites</Text>
              </TouchableOpacity>
            </View>
          </>
        )}

        <Text style={{ fontSize: 13, color: Colors.textDark, marginBottom: 8 }}>Ingrédient</Text>
        <TextInput
          placeholder="Rechercher un ingrédient..."
          value={ingredientSearch}
          onChangeText={setIngredientSearch}
          style={{
            borderWidth: 0.5,
            borderColor: Colors.border,
            borderRadius: 8,
            padding: 8,
            marginBottom: 8,
            color: Colors.textDark,
            backgroundColor: Colors.background,
          }}
          placeholderTextColor={Colors.textMedium}
        />
        <ScrollView style={{ maxHeight: 120, marginBottom: 16 }}>
          {filteredIngredients.map((ing) => (
            <TouchableOpacity
              key={ing.id}
              onPress={() => setLocalFilters((prev) => ({
                ...prev,
                ingredient: prev.ingredient === ing.name ? '' : ing.name,
              }))}
              style={{
                padding: 8,
                borderRadius: 8,
                backgroundColor: localFilters.ingredient === ing.name ? Colors.border : Colors.white,
              }}
            >
              <Text style={{ color: localFilters.ingredient === ing.name ? Colors.primaryDark : Colors.textDark, fontSize: 13 }}>
                {ing.name}
              </Text>
            </TouchableOpacity>
          ))}
        </ScrollView>

        <TouchableOpacity
          onPress={() => { onApply(localFilters); onClose(); }}
          style={{
            backgroundColor: Colors.primary,
            padding: 12,
            borderRadius: 8,
            alignItems: 'center',
          }}
        >
          <Text style={{ color: Colors.white, fontWeight: '600', fontSize: 15 }}>Appliquer</Text>
        </TouchableOpacity>
      </View>
    </Modal>
  );
}