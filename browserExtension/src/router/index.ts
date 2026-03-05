import { createRouter, createWebHashHistory } from "vue-router";
import PasswordListPage from "@/components/PasswordListPage.vue";
import ConfigurationPage from "@/components/ConfigurationPage.vue";
import InfoPage from "@/components/InfoPage.vue";

const routes = [
  {
    path: "/",
    component: PasswordListPage,
  },
  {
    path: "/config",
    component: ConfigurationPage,
  },
  {
    path: "/info",
    component: InfoPage,
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

export default router;
