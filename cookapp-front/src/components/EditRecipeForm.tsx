import { API_URL } from '@/constants/api';
import { Colors } from '@/constants/colors';
import { useAuth } from '@/context/AuthContext';
import { useBreakpoint } from '@/hooks/useBreakpoint';
import * as ImagePicker from 'expo-image-picker';
import { useState } from 'react';
import { Image, ScrollView, Text, TextInput, TouchableOpacity, View } from 'react-native';

type Ingredient = {
  name: string;
  quantity: string;
  unit: string;
};

type Step = {
  description: string;
};

type Recipe = {
  id: number;
  title: string;
  image_path: string;
  ingredients: { id: number; name: string; quantity: string; unit: string }[];
  steps: { id: number; order: number; description: string }[];
};

type Props = {
  recipe: Recipe;
  onClose: () => void;
};

export default function EditRecipeForm({ recipe, onClose }: Props) {
  const { fetchWithAuth } = useAuth();
  const { isMobile } = useBreakpoint();
  const [title, setTitle] = useState(recipe.title);
  const [ingredients, setIngredients] = useState<Ingredient[]>(
    recipe.ingredients.map((i) => ({ name: i.name, quantity: i.quantity, unit: i.unit }))
  );
  const [steps, setSteps] = useState<Step[]>(
    recipe.steps.map((s) => ({ description: s.description }))
  );
  const [image, setImage] = useState<{ uri: string; name: string; type: string } | null>(null);
  const [error, setError] = useState('');

  const updateIngredient = (index: number, field: keyof Ingredient, value: string) => {
    const updated = [...ingredients];
    updated[index][field] = value;
    setIngredients(updated);
  };

  const addIngredient = () => setIngredients([...ingredients, { name: '', quantity: '', unit: '' }]);
  const removeIngredient = (index: number) => setIngredients(ingredients.filter((_, i) => i !== index));

  const updateStep = (index: number, value: string) => {
    const updated = [...steps];
    updated[index].description = value;
    setSteps(updated);
  };

  const addStep = () => setSteps([...steps, { description: '' }]);
  const removeStep = (index: number) => setSteps(steps.filter((_, i) => i !== index));

  const pickImage = async () => {
    const result = await ImagePicker.launchImageLibraryAsync({
      mediaTypes: ImagePicker.MediaTypeOptions.Images,
      allowsEditing: true,
      quality: 0.8,
    });
    if (!result.canceled) {
      const asset = result.assets[0];
      const filename = asset.uri.split('/').pop() ?? 'image.jpg';
      const type = asset.mimeType ?? 'image/jpeg';
      setImage({ uri: asset.uri, name: filename, type });
    }
  };

  const handleSubmit = async () => {
    setError('');
    try {
      const formData = new FormData();
      formData.append('title', title);
      formData.append('ingredients', JSON.stringify(ingredients.filter((i) => i.name.trim() !== '')));
      formData.append('steps', JSON.stringify(
        steps.filter((s) => s.description.trim() !== '').map((s, index) => ({ ...s, order: index + 1 }))
      ));
      if (image) {
        const response = await fetch(image.uri);
        const blob = await response.blob();
        formData.append('image', blob, image.name);
      }
      const res = await fetchWithAuth(`${API_URL}/recipes/${recipe.id}`, { method: 'PUT', body: formData });
      if (!res.ok) { setError('Erreur lors de la modification'); return; }
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
        <Text style={{ fontSize: 17, fontWeight: '600', color: Colors.textDark }}>Modifier la recette</Text>
        <TouchableOpacity onPress={onClose}>
          <Text style={{ color: Colors.primaryDark, fontSize: 14 }}>Annuler</Text>
        </TouchableOpacity>
      </View>

      <ScrollView contentContainerStyle={{ padding: 16, maxWidth: isMobile ? undefined : 640, alignSelf: 'center', width: '100%' }}>

        <Text style={{ fontSize: 13, color: Colors.textDark, marginBottom: 6 }}>Titre</Text>
        <TextInput
          placeholder="Nom de la recette"
          value={title}
          onChangeText={setTitle}
          style={{ ...inputStyle, marginBottom: 20 }}
          placeholderTextColor={Colors.textMedium}
        />

        <TouchableOpacity
          onPress={pickImage}
          style={{
            borderWidth: 0.5,
            borderColor: Colors.border,
            borderRadius: 8,
            padding: 12,
            alignItems: 'center',
            marginBottom: 12,
            backgroundColor: Colors.white,
          }}
        >
          <Text style={{ color: Colors.primaryDark, fontSize: 14 }}>📷 Changer l'image</Text>
        </TouchableOpacity>
        {image ? (
          <Image source={{ uri: image.uri }} style={{ width: '100%', aspectRatio: 16 / 9, borderRadius: 8, marginBottom: 20 }} resizeMode="cover" />
        ) : recipe.image_path ? (
          <Image source={{ uri: `${API_URL}/${recipe.image_path}` }} style={{ width: '100%', aspectRatio: 16 / 9, borderRadius: 8, marginBottom: 20 }} resizeMode="cover" />
        ) : null}

        <Text style={{ fontSize: 15, fontWeight: '600', color: Colors.textDark, marginBottom: 12 }}>Ingrédients</Text>
        {ingredients.map((ing, index) => (
          <View key={index} style={{
            backgroundColor: Colors.white,
            borderRadius: 8,
            borderWidth: 0.5,
            borderColor: Colors.border,
            padding: 12,
            marginBottom: 8,
          }}>
            <TextInput
              placeholder="Nom"
              value={ing.name}
              onChangeText={(v) => updateIngredient(index, 'name', v)}
              style={{ ...inputStyle, marginBottom: 8 }}
              placeholderTextColor={Colors.textMedium}
            />
            <View style={{ flexDirection: 'row', gap: 8, marginBottom: 8 }}>
              <TextInput
                placeholder="Quantité"
                value={ing.quantity}
                onChangeText={(v) => updateIngredient(index, 'quantity', v)}
                style={{ ...inputStyle, flex: 1 }}
                placeholderTextColor={Colors.textMedium}
              />
              <TextInput
                placeholder="Unité"
                value={ing.unit}
                onChangeText={(v) => updateIngredient(index, 'unit', v)}
                style={{ ...inputStyle, flex: 1 }}
                placeholderTextColor={Colors.textMedium}
              />
            </View>
            <TouchableOpacity onPress={() => removeIngredient(index)}>
              <Text style={{ color: Colors.textMedium, fontSize: 13 }}>Supprimer</Text>
            </TouchableOpacity>
          </View>
        ))}
        <TouchableOpacity
          onPress={addIngredient}
          style={{
            borderWidth: 0.5,
            borderColor: Colors.border,
            borderRadius: 8,
            padding: 10,
            alignItems: 'center',
            marginBottom: 20,
            backgroundColor: Colors.white,
          }}
        >
          <Text style={{ color: Colors.primaryDark, fontSize: 14 }}>+ Ajouter un ingrédient</Text>
        </TouchableOpacity>

        <Text style={{ fontSize: 15, fontWeight: '600', color: Colors.textDark, marginBottom: 12 }}>Étapes</Text>
        {steps.map((step, index) => (
          <View key={index} style={{
            backgroundColor: Colors.white,
            borderRadius: 8,
            borderWidth: 0.5,
            borderColor: Colors.border,
            padding: 12,
            marginBottom: 8,
          }}>
            <Text style={{ fontSize: 13, fontWeight: '500', color: Colors.textMedium, marginBottom: 6 }}>
              Étape {index + 1}
            </Text>
            <TextInput
              placeholder="Description"
              value={step.description}
              onChangeText={(v) => updateStep(index, v)}
              multiline
              style={{ ...inputStyle, minHeight: 80, marginBottom: 8 }}
              placeholderTextColor={Colors.textMedium}
            />
            <TouchableOpacity onPress={() => removeStep(index)}>
              <Text style={{ color: Colors.textMedium, fontSize: 13 }}>Supprimer</Text>
            </TouchableOpacity>
          </View>
        ))}
        <TouchableOpacity
          onPress={addStep}
          style={{
            borderWidth: 0.5,
            borderColor: Colors.border,
            borderRadius: 8,
            padding: 10,
            alignItems: 'center',
            marginBottom: 24,
            backgroundColor: Colors.white,
          }}
        >
          <Text style={{ color: Colors.primaryDark, fontSize: 14 }}>+ Ajouter une étape</Text>
        </TouchableOpacity>

        {error ? <Text style={{ color: 'red', fontSize: 13, marginBottom: 12 }}>{error}</Text> : null}

        <TouchableOpacity
          onPress={handleSubmit}
          style={{
            backgroundColor: Colors.primary,
            padding: 14,
            borderRadius: 8,
            alignItems: 'center',
            marginBottom: 24,
          }}
        >
          <Text style={{ color: Colors.white, fontWeight: '600', fontSize: 15 }}>Sauvegarder</Text>
        </TouchableOpacity>

      </ScrollView>
    </View>
  );
}