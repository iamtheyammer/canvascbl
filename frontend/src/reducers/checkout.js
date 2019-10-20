import {CHECKOUT_GOT_CHECKOUT_TOKEN} from "../actions/checkout";

export default function checkout(state = {}, action) {
  switch(action.type) {
    case CHECKOUT_GOT_CHECKOUT_TOKEN:
      return {
        ...state,
        ...{
          checkoutSession: action.checkoutSession
        }
      };
    default:
      return state;
  }
}