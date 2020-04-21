export function set(location) {
  localStorage.loginReturnTo =
    location.pathname + location.hash + location.search;
}

export function get() {
  return localStorage.loginReturnTo;
}

export function clear() {
  localStorage.loginReturnTo = "";
}
