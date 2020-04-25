import {
  HOME_SET_DESTINATION,
  HOME_SET_REDIRECT_OK,
  HOME_SET_SIGN_IN_BUTTON_AVAILABILITY
} from '../actions/home';

export default function home(state = {}, action) {
  switch (action.type) {
    case HOME_SET_SIGN_IN_BUTTON_AVAILABILITY:
      return {
        ...state,
        signInButtonAvailability: !!action.availability
      };
    case HOME_SET_DESTINATION:
      return {
        ...state,
        destination: action.dest
      };
    case HOME_SET_REDIRECT_OK:
      return {
        ...state,
        redirectOk: action.ok
      };
    default:
      return state;
  }
}
