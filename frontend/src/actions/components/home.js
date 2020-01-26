export const HOME_SET_GET_SESSION_ID = 'HOME_SET_GET_SESSION_ID';

export const HOME_SET_SIGN_IN_BUTTON_AVAILABILITY =
  'HOME_SET_SIGN_IN_BUTTON_AVAILABILITY';

export function setGetSessionId(id) {
  return {
    type: HOME_SET_GET_SESSION_ID,
    id
  };
}

export function setSigninButtonAvailability(availability) {
  return {
    type: HOME_SET_SIGN_IN_BUTTON_AVAILABILITY,
    availability
  };
}
