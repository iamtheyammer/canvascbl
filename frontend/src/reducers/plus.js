import { PLUS_GOT_SESSION_INFORMATION } from '../actions/plus';

export default function plus(state = [], action) {
  switch (action.type) {
    case PLUS_GOT_SESSION_INFORMATION:
      return {
        ...state,
        session: action.sessionInformation
      };
    default:
      return state;
  }
}
