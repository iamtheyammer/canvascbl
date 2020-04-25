import {
  CANVAS_GET_USER_PROFILE,
  CANVAS_GET_USER_PROFILE_ERROR,
  CANVAS_GOT_USER_PROFILE,
  CANVAS_LOGOUT,
  CANVAS_LOGOUT_ERROR,
  CANVAS_LOGGED_OUT
} from '../actions/canvas';

export default function canvas(state = {}, action) {
  switch (action.type) {
    case CANVAS_GET_USER_PROFILE:
      return {
        ...state,
        loadingUserProfile: true
      };
    case CANVAS_GET_USER_PROFILE_ERROR:
      return {
        ...state,
        loadingUserProfile: false,
        getUserProfileError: action.e
      };
    case CANVAS_GOT_USER_PROFILE:
      return {
        ...state,
        loadingUserProfile: false,
        profile: action.profile
      };
    case CANVAS_LOGOUT:
      return {
        ...state,
        loadingLogout: true
      };
    case CANVAS_LOGOUT_ERROR:
      return {
        ...state,
        loadingLogout: false,
        logoutError: action.e
      };
    case CANVAS_LOGGED_OUT:
      return {
        loggedOut: true
      };
    default:
      return state;
  }
}
