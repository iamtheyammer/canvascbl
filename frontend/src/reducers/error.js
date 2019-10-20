import { CANVAS_PROXY_ERROR, CHECKOUT_ERROR } from '../actions/error';

export default function error(state = {}, action) {
  switch (action.type) {
    case CANVAS_PROXY_ERROR:
      return {
        ...state,
        [action.id]: {
          res: action.res
        }
      };
    case CHECKOUT_ERROR:
      return {
        ...state,
        [action.id]: {
          res: action.res
        }
      };
    default:
      return state;
  }
}
