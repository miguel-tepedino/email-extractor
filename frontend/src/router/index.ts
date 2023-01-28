import { createRouter, createWebHistory } from "vue-router";
import HomeView from "../views/HomeView.vue";
import { validateIfLogged } from "./guards";
import NotFoundView from "../views/NotFoundView.vue";

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
      path: "/mails",
      name: "mails",
      meta: {
        layout: "MainLayout",
      },
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import("../views/MailsView.vue"),
      beforeEnter: validateIfLogged,
    },
    { path: '/:pathMatch(.*)*', name: 'NotFound', component: NotFoundView },
  ],
});

export default router;
