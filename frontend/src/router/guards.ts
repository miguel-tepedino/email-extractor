import type { RouteLocation } from "vue-router";
import authstore from "../stores/auth";

export function validateIfLogged(
  to: RouteLocation,
  from: RouteLocation,
  next: Function
) {
  const validated = authstore().validated;
  if (validated && to.name !== "mails") {
    return next("/mails");
  } else if (!validated && to.name !== "home") {
    return next("/");
  } else {
    return next();
  }
}
