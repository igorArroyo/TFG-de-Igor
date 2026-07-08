// hooks/useHealthCheck.js
import axios from "axios";
import { useState } from "react";

export const useHealthCheck = () => {
  // Estado que almacena el código HTTP de respuesta o un indicador de error de red
  const [status, setStatus] = useState(null);

  // Realiza la petición GET al endpoint de comprobación de estado del servidor
  const checkHealth = async () => {
    try {
      /*
       * NOTA DE CONFIGURACIÓN DE RED (React Native):
       * - Web: 'http://localhost:8080'
       * - Emulador Android: 'http://10.0.2.2:8080'
       * - Dispositivo Físico: IP local del host (ej. 'http://192.168.x.x:8080')
       */
      const response = await axios.get(
        "http://localhost:8080/api/v1/healthcheck",
      );

      // Actualiza el estado con el código de respuesta exitoso (ej. 200 OK)
      setStatus(response.status);
    } catch (error) {
      // Diferencia entre errores devueltos por el servidor y errores de conectividad
      if (error.response) {
        // El servidor respondió con un código de error (ej. 404, 500)
        setStatus(error.response.status);
      } else {
        // El servidor está inaccesible o no hay conexión de red
        setStatus("ERROR_NETWORK");
      }
    }
  };

  // Expone el estado y el método disparador para ser consumidos por los componentes
  return { status, checkHealth };
};
