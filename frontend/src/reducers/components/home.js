import {
  HOME_SET_SIGN_IN_BUTTON_AVAILABILITY,
  HOME_SET_GET_SESSION_ID
} from '../../actions/components/home';

export default function home(state = {}, action) {
  switch (action.type) {
    case HOME_SET_GET_SESSION_ID:
      return {
        ...state,
        getSessionId: action.id
      };
    case HOME_SET_SIGN_IN_BUTTON_AVAILABILITY:
      return {
        ...state,
        signInButtonAvailability: !!action.availability
      };
    default:
      return state;
  }
}
