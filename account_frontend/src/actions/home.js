export const HOME_SET_SIGN_IN_BUTTON_AVAILABILITY =
  'HOME_SET_SIGN_IN_BUTTON_AVAILABILITY';

export const HOME_SET_REDIRECT_OK = 'HOME_SET_REDIRECT_OK';

export const HOME_SET_DESTINATION = 'HOME_SET_DEST';

export function setSignInButtonAvailability(availability) {
  return {
    type: HOME_SET_SIGN_IN_BUTTON_AVAILABILITY,
    availability
  };
}

export function setDestination(dest) {
  return {
    type: HOME_SET_DESTINATION,
    dest
  };
}

export function setRedirectOk(ok) {
  return {
    type: HOME_SET_REDIRECT_OK,
    ok
  };
}
