# 🍏 SugAIr

Your local, AI-powered nutritional assistant for diabetes management. Fully private, no cloud API keys needed. 

### 🛠 Tech Stack
* **App:** React Native (Expo)
* **API:** Go
* **AI Engine:** Ollama (Local Multimodal - Llava)

---

### 🚀 Quick Start

To run this project locally, you'll need 3 terminal tabs open.

**1. Spin up the AI (Ollama)**
Make sure you have Ollama installed on your machine. Go to the `ollama` folder and build the custom assistant:
```bash
ollama create llava-diabetes -f modelfile
ollama run llava-diabetes
```

**2. Start the Backend (Go)**
Open a new terminal tab, navigate to the Go server folder, and run it:
```bash
cd go
go run .
```
*It should start listening on `http://localhost:8080`.*

**3. Launch the App (Expo)**
Open a third terminal, navigate to the React Native folder, install the dependencies, and start the bundler:
```bash
cd React_native/front-end
npm install
npx expo start
```

📱 **To test on your phone:** Download the **Expo Go** app (iOS/Android), scan the QR code generated in the terminal, and make sure your phone and computer are on the same Wi-Fi network. 

---
*Built for my final degree project (TFG).* 🎓
