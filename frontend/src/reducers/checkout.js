import {
  CHECKOUT_GOT_CHECKOUT_SESSION,
  CHECKOUT_GOT_PRODUCTS
} from '../actions/checkout';

export default function checkout(state = {}, action) {
  switch (action.type) {
    case CHECKOUT_GOT_PRODUCTS:
      return {
        ...state,
        products: action.products
      };
    case CHECKOUT_GOT_CHECKOUT_SESSION:
      return {
        ...state,
        ...{
          session: action.session
        }
      };
    default:
      return state;
  }
}
