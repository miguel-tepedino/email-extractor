import { defineStore } from "pinia";
import { useLocalStorage } from "@vueuse/core";

export default defineStore("AuthStore", () => {
  const validated = useLocalStorage("validated", false);

  function login(payload: { email: string; password: string }) {
    if (payload.email == "admin@gmail.com" && payload.password == "12345678") {
      validated.value = true;
    } else {
      throw Error("Invalid credentials");
    }
  }

  function logout() {
    validated.value = false;
  }

  return {
    validated,
    login,
    logout,
  };
});
