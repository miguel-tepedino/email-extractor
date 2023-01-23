import { createApp, defineAsyncComponent } from "vue";
import { createPinia } from "pinia";

import App from "./App.vue";
import router from "./router";

import "./assets/main.css";

const app = createApp(App);

app.use(createPinia());
app.use(router);
app.component(
  "DefaultLayout",
  defineAsyncComponent({
    loader: () => import("./layouts/DefaultLayout.vue"),
  })
);

app.component(
  "MainLayout",
  defineAsyncComponent({
    loader: () => import("./layouts/MainLayout.vue"),
  })
);

app.mount("#app");
