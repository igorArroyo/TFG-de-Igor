import { Ionicons } from "@expo/vector-icons";
import * as ImagePicker from "expo-image-picker";
import React, { useState } from "react";
import {
  ActivityIndicator,
  Image,
  KeyboardAvoidingView,
  Platform,
  ScrollView,
  StyleSheet,
  Text,
  TextInput,
  TouchableOpacity,
  View,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

export default function App() {
  const [inputText, setInputText] = useState("");
  const [loading, setLoading] = useState(false);
  const [status, setStatus] = useState<number | null>(null);
  const [message, setMessage] = useState<string>(
    "Escribe un mensaje o adjunta una imagen para empezar...",
  );

  // Estado para gestionar la visualización local de la imagen y su codificación
  const [selectedImage, setSelectedImage] = useState<{
    uri: string;
    base64: string;
  } | null>(null);

  /**
   * Gestiona los permisos del sistema y la apertura de la galería multimedia.
   * Aplica compresión inicial y extrae la cadena Base64 de la imagen seleccionada.
   */
  const pickImage = async () => {
    const result = await ImagePicker.launchImageLibraryAsync({
      mediaTypes: ["images"],
      allowsEditing: true,
      quality: 0.5,
      base64: true,
    });

    if (!result.canceled && result.assets && result.assets.length > 0) {
      setSelectedImage({
        uri: result.assets[0].uri,
        base64: result.assets[0].base64 || "",
      });
    }
  };

  /**
   * Orquesta la validación de entrada y la comunicación HTTP con el servidor API Gateway.
   * Controla los estados de carga y renderiza la respuesta o las excepciones de red.
   */
  const sendChat = async () => {
    if (!inputText.trim() && !selectedImage) return;

    const textToSend = inputText;
    const imageToSend = selectedImage?.base64 || "";

    setInputText("");
    setSelectedImage(null);
    setLoading(true);
    setStatus(null);

    try {
      const response = await fetch("http://localhost:8080/api/v1/chat", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          content: textToSend,
          image: imageToSend,
        }),
      });

      setStatus(response.status);

      if (response.ok) {
        const data = await response.json();
        setMessage(data.message);
      } else {
        setMessage(`Error del servidor.\nStatus: ${response.status}`);
      }
    } catch (error: any) {
      setStatus(500);
      setMessage(`Error de red:\n${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <SafeAreaView style={styles.container} edges={["bottom"]}>
      <KeyboardAvoidingView
        style={styles.container}
        behavior={Platform.OS === "ios" ? "padding" : "height"}
        keyboardVerticalOffset={Platform.OS === "ios" ? 60 : 0}
      >
        {/* Contenedor principal de mensajes */}
        <ScrollView
          contentContainerStyle={styles.responseContainer}
          style={styles.responseBox}
          keyboardShouldPersistTaps="handled"
        >
          {loading ? (
            <ActivityIndicator size="large" color="#0096D6" />
          ) : (
            <Text
              style={[
                styles.responseText,
                status === 200
                  ? styles.success
                  : status
                    ? styles.error
                    : styles.neutral,
              ]}
            >
              {message}
            </Text>
          )}
        </ScrollView>

        {/* Módulo de previsualización de imagen adjunta */}
        {selectedImage && (
          <View style={styles.previewContainer}>
            <Image
              source={{ uri: selectedImage.uri }}
              style={styles.previewImage}
            />
            <TouchableOpacity
              style={styles.removeImageButton}
              onPress={() => setSelectedImage(null)}
            >
              <Ionicons name="close-circle" size={24} color="#FF3B30" />
            </TouchableOpacity>
          </View>
        )}

        {/* Interfaz de entrada de datos de usuario */}
        <View style={styles.inputRow}>
          <TouchableOpacity
            style={styles.attachButton}
            onPress={pickImage}
            disabled={loading}
          >
            <Ionicons
              name="attach"
              size={28}
              color={loading ? "#ccc" : "#0096D6"}
            />
          </TouchableOpacity>

          <TextInput
            style={styles.textInput}
            placeholder="Escribe tu mensaje..."
            value={inputText}
            onChangeText={setInputText}
            onSubmitEditing={sendChat}
          />

          <TouchableOpacity
            style={styles.sendButton}
            onPress={sendChat}
            disabled={loading || (!inputText.trim() && !selectedImage)}
          >
            <Ionicons
              name="send"
              size={20}
              color={
                loading || (!inputText.trim() && !selectedImage)
                  ? "#ccc"
                  : "#0096D6"
              }
            />
          </TouchableOpacity>
        </View>
      </KeyboardAvoidingView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#ffffff",
  },
  responseBox: {
    flex: 1,
    margin: 20,
    borderWidth: 1,
    borderColor: "#e0e0e0",
    borderRadius: 8,
    backgroundColor: "#f9f9f9",
  },
  responseContainer: {
    padding: 20,
    justifyContent: "center",
    alignItems: "center",
    flexGrow: 1,
  },
  responseText: {
    fontSize: 16,
    textAlign: "left",
    width: "100%",
  },
  success: {
    color: "#333",
  },
  error: {
    color: "#F44336",
    fontWeight: "bold",
    textAlign: "center",
  },
  neutral: {
    fontSize: 16,
    color: "#888",
    textAlign: "center",
  },
  previewContainer: {
    paddingHorizontal: 20,
    paddingTop: 10,
    flexDirection: "row",
    alignItems: "flex-start",
  },
  previewImage: {
    width: 60,
    height: 60,
    borderRadius: 8,
    borderWidth: 1,
    borderColor: "#ccc",
  },
  removeImageButton: {
    marginLeft: -10,
    marginTop: -10,
    backgroundColor: "#fff",
    borderRadius: 12,
  },
  inputRow: {
    flexDirection: "row",
    padding: 15,
    borderTopWidth: 1,
    borderColor: "#e0e0e0",
    alignItems: "center",
    backgroundColor: "#fff",
  },
  attachButton: {
    marginRight: 10,
    justifyContent: "center",
    alignItems: "center",
  },
  textInput: {
    flex: 1,
    borderWidth: 1,
    borderColor: "#ccc",
    borderRadius: 20,
    paddingHorizontal: 15,
    paddingVertical: 10,
    fontSize: 16,
    marginRight: 10,
  },
  sendButton: {
    padding: 10,
    justifyContent: "center",
    alignItems: "center",
  },
});
