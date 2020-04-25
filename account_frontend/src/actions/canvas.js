export const CANVAS_GET_USER_PROFILE = 'CANVAS_GET_USER_PROFILE';
export const CANVAS_GET_USER_PROFILE_ERROR = 'CANVAS_GET_USER_PROFILE_ERROR';
export const CANVAS_GOT_USER_PROFILE = 'CANVAS_GOT_USER_PROFILE';

export const CANVAS_LOGOUT = 'CANVAS_LOGOUT';
export const CANVAS_LOGOUT_ERROR = 'CANVAS_LOGOUT_ERROR';
export const CANVAS_LOGGED_OUT = 'CANVAS_LOGGED_OUT';

export function getUserProfile() {
  return {
    type: CANVAS_GET_USER_PROFILE
  };
}

export function getUserProfileError(e) {
  return {
    type: CANVAS_GET_USER_PROFILE_ERROR,
    e
  };
}

export function gotUserProfile(profile) {
  return {
    type: CANVAS_GOT_USER_PROFILE,
    profile
  };
}

export function logout() {
  return {
    type: CANVAS_LOGOUT
  };
}

export function logoutError(e) {
  return {
    type: CANVAS_LOGOUT_ERROR,
    e
  };
}

export function loggedOut() {
  return {
    type: CANVAS_LOGGED_OUT
  };
}
