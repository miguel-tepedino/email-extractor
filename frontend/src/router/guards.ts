import type { RouteLocation } from "vue-router";
import authstore from "../stores/auth";

export function validateIfLogged(
  to: RouteLocation,
  from: RouteLocation,
  next: Function
) {
  const validated = authstore().validated;
  if (validated && to.name !== "about") {
    return next("/about");
  } else if (!validated && to.name !== "home") {
    return next("/");
  } else {
    return next();
  }
}
