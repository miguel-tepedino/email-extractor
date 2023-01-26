import { createRouter, createWebHistory } from "vue-router";
import HomeView from "../views/HomeView.vue";
import { validateIfLogged } from "./guards";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      component: HomeView,
      beforeEnter: validateIfLogged,
    },
    {
      path: "/about",
      name: "about",
      meta: {
        layout: "MainLayout",
      },
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import("../views/MailsView.vue"),
      beforeEnter: validateIfLogged,
    },
  ],
});

export default router;
